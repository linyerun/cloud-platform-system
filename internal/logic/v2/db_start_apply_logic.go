package v2

import (
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DbStartApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDbStartApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DbStartApplyLogic {
	return &DbStartApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DbStartApplyLogic) DbStartApply(req *types.DbStartApplyReq) (resp *types.CommonResponse, err error) {
	// 判断image是否存在
	ret := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbImageDocument).FindOne(l.ctx, bson.D{{"_id", req.ImageId}})
	if err = ret.Err(); err != nil && err != mongo.ErrNoDocuments {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "get image data error")
	} else if err == mongo.ErrNoDocuments {
		return nil, errorx.NewCodeError(400, "不存在该DbImage")
	}

	// 保存申请单
	form := &models.DbApplicationForm{
		Id:          utils.GetSnowFlakeIdAndBase64(),
		UserId:      l.ctx.Value("user").(*models.User).Id,
		Explanation: req.Explanation,
		ImageId:     req.ImageId,
		DbName:      req.DbName,
		Status:      models.DbApplicationFormStatusIng,
		CreateAt:    time.Now().UnixMilli(),
	}
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbApplicationFormDocument).InsertOne(l.ctx, form)
	if err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "save data error")
	}

	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}
