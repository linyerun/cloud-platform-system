package v3

import (
	"bytes"
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/utils"
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os/exec"
	"strings"
	"time"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeDbApplicationStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeDbApplicationStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeDbApplicationStatusLogic {
	return &ChangeDbApplicationStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeDbApplicationStatusLogic) ChangeDbApplicationStatus(req *types.ChangeDbApplicationStatusReq) (resp *types.CommonResponse, err error) {
	// 获取application数据
	ret := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbApplicationFormDocument).FindOne(l.ctx, bson.D{{"_id", req.Id}})
	if err = ret.Err(); err == mongo.ErrNoDocuments {
		return nil, errorx.NewCodeError(400, "this application no exists")
	} else if err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "get application msg error")
	}
	// decode msg
	form := new(models.DbApplicationForm)
	if err = ret.Decode(form); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "decode msg error")
	}

	// 判断该审核是否可以进行
	if form.Status != models.DbApplicationFormStatusIng {
		return nil, errorx.NewCodeError(400, "this application has been checked!")
	}

	switch req.Status {
	case models.DbApplicationFormStatusOk:
		// 获取DbImage信息
		ret = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbImageDocument).FindOne(l.ctx, bson.D{{"_id", form.ImageId}})
		if err = ret.Err(); err == mongo.ErrNoDocuments {
			return nil, errorx.NewCodeError(400, "no db image")
		} else if err != nil {
			l.Logger.Error(err)
			return nil, errorx.NewCodeError(500, "get db image error")
		}
		image := new(models.DbImage)
		if err = ret.Decode(image); err != nil {
			l.Logger.Error(err)
			return nil, errorx.NewCodeError(500, "decode db image error")
		}

		// 获取端口
		var port uint
		port, err = l.svcCtx.PortManager.GetSinglePort()
		if err != nil {
			l.Error(err)
			return nil, errorx.NewCodeError(500, "get port from host error")
		}

		// 获取启动容器指令
		dbcName := image.Type + "-" + utils.GetSnowFlakeIdAndBase64()
		command, ok := l.getDbCmd(image.Type, image.ImageId, dbcName, image.Password, uint(image.Port), port)
		if ok {
			l.Logger.Error("image type error, image_id is " + image.Id)
			return nil, errorx.NewCodeError(500, "db image error")
		}

		// 运行启动指令
		cmd := exec.Command("docker", command...)
		var outputBuf, errorBuf bytes.Buffer
		cmd.Stdout = &outputBuf
		cmd.Stderr = &errorBuf
		l.Logger.Info("docker", command) // 记录执行的docker指令
		if err = cmd.Run(); err != nil {
			// 归还端口
			l.repayPort(port)
			// 记录错误日志
			l.Logger.Error(err)
			return nil, errorx.NewCodeError(500, "run db container error")
		}

		// 保存容器信息
		dbContainer := &models.DbContainer{
			Id:              utils.GetSnowFlakeIdAndBase64(),
			UserId:          form.UserId,
			Name:            dbcName,
			DbContainerName: form.DbName,
			ImageId:         image.ImageId,

			CreateAt:  time.Now().UnixMilli(),
			StartTime: time.Now().UnixMilli(),
			StopTime:  0,
			Status:    models.DbContainerStatusRunning,

			Host: l.svcCtx.Config.Container.Host,
			Port: port,

			Type:     image.Type,
			Username: "root",
			Password: image.Password,
		}
		_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbContainerDocument).InsertOne(l.ctx, dbContainer)
		if err != nil {
			// 归还端口
			l.repayPort(port)
			// 删除容器
			l.removeDockerContainer(dbcName)
			// 记录异常
			l.Logger.Error(err)
			return nil, errorx.NewCodeError(500, "start db error")
		}

		// 修改容器申请单状态
		filter := bson.D{{"_id", form.Id}}
		update := bson.D{{"$set", bson.M{"status": models.DbApplicationFormStatusOk, "finish_at": time.Now().UnixMilli()}}}
		_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbApplicationFormDocument).UpdateOne(l.ctx, filter, update)
		if err != nil {
			l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]any{"table": models.DbContainerDocument, "_id": req.Id, "updated": update}, "手动修改相关数据(可能被处理多次, 用户被新增多个Db)"))
			l.Logger.Error(err)
			return nil, errorx.NewCodeError(500, "系统异常, 修改容器状态失败")
		}
	case models.LinuxApplicationFormStatusReject:
		filter := bson.D{{"_id", req.Id}}
		update := bson.D{{"$set", bson.M{"status": models.DbApplicationFormStatusReject, "finish_at": time.Now().UnixMilli(), "reject_reason": req.RejectReason}}}
		_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbApplicationFormDocument).UpdateOne(l.ctx, filter, update)
		if err != nil {
			l.Logger.Error(err)
			return nil, errorx.NewCodeError(500, "拒绝申请失败")
		}
	default:
		return &types.CommonResponse{Code: 400, Msg: "status参数有误"}, nil
	}

	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}

func (l *ChangeDbApplicationStatusLogic) getDbCmd(dbType, imageId, dbName, password string, port1, port2 uint) ([]string, bool) {
	switch dbType {
	case models.DbImageTypeMySql:
		str := fmt.Sprintf("run --privileged=true --restart unless-stopped --name %s -d -e MYSQL_ROOT_PASSWORD=%s -p %d:%d %s", dbName, password, port1, port2, imageId)
		return strings.Split(str, " "), true
	case models.DbImageTypeRedis:
		str := fmt.Sprintf("run --privileged=true --restart unless-stopped --name %s -d -p %d:%d %s --requirepass %s", dbName, port1, port2, imageId, password)
		return strings.Split(str, " "), true
	case models.DbImageTypeMongoDb:
		str := fmt.Sprintf("run --privileged=true --restart unless-stopped --name %s -d -p %d:%d -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=%s %s", dbName, port1, port2, password, imageId)
		return strings.Split(str, " "), true
	}
	return nil, false
}

func (l *ChangeDbApplicationStatusLogic) repayPort(port uint) {
	if err := l.svcCtx.PortManager.RepaySinglePort(port); err != nil {
		l.Logger.Error(err)
		l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]uint{"from": port, "to": port}, "手动把from到to保存到port_recycle中"))
	}
}

func (l *ChangeDbApplicationStatusLogic) removeDockerContainer(dbcName string) {
	// 删除docker中指定的容器(docker相关)
	err := l.svcCtx.DockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: dbcName, Force: true})
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "移除docker容器失败, 需要手动删除并手动归还端口到mongo对应的文档中。container name: "+dbcName))
		l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(dbcName, "在docker中移除该名字的容器"))
	}
}
