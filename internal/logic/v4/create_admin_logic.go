package v4

import (
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminLogic) CreateAdmin(req *types.CreateAdminRequest) (resp *types.CommonResponse, err error) {
	// 校验参数
	if !utils.IsNormalEmail(req.Email) || len(req.Password) < 1 || len(req.Name) < 1 {
		return &types.CommonResponse{Code: 400, Msg: "参数有误"}, nil
	}

	// 判断邮箱是否重复注册
	err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).FindOne(l.ctx, bson.D{{"email", req.Email}}).Err()
	if err != mongo.ErrNoDocuments {
		return nil, errorx.NewCodeError(400, "该邮箱已被用于注册")
	}

	// 以游客身份注册
	admin := &models.User{
		Id:       utils.GetSnowFlakeIdAndBase64(),
		Email:    req.Email,
		Password: utils.DoHashAndBase64(l.svcCtx.Config.Salt, req.Password),
		Name:     req.Name,
		Auth:     models.AdminAuth,
	}
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).InsertOne(l.ctx, admin)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "insert user error, email is "+req.Email))
		return &types.CommonResponse{Code: 500, Msg: err.Error()}, nil
	}

	utils.SendTextByEmail(req.Email, "管理员账户生成成功通知", "名称: "+req.Name+"，密码: "+req.Password)

	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}
