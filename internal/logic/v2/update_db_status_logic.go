package v2

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDbStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDbStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDbStatusLogic {
	return &UpdateDbStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDbStatusLogic) UpdateDbStatus(req *types.UpdateDbStatusReq) (resp *types.CommonResponse, err error) {
	// 根据db_id查找数据库容器
	ret := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbContainerDocument).FindOne(l.ctx, bson.D{{"_id", req.DbId}})
	if err = ret.Err(); err == mongo.ErrNoDocuments {
		return nil, errorx.NewBaseError(400, "this db container not exists")
	} else if err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewBaseError(500, "get db msg error")
	}

	//decode
	db := new(models.DbContainer)
	if err = ret.Decode(db); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewBaseError(500, "decode db msg error")
	}

	// 判断status是否符合要求
	if req.Status == db.Status {
		return nil, errorx.NewBaseError(400, "当前数据库正处于此状态")
	}

	var updateData bson.M
	var rollbackFunc func() error
	switch req.Status {
	case models.DbContainerStatusSleeping:
		if err = l.svcCtx.DockerClient.StopContainerWithContext(db.Name, 0, l.ctx); err != nil {
			l.Logger.Error(errors.Wrap(err, "停止db容器失败,name: "+db.Name))
			return &types.CommonResponse{Code: 500, Msg: "停止db容器失败"}, nil
		}
		rollbackFunc = func() error {
			return l.svcCtx.DockerClient.StartContainerWithContext(db.Name, &docker.HostConfig{}, l.ctx)
		}
		updateData = bson.M{"status": req.Status, "stop_time": time.Now().UnixMilli()}
	case models.DbContainerStatusRunning:
		if err = l.svcCtx.DockerClient.StartContainerWithContext(db.Name, &docker.HostConfig{}, l.ctx); err != nil {
			l.Logger.Error(errors.Wrap(err, "启动db容器失败,name: "+db.Name))
			return &types.CommonResponse{Code: 500, Msg: "启动db容器失败"}, nil
		}
		rollbackFunc = func() error { return l.svcCtx.DockerClient.StopContainerWithContext(db.Name, 0, l.ctx) }
		updateData = bson.M{"status": req.Status, "start_time": time.Now().UnixMilli()}
	case models.DbContainerDel:
		// 删除指定文档(mongo相关)
		err = common.DelData(l.Logger, l.svcCtx, func() (any, error) {
			_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbContainerDocument).DeleteOne(l.ctx, bson.D{{"_id", db.Id}})
			if err != nil {
				l.Logger.Error(errors.Wrap(err, "delete DbContainerDocument data err"))
				return nil, err
			}
			return db, nil
		}, models.DbContainerDocument)
		if err != nil && err != common.SaveMongoDelDataError {
			l.Logger.Error(errors.Wrap(err, "delete data error"))
			return &types.CommonResponse{Code: 500, Msg: "delete data error"}, nil
		}

		// 归还端口(redis相关)
		if err = l.svcCtx.PortManager.RepaySinglePort(db.Port); err != nil {
			l.Logger.Error(err)
			l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]any{"_id": fmt.Sprintf("%d-%d", db.Port, db.Port), "table": models.PortRecycleDocument, "type": models.SinglePortType}, "手动把单个端口加入复用表"))
		}

		// 删除docker中指定的容器(docker相关)
		err = l.svcCtx.DockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: db.Name, Force: true})
		if err != nil {
			l.Logger.Error(errors.Wrap(err, "移除docker容器失败, 需要手动删除并手动归还端口到mongo对应的文档中。container name: "+db.Name))
			l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(db.Name, "在docker中移除该名字的容器"))
		}

		// 响应结果(其实只是逻辑删除)
		return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
	default:
		return nil, errorx.NewBaseError(400, "status参数有误")
	}

	// 执行更新db操作
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbContainerDocument).UpdateByID(l.ctx, db.Id, bson.D{{"$set", updateData}})
	if err != nil {
		l.Logger.Error(err)
		if err = rollbackFunc(); err != nil { // 可以回滚成功那就最好, 不行那也没办法了
			l.Logger.Error(err)
			l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]any{"_id": db.Id, "set_data": updateData}, "手动修改db_containers文档数据"))
		}
		return nil, errorx.NewBaseError(500, "系统异常")
	}

	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}
