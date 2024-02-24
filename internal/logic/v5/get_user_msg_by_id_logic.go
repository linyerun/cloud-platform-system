package v5

import (
	"cloud-platform-system/internal/models"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserMsgByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserMsgByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMsgByIdLogic {
	return &GetUserMsgByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserMsgByIdLogic) GetUserMsgById(req *types.GetUserMsgByIdReq) (resp *types.CommonResponse, err error) {
	result := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).FindOne(l.ctx, bson.D{{"_id", req.UserId}, {"status", models.UserApplicationFormStatusOk}}, options.FindOne().SetProjection(bson.D{{"password", 0}}))

	// 处理错误
	if err = result.Err(); err == mongo.ErrNoDocuments {
		return &types.CommonResponse{Code: 200, Msg: "查无此人", Data: map[string]any{"user": nil}}, nil
	} else if err != nil {
		l.Logger.Error(errors.Wrap(err, "find models.UserDocument error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// decode
	user := new(models.User)
	if err = result.Decode(user); err != nil {
		l.Logger.Error(errors.Wrap(err, "decode models.UserDocument error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 返回结果
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"user": user}}, nil
}
