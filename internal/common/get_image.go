package common

import (
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllImage(svcCtx *svc.ServiceContext, Logger logx.Logger, ctx context.Context) (*types.CommonResponse, error) {
	// 获取镜像
	cur, err := svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.ImageDocument).Find(ctx, bson.D{})
	if err != nil {
		Logger.Error(errors.Wrap(err, "find images error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	var creators = make(map[string]*models.User) // key: admin_id, value: admin
	type ImageDto struct {
		Id      string       `json:"id"`
		Creator *models.User `json:"creator"`
		Name    string       `json:"name"`
		Tag     string       `json:"tag"`
		ImageId string       `json:"image_id"`
		Size    int64        `json:"size"`
	}
	var respData []*ImageDto

	for cur.Next(ctx) {
		// 解析image数据
		image := new(models.Image)
		if err = cur.Decode(image); err != nil {
			Logger.Error(errors.Wrap(err, "decode images error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
		// 组装creator信息
		if _, ok := creators[image.CreatorId]; !ok {
			result := svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).FindOne(ctx, bson.D{{"_id", image.CreatorId}}, options.FindOne().SetProjection(bson.D{{"password", 0}, {"auth", 0}}))
			if err = result.Err(); err != nil {
				Logger.Error(errors.Wrap(err, "get user error"))
				return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
			}
			user := new(models.User)
			if err = result.Decode(user); err != nil {
				Logger.Error(errors.Wrap(err, "decode user error"))
				return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
			}
			creators[user.Id] = user
		}
		respData = append(respData, &ImageDto{Id: image.Id, ImageId: image.ImageId, Name: image.Name, Tag: image.Tag, Size: image.Size, Creator: creators[image.CreatorId]})
	}

	// 返回结果
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"images": respData}}, nil
}
