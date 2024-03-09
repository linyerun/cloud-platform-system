package asynctask

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/utils"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	PwdStr          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	DbImagePullType = "db_image_pull_type"
)

type DbImagePullReq struct {
	CreatorId string `json:"creator_id"`
	ImageName string `json:"image_name"`
	ImageTag  string `json:"image_tag"`
	Type      string `json:"type"`
	Port      uint   `json:"port"`
}

func GetDbImagePullReqJson(cId, imageName, imageTag, imageType string, Port uint) (string, error) {
	obj := &DbImagePullReq{
		CreatorId: cId,
		ImageName: imageName,
		ImageTag:  imageTag,
		Type:      imageType,
		Port:      Port,
	}
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func GetDbImagePullReqByJsonStr(str string) (*DbImagePullReq, error) {
	obj := new(DbImagePullReq)
	err := json.Unmarshal([]byte(str), obj)
	return obj, err
}

type DbImagePullHandler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDbImagePullHandler(ctx context.Context, srvCtx *svc.ServiceContext) IAsyncTaskHandler {
	return &DbImagePullHandler{ctx: ctx, svcCtx: srvCtx}
}

func (i *DbImagePullHandler) Execute(args string) (respData *RespData, status uint) {
	// decode msg
	req, err := GetDbImagePullReqByJsonStr(args)
	if err != nil {
		logx.Error("decode DbImagePullReq Json Msg error")
		return &RespData{Code: 500, Msg: "decode DbImagePullReq Json Msg error"}, models.AsyncTaskFail
	}

	// 异步任务被执行, 归还分布式锁
	defer func() {
		err = i.svcCtx.RedisClient.Del(i.ctx, fmt.Sprintf(i.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag)).Err()
		if err != nil {
			logx.Error(errors.Wrap(err, fmt.Sprintf(i.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag)+", 删除失败请即使清除避免系统异常"))
			i.svcCtx.RedisClient.RPush(i.ctx, i.svcCtx.ExceptionList, common.NewJsonMsgString(fmt.Sprintf(i.svcCtx.ImagePrefix, req.ImageName+":"+req.ImageTag), "需要抓紧手动删除该分布式锁"))
		}
	}()

	// 使用docker拉取镜像
	err = i.svcCtx.DockerClient.PullImage(docker.PullImageOptions{Repository: req.ImageName, Tag: req.ImageTag}, docker.AuthConfiguration{})
	if err != nil {
		logx.Error(errors.Wrap(err, "pull image error"))
		return &RespData{Code: 500, Msg: "拉取镜像异常"}, models.AsyncTaskFail
	}

	// 获取镜像信息
	dockerImage, err := i.svcCtx.DockerClient.InspectImage(req.ImageName + ":" + req.ImageTag)
	if err != nil {
		logx.Error("从docker中获取镜像信息失败, image_name: " + req.ImageName + ":" + req.ImageTag)
		return &RespData{Code: 500, Msg: "从docker中获取镜像信息失败, image_name: " + req.ImageName + ":" + req.ImageTag}, models.AsyncTaskFail
	}

	// 将镜像信息保存到DbImage
	image := &models.DbImage{
		Id:        utils.GetSnowFlakeIdAndBase64(),
		CreatorId: req.CreatorId,

		Type: req.Type,

		Name: req.ImageName,
		Tag:  req.ImageTag,

		Username: models.DbUsername,
		Password: getPwd(10),

		ImageId: dockerImage.ID,
		Size:    dockerImage.Size,

		Port: req.Port,

		CreatedAt: time.Now().UnixMilli(),
		IsDeleted: false,
	}
	_, err = i.svcCtx.MongoClient.Database(i.svcCtx.Config.Mongo.DbName).Collection(models.DbImageDocument).InsertOne(i.ctx, image)
	if err != nil {
		logx.Error(err)
		return &RespData{Code: 500, Msg: "Save DbImageError"}, models.AsyncTaskFail
	}

	return &RespData{Code: 200, Msg: "拉取镜像成功"}, models.AsyncTaskOk
}

func getPwd(pwdLen int) string {
	if pwdLen < 6 {
		pwdLen = 6
	}
	var pwd []byte
	for i := 0; i < pwdLen; i++ {
		pwd = append(pwd, PwdStr[rand.Intn(len(PwdStr))])
	}
	return string(pwd)
}
