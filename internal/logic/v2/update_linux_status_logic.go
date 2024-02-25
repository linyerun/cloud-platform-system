package v2

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"context"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLinuxStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLinuxStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLinuxStatusLogic {
	return &UpdateLinuxStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLinuxStatusLogic) UpdateLinuxStatus(req *types.UpdateLinuxStatusReq) (resp *types.CommonResponse, err error) {
	// 根据条件查找容器
	doc := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxContainerDocument)
	result := doc.FindOne(l.ctx, bson.D{{"_id", req.ContainerId}})
	if err = result.Err(); err != nil {
		l.Logger.Error(errors.Wrap(err, "get LinuxContainerDocument data err"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	container := new(models.LinuxContainer)
	if err = result.Decode(container); err != nil {
		l.Logger.Error(errors.Wrap(err, "decode LinuxContainerDocument data err"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 根据容器状态判断是否处于相反状态
	if req.Status == container.Status || (req.Status != models.LinuxSleep && req.Status != models.LinuxRunning) {
		return &types.CommonResponse{Code: 400, Msg: "status参数有误"}, nil
	}

	// 修改容器状态
	var updateData bson.M
	switch req.Status {
	case models.LinuxSleep:
		if err = l.svcCtx.DockerClient.StopContainerWithContext(container.Name, 0, l.ctx); err != nil {
			l.Logger.Error(errors.Wrap(err, "停止docker容器失败"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常，修改容器状态失败"}, nil
		}
		updateData = bson.M{"status": req.Status, "stop_time": time.Now().UnixMilli()}
	case models.LinuxRunning:
		if err = l.svcCtx.DockerClient.StartContainerWithContext(container.Name, &docker.HostConfig{}, l.ctx); err != nil {
			l.Logger.Error(errors.Wrap(err, "启动docker容器失败"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常，修改容器状态失败"}, nil
		}
		updateData = bson.M{"status": req.Status, "start_time": time.Now().UnixMilli()}
	}

	_, err = doc.UpdateByID(l.ctx, req.ContainerId, bson.D{{"$set", updateData}})
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "启动docker容器失败"))
		l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]any{"_id": container.Id, "set_data": updateData}, "手动修改linux_containers文档数据"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常，修改容器状态失败"}, nil
	}

	return &types.CommonResponse{Code: 200, Msg: "修改成功"}, nil
}
