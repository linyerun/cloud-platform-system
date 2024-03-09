package common

import (
	"context"

	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetDbImageList(svcCtx *svc.ServiceContext, Logger logx.Logger, ctx context.Context) (*types.CommonResponse, error) {
	cur, err := svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.DbImageDocument).Find(ctx, bson.D{{"is_deleted", false}})
	if err == mongo.ErrNoDocuments {
		return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
	} else if err == nil {
		Logger.Error(err)
		return nil, errorx.NewCodeError(500, "get data error")
	}

	var dbs []*models.DbImage
	for cur.Next(ctx) {
		db := new(models.DbImage)
		if err = cur.Decode(db); err != nil {
			Logger.Error(err)
			return nil, errorx.NewCodeError(500, "decode data error")
		}
		dbs = append(dbs, db)
	}

	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"db_images": dbs}}, nil
}
