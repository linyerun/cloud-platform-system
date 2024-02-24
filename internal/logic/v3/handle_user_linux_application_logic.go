package v3

import (
	"cloud-platform-system/internal/common"
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
	switch req.Status {
	case models.LinuxApplicationFormStatusOk:
		// 获取申请单并修改申请单
		filter := bson.D{{"_id", req.FormId}}
		update := bson.D{{"$set", bson.M{"status": models.LinuxApplicationFormStatusOk, "finish_at": time.Now().UnixMilli()}}}
		updateResult := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).FindOneAndUpdate(l.ctx, filter, update)
		if err = updateResult.Err(); err != nil && err != mongo.ErrNoDocuments {
			l.Logger.Error(errors.Wrap(err, "update LinuxApplicationFormDocument error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		} else if err == mongo.ErrNoDocuments {
			return &types.CommonResponse{Code: 400, Msg: "不存在该Linux服务器申请单"}, nil
		}
		form := new(models.LinuxApplicationForm)
		if err = updateResult.Decode(form); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode LinuxApplicationFrom error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}

		// 启动容器时需要做的配置
		containerName := strings.ReplaceAll(utils.GetSnowFlakeIdAndBase64(), "=", ".")
		nameOption := utils.WithNameOption(containerName)
		coreCountOption := utils.WithCpuCoreCountOption(form.CoreCount)
		memoryOption := utils.WithMemoryOption(form.Memory)
		memorySwapOption := utils.WithMemorySwapOption(form.MemorySwap)
		diskSizeOption := utils.WithDiskSizeOption(form.DiskSize)

		// 根据image_id获取镜像信息
		findResult := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).FindOne(l.ctx, bson.D{{"_id", form.ImageId}})
		if err = findResult.Err(); err != nil {
			// 复原form
			rollback01(req, l)
			// 记录错误日志
			l.Logger.Error(errors.Wrap(err, "get image error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
		image := new(models.LinuxImage)
		if err = findResult.Decode(image); err != nil {
			// 复原form
			rollback01(req, l)
			// 记录错误日志
			l.Logger.Error(errors.Wrap(err, "decode image error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}

		// 构建容器启动首次需要执行的命令
		containerRunCommandOption := utils.WithImageAndContainerCommand(image.ImageId, image.EnableCommands)

		// 构建端口映射(只能映射10个端口)
		from, to, err := l.svcCtx.PortManager.GetTenPort()
		if err != nil {
			// 复原form
			rollback01(req, l)
			// 记录错误日志
			l.Logger.Error(errors.Wrap(err, "get ten port error"))
			return &types.CommonResponse{Code: 500, Msg: err.Error()}, nil
		}
		var portMappingOptions []utils.ContainerRunCommandOption
		var portsMapping = make(map[int64]int64)
		for i := from; i < to; i++ {
			if len(image.MustExportPorts) <= int(i-from) { // 超出范围的
				portMappingOptions = append(portMappingOptions, utils.WithPortMappingOption(int64(i), int64(i)))
				portsMapping[int64(i)] = int64(i)
				continue
			}
			portsMapping[int64(i)] = image.MustExportPorts[i-from]
			portMappingOptions = append(portMappingOptions, utils.WithPortMappingOption(int64(i), image.MustExportPorts[i-from]))
		}

		// 运行Linux容器
		commands := utils.CreateContainerRunCommand(append(portMappingOptions, nameOption, coreCountOption, memoryOption, memorySwapOption, diskSizeOption, containerRunCommandOption)...)
		l.Logger.Infof("%v", commands)
		output, err := exec.Command("docker", commands...).Output()
		if err != nil {
			// 复原form、归还端口
			rollback02(req, l, from, to)

			// 记录错误日志
			l.Logger.Error(errors.Wrap(err, "run LinuxContainer error: "+string(output)))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}

		// 将容器信息保存到表单中
		linux := models.LinuxContainer{
			Id:                utils.GetSnowFlakeIdAndBase64(),
			UserId:            form.UserId,
			Name:              containerName,
			UserContainerName: form.ContainerName,
			ImageId:           form.ImageId,

			CreateAt:  time.Now().UnixMilli(),
			StartTime: time.Now().UnixMilli(),
			Status:    models.LinuxSleep,

			Host:         l.svcCtx.Config.Container.Host,
			PortsMapping: portsMapping,

			InitUsername: l.svcCtx.Config.Container.InitUsername,
			InitPassword: l.svcCtx.Config.Container.InitPassword,

			Memory:     form.Memory,
			MemorySwap: form.MemorySwap,

			CoreCount: form.CoreCount,
			DiskSize:  form.DiskSize,
		}
		_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxContainerDocument).InsertOne(l.ctx, linux)
		if err != nil {
			// 回滚前面的修改
			rollback03(req, l, containerName, from, to)
			l.Logger.Error(errors.Wrap(err, ""))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}

		// 响应信息
		return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
	case models.LinuxApplicationFormStatusReject:
		filter := bson.D{{"_id", req.FormId}}
		update := bson.D{{"$set", bson.M{"status": models.LinuxApplicationFormStatusReject, "finish_at": time.Now().UnixMilli()}}}
		_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).UpdateOne(l.ctx, filter, update)
		if err != nil && err != mongo.ErrNoDocuments {
			l.Logger.Error(errors.Wrap(err, "update LinuxApplicationFormDocument error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		} else if err == mongo.ErrNoDocuments {
			return &types.CommonResponse{Code: 400, Msg: "不存在该Linux服务器申请单"}, nil
		}
	default:
		return &types.CommonResponse{Code: 400, Msg: "status存在问题"}, nil
	}
	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}

func rollback01(req *types.HandleUserLinuxApplicationReq, l *HandleUserLinuxApplicationLogic) {
	// 复原form
	filter := bson.D{{"_id", req.FormId}}
	update := bson.D{{"$set", bson.M{"status": models.LinuxApplicationFormStatusIng, "finish_at": 0}}}
	_, err := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).UpdateOne(l.ctx, filter, update)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "update LinuxApplicationFormDocument error"))
		l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(update, "把update的信息手动操作到mongo的"+models.LinuxApplicationFormDocument+"表中"))
	}
}

func rollback02(req *types.HandleUserLinuxApplicationReq, l *HandleUserLinuxApplicationLogic, from, to uint) {
	rollback01(req, l)

	// 归还端口
	err := l.svcCtx.PortManager.RepayTenPort(from, to)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, fmt.Sprintf("归还端口失败, 需要手动输入到port_recycle表. from=%d, to=%d", from, to)))
		l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]uint{"from": from, "to": to}, "手动把from到to保存到port_recycle中"))
	}
}

func rollback03(req *types.HandleUserLinuxApplicationReq, l *HandleUserLinuxApplicationLogic, containerName string, from, to uint) {
	rollback02(req, l, from, to)

	// 删除容器
	err := l.svcCtx.DockerClient.RemoveContainer(docker.RemoveContainerOptions{Context: l.ctx, ID: containerName, Force: true})
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "remove docker container error"))
		l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]string{"name": containerName}, "需要手动删除这个容器在docker中"))
	}
}
