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
	"time"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

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

	// 避免重复拉取镜像
	err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).FindOne(l.ctx, bson.D{{"name", req.ImageName}, {"tag", req.ImageTag}}).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		l.Logger.Error(errors.Wrap(err, "find data in mongo error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	} else if err == nil {
		return &types.CommonResponse{Code: 400, Msg: "镜像已经存在，无需拉取！"}, nil
	}

	// 使用锁来保证同一个镜像只能有一个管理员拉取
	res := l.svcCtx.RedisClient.SetNX(l.ctx, fmt.Sprintf(l.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag), "1", time.Minute*30)
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
		return &types.CommonResponse{Code: 401, Msg: "已有管理员拉取此镜像，无所重复拉取镜像！"}, nil
	}
	// 归还分布式锁
	defer l.svcCtx.RedisClient.Del(l.ctx, fmt.Sprintf(l.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag))

	// 执行拉取操作
	err = l.svcCtx.DockerClient.PullImage(docker.PullImageOptions{Repository: req.ImageName, Tag: req.ImageTag}, docker.AuthConfiguration{})
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "pull image error"))
		return &types.CommonResponse{Code: 500, Msg: "拉取镜像异常"}, nil
	}

	// 将镜像信息保存到mongo中
	dockerImage, err := l.svcCtx.DockerClient.InspectImage(req.ImageName + ":" + req.ImageTag)
	if err != nil {
		err = l.svcCtx.DockerClient.RemoveImage(req.ImageName + ":" + req.ImageTag)
		if err != nil {
			l.Logger.Error(errors.Wrap(err, "image remove error"))
			l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(req.ImageName+":"+req.ImageTag, "需要删除这个镜像"))
		}
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	image := &models.LinuxImage{
		Id:        utils.GetSnowFlakeIdAndBase64(),
		CreatorId: l.ctx.Value("user").(*models.User).Id,
		Name:      req.ImageName,
		Tag:       req.ImageTag,
		ImageId:   dockerImage.ID,
		Size:      dockerImage.Size,
	}
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).InsertOne(l.ctx, image)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "save image msg to mongo error"))
		if err = l.svcCtx.DockerClient.RemoveImage(req.ImageName + ":" + req.ImageTag); err != nil {
			l.Logger.Error(errors.Wrap(err, "image remove error"))
			l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(req.ImageName+":"+req.ImageTag, "需要删除这个镜像"))
		}
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 返回结果
	return &types.CommonResponse{Code: 200, Msg: "拉取镜像成功"}, nil
}
