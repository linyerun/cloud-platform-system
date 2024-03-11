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
	// get db image
	ret := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbImageDocument).FindOne(l.ctx, bson.D{{"_id", req.Id}})
	if err = ret.Err(); err == mongo.ErrNoDocuments {
		return &types.CommonResponse{Code: 400, Msg: "id有误获取失败"}, nil
	} else if err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewBaseError(500, "get data error")
	}

	db := new(models.DbImage)
	if err = ret.Decode(db); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewBaseError(500, "decode data error")
	}

	// get creator
	ret = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).FindOne(l.ctx, bson.D{{"_id", db.CreatorId}})
	if err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewBaseError(500, "get data error")
	}

	user := new(models.User)
	if err = ret.Decode(user); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewBaseError(500, "decode data error")
	}

	// return data
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"db_image": map[string]any{"db_image": db, "user_msg": map[string]any{"email": user.Email, "name": user.Name}}}}, nil
}
