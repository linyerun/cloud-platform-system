package v3

import (
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDbImageByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDbImageByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDbImageByIdLogic {
	return &GetDbImageByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDbImageByIdLogic) GetDbImageById(req *types.GetDbImageByIdReq) (resp *types.CommonResponse, err error) {
	ret := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbImageDocument).FindOne(l.ctx, bson.D{{"_id", req.Id}})
	if err = ret.Err(); err == mongo.ErrNoDocuments {
		return &types.CommonResponse{Code: 400, Msg: "id有误获取失败"}, nil
	} else if err == nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "get data error")
	}

	db := new(models.DbImage)
	if err = ret.Decode(db); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "decode data error")
	}

	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"db_image": db}}, nil
}
