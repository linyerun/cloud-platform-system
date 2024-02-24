package v2

import (
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LinuxStartApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLinuxStartApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LinuxStartApplyLogic {
	return &LinuxStartApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LinuxStartApplyLogic) LinuxStartApply(req *types.LinuxStartApplyRequest) (resp *types.CommonResponse, err error) {
	// 校验参数
	if req.Memory <= 0 || (req.MemorySwap != -1 && req.MemorySwap <= req.Memory) || req.CoreCount <= 0 || req.DiskSize <= 0 || len(req.ContainerName) == 0 {
		return &types.CommonResponse{Code: 400, Msg: "参数有误"}, nil
	}
	for _, port := range req.ExportPorts {
		if port < 0 || port >= 65535 {
			return &types.CommonResponse{Code: 400, Msg: "参数有误, 端口范围是[0, 65535]"}, nil
		}
	}

	// 判断image_id是否存在
	err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).FindOne(l.ctx, bson.D{{"_id", req.ImageId}}).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		l.Logger.Error(errors.Wrap(err, "get models.LinuxImageDocument error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	} else if err == mongo.ErrNoDocuments {
		return &types.CommonResponse{Code: 400, Msg: "镜像不存在"}, nil
	}

	// 创建Linux服务器申请单
	form := &models.LinuxApplicationForm{
		Id:            utils.GetSnowFlakeIdAndBase64(),
		UserId:        l.ctx.Value("user").(*models.User).Id,
		Explanation:   req.Explanation,
		ImageId:       req.ImageId,
		ContainerName: req.ContainerName,
		ExportPorts:   req.ExportPorts,
		Memory:        req.Memory,
		MemorySwap:    req.MemorySwap,
		CoreCount:     req.CoreCount,
		DiskSize:      req.DiskSize,
		Status:        models.LinuxApplicationFormStatusIng,
		CreateAt:      time.Now().UnixMilli(),
	}

	// 支持插入操作
	if _, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).InsertOne(l.ctx, form); err != nil {
		l.Logger.Error(errors.Wrap(err, "insert LinuxApplicationForm error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 返回结果
	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}
