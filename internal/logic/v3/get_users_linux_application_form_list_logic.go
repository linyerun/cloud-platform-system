package v3

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

type GetUsersLinuxApplicationFormListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersLinuxApplicationFormListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersLinuxApplicationFormListLogic {
	return &GetUsersLinuxApplicationFormListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersLinuxApplicationFormListLogic) GetUsersLinuxApplicationFormList() (resp *types.CommonResponse, err error) {
	adminId := l.ctx.Value("user").(*models.User).Id

	// 获取名下用户
	cursor, err := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserApplicationFormDocument).Find(l.ctx, bson.D{{"admin_id", adminId}, {"status", models.UserApplicationFormStatusOk}}, options.Find().SetProjection(bson.D{{"user_id", 1}}))
	if err != nil && err != mongo.ErrNoDocuments {
		l.Logger.Error(errors.Wrap(err, "get models.UserApplicationFormDocument error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	} else if err == mongo.ErrNoDocuments {
		return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"forms": nil}}, nil
	}

	// decode data
	type UserIdTmp struct {
		UserId string `bson:"user_id"`
	}
	ut := new(UserIdTmp)
	var ids []string
	for cursor.Next(l.ctx) {
		// decode
		if err = cursor.Decode(ut); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode models.UserApplicationFormDocument error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
		ids = append(ids, ut.UserId)
	}

	// 名下没有用户, 那就没数据了
	if len(ids) == 0 {
		return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
	}

	// get form msg
	filter := bson.D{{"user_id", bson.M{"$in": ids}}}
	cursor, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).Find(l.ctx, filter)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get models.LinuxApplicationFormDocument error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	var forms []*models.LinuxApplicationForm
	for cursor.Next(l.ctx) {
		form := new(models.LinuxApplicationForm)
		if err = cursor.Decode(form); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode models.LinuxApplicationFormDocument error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
		forms = append(forms, form)
	}

	// 返回结果
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"linux_application_forms": forms}}, nil
}
