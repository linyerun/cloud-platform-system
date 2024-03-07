package asynctask

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"os/exec"
	"strings"
	"time"
)

const ContainerRunType = "container_run"

type ContainerRunArgs struct {
	FormId string `json:"form_id"`
	Status uint   `json:"status"`
}

func (args *ContainerRunArgs) JsonMarshal() string {
	bytes, _ := json.Marshal(args)
	return string(bytes)
}

func (args *ContainerRunArgs) JsonUnmarshal(data string) {
	json.Unmarshal([]byte(data), args)
}

type ContainerRunArgsHandler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewContainerRunArgsHandler(ctx context.Context, srvCtx *svc.ServiceContext) IAsyncTaskHandler {
	return &ContainerRunArgsHandler{ctx: ctx, svcCtx: srvCtx}
}

func (l *ContainerRunArgsHandler) Execute(args string) (respData *RespData, status uint) {
	// 获取参数
	req := new(ContainerRunArgs)
	req.JsonUnmarshal(args)

	// 核心处理逻辑
	var err error

	// 获取form
	filter := bson.D{{"_id", req.FormId}}
	result := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).FindOne(l.ctx, filter)
	if err = result.Err(); err != nil {
		logx.Error(errors.Wrap(err, "get LinuxApplicationFormDocument error"))
		return &RespData{Code: 500, Msg: "系统异常"}, models.AsyncTaskFail
	}
	// decode
	form := new(models.LinuxApplicationForm)
	if err = result.Decode(form); err != nil {
		logx.Error(errors.Wrap(err, "decode LinuxApplicationFrom error"))
		return &RespData{Code: 500, Msg: "系统异常"}, models.AsyncTaskFail
	}

	switch req.Status {
	case models.LinuxApplicationFormStatusOk:
		var containerName string
		var from, to uint

		// 最后根据处理结果来修改订单状态
		defer func() {
			formStatus := models.LinuxApplicationFormStatusOk
			if status == models.AsyncTaskFail { // 对于执行期间出现异常的，直接拒绝申请处理
				formStatus = models.LinuxApplicationFormStatusReject
			}
			// 最后才修改form的状态
			update := bson.D{{"$set", bson.M{"status": formStatus, "finish_at": time.Now().UnixMilli()}}}
			_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).UpdateOne(l.ctx, filter, update)
			if err != nil {
				// 需要回滚删除很多东西了
				containerRunRollbackPortsContainer(l.ctx, l.svcCtx, containerName, from, to)
				logx.Error(errors.Wrap(err, "update LinuxApplicationFormDocument error"))
				respData = &RespData{Code: 500, Msg: "update LinuxApplicationFormDocument error"}
				status = models.AsyncTaskFail
			}
		}()

		// 启动容器时需要做的配置
		containerName = strings.ReplaceAll(utils.GetSnowFlakeIdAndBase64(), "=", ".") // =不能作为container name的符号但是.可以
		nameOption := utils.WithNameOption(containerName)
		coreCountOption := utils.WithCpuCoreCountOption(form.CoreCount)
		memoryOption := utils.WithMemoryOption(form.Memory)
		memorySwapOption := utils.WithMemorySwapOption(form.MemorySwap)
		diskSizeOption := utils.WithDiskSizeOption(form.DiskSize)

		// 根据image_id获取镜像信息
		findResult := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).FindOne(l.ctx, bson.D{{"_id", form.ImageId}})
		if err = findResult.Err(); err != nil {
			// 记录错误日志
			logx.Error(errors.Wrap(err, "get image error"))
			return &RespData{Code: 500, Msg: "get image error"}, models.AsyncTaskFail
		}
		image := new(models.LinuxImage)
		if err = findResult.Decode(image); err != nil {
			// 记录错误日志
			logx.Error(errors.Wrap(err, "decode image error"))
			return &RespData{Code: 500, Msg: "decode image error"}, models.AsyncTaskFail
		}
		// 把用户自定义开启的端口放入image.MustExportPorts中方便后面操作
		image.MustExportPorts = append(image.MustExportPorts, form.ExportPorts...)

		// 构建容器启动首次需要执行的命令
		containerRunCommandOption := utils.WithImageAndContainerCommand(image.ImageId, image.EnableCommands)

		// 构建端口映射(只能映射10个端口)
		from, to, err = l.svcCtx.PortManager.GetTenPort()
		if err != nil {
			// 记录错误日志
			logx.Error(errors.Wrap(err, "get ten port error"))
			return &RespData{Code: 500, Msg: "get ten port error. " + err.Error()}, models.AsyncTaskFail
		}
		var portMappingOptions []utils.ContainerRunCommandOption
		var portsMapping = make(map[int64]int64)
		for i := from; i <= to; i++ {
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
		logx.Infof("运行指令: %v", commands)
		output, err := exec.Command("docker", commands...).Output()
		if err != nil {
			// 归还端口、删除容器(没有得删那就不管了)
			containerRunRollbackPortsContainer(l.ctx, l.svcCtx, containerName, from, to)

			// 记录错误日志
			logx.Error(errors.Wrap(err, "run LinuxContainer error: "+string(output)))
			return &RespData{Code: 500, Msg: "系统异常"}, models.AsyncTaskFail
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
			Status:    models.LinuxRunning,

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
			containerRunRollbackPortsContainer(l.ctx, l.svcCtx, containerName, from, to)
			logx.Error(errors.Wrap(err, ""))
			return &RespData{Code: 500, Msg: "系统异常"}, models.AsyncTaskFail
		}

		// 响应信息
		return &RespData{Code: 200, Msg: "成功"}, models.AsyncTaskOk
	case models.LinuxApplicationFormStatusReject:
		filter = bson.D{{"_id", req.FormId}}
		update := bson.D{{"$set", bson.M{"status": models.LinuxApplicationFormStatusReject, "finish_at": time.Now().UnixMilli()}}}
		_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).UpdateOne(l.ctx, filter, update)
		if err != nil {
			logx.Error(errors.Wrap(err, "update LinuxApplicationFormDocument error"))
			return &RespData{Code: 500, Msg: "update LinuxApplicationFormDocument error"}, models.AsyncTaskFail
		}
	default:
		return &RespData{Code: 400, Msg: "status存在问题"}, models.AsyncTaskFail
	}
	return
}

func containerRunRollbackPorts(ctx context.Context, svcCtx *svc.ServiceContext, from, to uint) {
	// 归还端口
	err := svcCtx.PortManager.RepayTenPort(from, to)
	if err != nil {
		logx.Error(errors.Wrap(err, fmt.Sprintf("归还端口失败, 需要手动输入到port_recycle表. from=%d, to=%d", from, to)))
		svcCtx.RedisClient.RPush(ctx, svcCtx.ExceptionList, common.NewJsonMsgString(map[string]uint{"from": from, "to": to}, "手动把from到to保存到port_recycle中"))
	}
}

func containerRunRollbackPortsContainer(ctx context.Context, svcCtx *svc.ServiceContext, containerName string, from, to uint) {
	// 归还端口
	containerRunRollbackPorts(ctx, svcCtx, from, to)
	// 删除容器
	err := svcCtx.DockerClient.RemoveContainer(docker.RemoveContainerOptions{Context: ctx, ID: containerName, Force: true})
	if err != nil {
		logx.Error(errors.Wrap(err, "remove docker container error"))
		svcCtx.RedisClient.RPush(ctx, svcCtx.ExceptionList, common.NewJsonMsgString(map[string]string{"name": containerName}, "需要手动删除这个容器在docker中"))
	}
}
