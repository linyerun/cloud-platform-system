package asynctask

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/utils"
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const ImagePullType = "image_pull"

type ImagePullArgs struct {
	ImageName            string   `bson:"image_name"`
	ImageTag             string   `bson:"image_tag"`
	UserId               string   `bson:"user_id"`
	ImageEnabledCommands []string `bson:"image_enabled_commands"`
	ImageMustExportPorts []int64  `bson:"image_must_export_ports"`
}

type ImagePullHandler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImagePullHandler(ctx context.Context, srvCtx *svc.ServiceContext) IAsyncTaskHandler {
	return &ImagePullHandler{ctx: ctx, svcCtx: srvCtx}
}

func (i *ImagePullHandler) Execute(args any) (respData *RespData, status uint) {
	// 获取参数
	req := args.(*ImagePullArgs)

	// 异步任务被执行, 归还分布式锁
	defer func() {
		err := i.svcCtx.RedisClient.Del(i.ctx, fmt.Sprintf(i.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag)).Err()
		if err != nil {
			logx.Error(errors.Wrap(err, fmt.Sprintf(i.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag)+", 删除失败请即使清除避免系统异常"))
			i.svcCtx.RedisClient.RPush(i.ctx, i.svcCtx.ExceptionList, common.NewJsonMsgString(fmt.Sprintf(i.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag), "需要抓紧手动删除该分布式锁"))
		}
	}()

	// 执行拉取操作
	err := i.svcCtx.DockerClient.PullImage(docker.PullImageOptions{Repository: req.ImageName, Tag: req.ImageTag}, docker.AuthConfiguration{})
	if err != nil {
		logx.Error(errors.Wrap(err, "pull image error"))
		return &RespData{Code: 500, Msg: "拉取镜像异常"}, models.AsyncTaskFail
	}

	// 将镜像信息保存到mongo中
	dockerImage, err := i.svcCtx.DockerClient.InspectImage(req.ImageName + ":" + req.ImageTag)
	if err != nil {
		// 删除镜像，因为拉取失败了
		if err = i.svcCtx.DockerClient.RemoveImage(req.ImageName + ":" + req.ImageTag); err != nil {
			logx.Error(errors.Wrap(err, "image remove error"))
			i.svcCtx.RedisClient.RPush(i.ctx, i.svcCtx.ExceptionList, common.NewJsonMsgString(req.ImageName+":"+req.ImageTag, "需要删除这个镜像"))
		}
		return &RespData{Code: 500, Msg: "系统异常"}, models.AsyncTaskFail
	}
	image := &models.LinuxImage{
		Id:              utils.GetSnowFlakeIdAndBase64(),
		CreatorId:       req.UserId,
		Name:            req.ImageName,
		Tag:             req.ImageTag,
		ImageId:         dockerImage.ID,
		Size:            dockerImage.Size,
		EnableCommands:  req.ImageEnabledCommands,
		MustExportPorts: req.ImageMustExportPorts,
	}
	_, err = i.svcCtx.MongoClient.Database(i.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).InsertOne(i.ctx, image)
	if err != nil {
		logx.Error(errors.Wrap(err, "save image msg to mongo error"))
		if err = i.svcCtx.DockerClient.RemoveImage(req.ImageName + ":" + req.ImageTag); err != nil {
			logx.Error(errors.Wrap(err, "image remove error"))
			i.svcCtx.RedisClient.RPush(i.ctx, i.svcCtx.ExceptionList, common.NewJsonMsgString(req.ImageName+":"+req.ImageTag, "需要删除这个镜像"))
		}
		return &RespData{Code: 500, Msg: "系统异常"}, models.AsyncTaskFail
	}

	// 返回结果
	return &RespData{Code: 200, Msg: "拉取镜像成功"}, models.AsyncTaskOk
}
