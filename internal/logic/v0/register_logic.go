package v0

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.UserRegisterRequest) (resp *types.CommonResponse, err error) {
	// 校验参数
	if !utils.IsNormalEmail(req.Email) || len(req.Password) < 1 || len(req.Name) < 1 || !utils.IsValidCaptcha(l.svcCtx.RedisClient, l.svcCtx.CAPTCHA, req.Captcha) {
		return &types.CommonResponse{Code: 400, Msg: "参数有误"}, nil
	}

	// 判断邮箱是否重复注册(因为email是唯一索引, 所以无需自己校验)

	// 以游客身份注册
	visitor := &models.User{
		Id:       utils.GetSnowFlakeIdAndBase64(),
		Email:    req.Email,
		Password: utils.DoHashAndBase64(l.svcCtx.Config.Salt, req.Password),
		Name:     req.Name,
		Auth:     models.VisitorAuth,
	}
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserTable).InsertOne(l.ctx, visitor)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "insert user error, email is "+req.Email))
		return &types.CommonResponse{Code: 500, Msg: err.Error()}, nil
	}

	// 生成token, 返回给用户
	return common.GetToken(visitor, l.Logger, l.svcCtx)
}
