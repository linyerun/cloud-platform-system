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
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type HandleUserLinuxApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleUserLinuxApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleUserLinuxApplicationLogic {
	return &HandleUserLinuxApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandleUserLinuxApplicationLogic) HandleUserLinuxApplication(req *types.HandleUserLinuxApplicationReq) (resp *types.CommonResponse, err error) {
	// 使用锁来保证只会提交一个异步处理的任务
	res := l.svcCtx.RedisClient.SetNX(l.ctx, fmt.Sprintf("linux_application_form_handler: %s", req.FormId), "1", redis.KeepTTL)
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
		return &types.CommonResponse{Code: 401, Msg: "容器审核异步处理中"}, nil
	}

	// 把异步任务保存到mongo中等待被执行
	asyncTask := &models.AsyncTask{
		Id:       utils.GetSnowFlakeIdAndBase64(),
		UserId:   l.ctx.Value("user").(*models.User).Id,
		Type:     asynctask.ContainerRunType,
		Args:     (&asynctask.ContainerRunArgs{FormId: req.FormId, Status: req.Status}).JsonMarshal(),
		Status:   models.AsyncTaskIng,
		CreateAt: time.Now().UnixMilli(),
	}
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.AsyncTaskDocument).InsertOne(l.ctx, asyncTask)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "insert async task error"))
		// 释放分布式锁
		if err = l.svcCtx.RedisClient.Del(l.ctx, fmt.Sprintf("linux_application_form_handler: %s", req.FormId)).Err(); err != nil {
			logx.Error(errors.Wrap(err, fmt.Sprintf("linux_application_form_handler: %s", req.FormId)+", 删除失败请即使清除避免系统异常"))
			l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(fmt.Sprintf("linux_application_form_handler: %s", req.FormId), "需要抓紧手动删除该分布式锁"))
		}
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	return &types.CommonResponse{Code: 200, Msg: "已提交容器审核，等待异步处理"}, nil
}
