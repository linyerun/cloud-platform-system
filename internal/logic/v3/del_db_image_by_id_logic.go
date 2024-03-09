package v3

import (
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelDbImageByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelDbImageByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelDbImageByIdLogic {
	return &DelDbImageByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelDbImageByIdLogic) DelDbImageById(req *types.DelDbImageByIdReq) (resp *types.CommonResponse, err error) {
	update := bson.D{{"$set", bson.M{"is_deleted": true}}}
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbImageDocument).UpdateOne(l.ctx, bson.D{{"_id", req.Id}}, update)
	if err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "del image error")
	}
	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}
