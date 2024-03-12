package v0

import (
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeForgetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeForgetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeForgetPasswordLogic {
	return &ChangeForgetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeForgetPasswordLogic) ChangeForgetPassword(req *types.ChangeForgetPasswordReq) (resp *types.CommonResponse, err error) {
	err = utils.IsValidEmailCaptcha(l.svcCtx.RedisClient, l.svcCtx.CAPTCHA, req.Captcha, req.Email)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewBaseError(400, "captcha err"), "err: %v, req: %+v", err, req)
	}

	if len(req.Password) < 6 {
		return nil, errorx.NewBaseError(400, "password长度不可小于6")
	}

	// 执行修改操作
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).UpdateOne(l.ctx, bson.D{{"email", req.Email}}, bson.D{{"$set", bson.M{"password": utils.DoHashAndBase64(l.svcCtx.Config.Salt, req.Password)}}})
	if err != nil {
		return nil, errors.Wrapf(errorx.NewBaseError(500, "update user pwd err"), "err: %v, req: %+v", err, req)
	}

	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}
