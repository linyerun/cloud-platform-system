package v3

import (
	"cloud-platform-system/internal/asynctask"
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullLinuxImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullLinuxImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullLinuxImageLogic {
	return &PullLinuxImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullLinuxImageLogic) PullLinuxImage(req *types.ImagePullRequest) (resp *types.CommonResponse, err error) {
	// 必须指明版本
	if req.ImageTag == "latest" {
		return &types.CommonResponse{Code: 400, Msg: "必须只能具体版本，不能使用latest作为版本号"}, nil
	}

	// 端口必须在[0, 65535]这个范围
	for _, port := range req.ImageMustExportPorts {
		if port < 0 || port > 65535 {
			return &types.CommonResponse{Code: 400, Msg: "每个端口范围必须在[0, 65535]里面"}, nil
		}
	}

	// 避免重复拉取镜像
	err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).FindOne(l.ctx, bson.D{{"name", req.ImageName}, {"tag", req.ImageTag}}).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		l.Logger.Error(errors.Wrap(err, "find data in mongo error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	} else if err == nil {
		return &types.CommonResponse{Code: 400, Msg: "镜像已经存在，无需拉取！"}, nil
	}

	// 使用锁来保证同一个镜像只能有一个管理员创建拉取任务
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

	// 创建异步任务
	asyncTask := models.AsyncTask{
		Id:     utils.GetSnowFlakeIdAndBase64(),
		UserId: l.ctx.Value("user").(*models.User).Id,
		Type:   asynctask.ImagePullType,
		Args: (&asynctask.ImagePullArgs{
			ImageName:            req.ImageName,
			ImageTag:             req.ImageTag,
			UserId:               l.ctx.Value("user").(*models.User).Id,
			ImageEnabledCommands: req.ImageEnabledCommands,
			ImageMustExportPorts: req.ImageMustExportPorts,
		}).JsonMarshal(),
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
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 返回结果
	return &types.CommonResponse{Code: 200, Msg: "创建拉取镜像任务成功，等待执行"}, nil
}
