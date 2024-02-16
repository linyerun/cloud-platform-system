package v0

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserLoginRequest) (resp *types.CommonResponse, err error) {
	// 校验参数
	if !utils.IsNormalEmail(req.Email) || len(req.Password) < 1 || !utils.IsValidCaptcha(l.svcCtx.RedisClient, l.svcCtx.CAPTCHA, req.Captcha) {
		return &types.CommonResponse{Code: 400, Msg: "参数有误"}, nil
	}

	// 查询具体信息
	var user = new(models.User)
	filter := bson.M{"email": req.Email, "password": utils.DoHashAndBase64(l.svcCtx.Config.Salt, req.Password)}
	err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserTable).FindOne(l.ctx, filter).Decode(user)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "find user error, email is "+req.Email))
		return &types.CommonResponse{Code: 500, Msg: err.Error()}, nil
	}

	// 把信息封装到token中并返回给用户
	return common.GetToken(user, l.Logger, l.svcCtx)
}
