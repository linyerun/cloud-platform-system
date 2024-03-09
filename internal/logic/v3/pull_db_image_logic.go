package v3

import (
	"cloud-platform-system/internal/asynctask"
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/utils"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullDbImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullDbImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullDbImageLogic {
	return &PullDbImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullDbImageLogic) PullDbImage(req *types.PullDbImageReq) (resp *types.CommonResponse, err error) {
	// image必须存在
	if len(req.ImageName) == 0 {
		return &types.CommonResponse{Code: 400, Msg: "image name error"}, nil
	}

	// 必须指明版本
	if req.ImageTag == "latest" || len(req.ImageTag) == 0 {
		return &types.CommonResponse{Code: 400, Msg: "必须只能具体版本，不能使用latest作为版本号"}, nil
	}

	// port必须符合要求
	if req.Port < 0 || req.Port > 65535 {
		return &types.CommonResponse{Code: 400, Msg: "每个端口范围必须在[0, 65535]里面"}, nil
	}

	// Type必须是存在的
	if req.Type != models.DbImageTypeMySql && req.Type != models.DbImageTypeRedis && req.Type != models.DbImageTypeMongoDb {
		return nil, errorx.NewCodeError(400, "type不符合要求")
	}

	// 避免重复拉取镜像
	err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbImageDocument).FindOne(l.ctx, bson.D{{"name", req.ImageName}, {"tag", req.ImageTag}}).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		l.Logger.Error(errors.Wrap(err, "find data in mongo error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	} else if err == nil {
		return &types.CommonResponse{Code: 400, Msg: "镜像已经存在，无需拉取！"}, nil
	}

	// 使用分布式锁保证不会重复拉取(保证同一个镜像只能有一个管理员创建拉取任务)
	res := l.svcCtx.RedisClient.SetNX(l.ctx, fmt.Sprintf(l.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag), "1", redis.KeepTTL)
	if err = res.Err(); err != nil {
		l.Logger.Error(errors.Wrap(err, "error to apply for a distributed lock"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	ok, err := res.Result()
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get result from redis error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	if !ok {
		return &types.CommonResponse{Code: 401, Msg: "已有管理员拉取此镜像或创建拉取该镜像的异步任务，无所重复拉取镜像！"}, nil
	}

	// 创建异步任务并保存
	creatorId := l.ctx.Value("user").(*models.User).Id
	args, err := asynctask.GetDbImagePullReqJson(creatorId, req.ImageName, req.ImageTag, req.Type, req.Port)
	if err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "create async task error")
	}
	asyncTask := &models.AsyncTask{
		Id:       utils.GetSnowFlakeIdAndBase64(),
		UserId:   creatorId,
		Type:     asynctask.DbImagePullType,
		Args:     args,
		Priority: 0,
		Status:   models.AsyncTaskIng,
		CreateAt: time.Now().UnixMilli(),
	}
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.AsyncTaskDocument).InsertOne(l.ctx, asyncTask)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "insert async task error"))
		// 释放分布式锁
		if err = l.svcCtx.RedisClient.Del(l.ctx, fmt.Sprintf(l.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag)).Err(); err != nil {
			logx.Error(errors.Wrap(err, fmt.Sprintf(l.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag)+", 删除失败请即使清除避免系统异常"))
			l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(fmt.Sprintf(l.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag), "需要抓紧手动删除该分布式锁"))
		}
		return &types.CommonResponse{Code: 500, Msg: "insert async task error"}, nil
	}

	// 返回结果
	return &types.CommonResponse{Code: 200, Msg: "创建拉取镜像任务成功, 等待执行"}, nil
}
