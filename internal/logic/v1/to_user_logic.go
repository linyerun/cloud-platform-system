package v1

import (
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToUserLogic {
	return &ToUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToUserLogic) ToUser(req *types.ApplicationFormPostRequest) (resp *types.CommonResponse, err error) {
	// 校验是否存在
	filter := bson.D{{"_id", req.AdminId}, {"email", req.AdminEmail}}
	err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).FindOne(l.ctx, filter).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		l.Logger.Error(errors.Wrap(err, "search admin in mongo error"))
		return &types.CommonResponse{Code: 500, Msg: "系统错误"}, nil
	} else if err == mongo.ErrNoDocuments {
		return &types.CommonResponse{Code: 400, Msg: "该管理员不存在"}, nil
	}

	// 校验是否重复申请
	user := l.ctx.Value("user").(*models.User)
	filter = bson.D{{"user_id", user.Id}, {"admin_id", req.AdminId}}
	err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.ApplicationFormDocument).FindOne(l.ctx, filter).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		l.Logger.Error(errors.Wrap(err, "search admin in mongo error"))
		return &types.CommonResponse{Code: 500, Msg: "系统错误"}, nil
	} else if err == nil {
		return &types.CommonResponse{Code: 400, Msg: "不可重复发出申请"}, nil
	}

	// 新增到文档当中
	af := &models.ApplicationForm{Id: utils.GetSnowFlakeIdAndBase64(), UserId: user.Id, AdminId: req.AdminId, Status: models.ApplicationFormStatusIng}
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.ApplicationFormDocument).InsertOne(l.ctx, af)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "insert application_form in mongo error"))
		return &types.CommonResponse{Code: 500, Msg: "系统错误"}, nil
	}

	// 发送申请信息到管理员邮箱(失败也无所谓)
	err = utils.SendTextByEmail(req.AdminEmail, "游客转用户申请通知", fmt.Sprintf("请即使进入管理员系统处理%s[%s]的申请请求", user.Name, user.Email))
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "send email error"))
		return &types.CommonResponse{Code: 500, Msg: "系统错误"}, nil
	}
	return &types.CommonResponse{Code: 200, Msg: "申请成功"}, nil
}
