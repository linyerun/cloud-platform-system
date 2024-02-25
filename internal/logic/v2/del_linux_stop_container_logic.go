package v2

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"context"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"sort"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelLinuxStopContainerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelLinuxStopContainerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelLinuxStopContainerLogic {
	return &DelLinuxStopContainerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelLinuxStopContainerLogic) DelLinuxStopContainer(req *types.DelLinuxStopContainerReq) (resp *types.CommonResponse, err error) {
	linux := new(models.LinuxContainer)

	// 删除指定容器数据，并把数据保存到删除表中
	err = common.DelData(l.Logger, l.svcCtx, func() (any, error) {
		result := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxContainerDocument).FindOneAndDelete(l.ctx, bson.D{{"_id", req.ContainerId}})
		if err = result.Err(); err != nil {
			l.Logger.Error(errors.Wrap(err, "get and delete LinuxContainerDocument data err"))
			return nil, err
		}
		if err = result.Decode(linux); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode LinuxContainerDocument data err"))
			return nil, err
		}
		return linux, nil
	}, models.LinuxContainerDocument)
	if err != nil && err != common.SaveMongoDelDataError {
		l.Logger.Error(errors.Wrap(err, "delete data error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 删除docker中对应的容器(如果失败了，需要管理员手动操作)
	err = l.svcCtx.DockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: linux.Name, Force: true})
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "移除docker容器失败, 需要手动删除并手动归还端口到mongo对应的文档中。container name: "+linux.Name))
		l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(linux.Name, "在docker中移除该名字的容器"))
		// 响应结果(其实只是逻辑删除)
		return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
	}

	// 归还端口
	portsMapping := linux.PortsMapping
	var ports []int
	for k := range portsMapping {
		ports = append(ports, int(k))
	}
	sort.Ints(ports)
	if err = l.svcCtx.PortManager.RepayTenPort(uint(ports[0]), uint(ports[len(ports)-1])); err != nil {
		l.Logger.Error(err)
		l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]uint{"start_port": uint(ports[0]), "stop_port": uint(ports[len(ports)-1])}, "把端口手动写入port_recycle"))
	}

	// 响应结果
	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}
