package v3

import (
	"cloud-platform-system/internal/models"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFormByStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFormByStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFormByStatusLogic {
	return &GetFormByStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFormByStatusLogic) GetFormByStatus(req *types.GetFormByStatusRequest) (resp *types.CommonResponse, err error) {
	admin := l.ctx.Value("user").(*models.User)
	// 校验status的范围
	if req.Status != models.ApplicationFormStatusIng && req.Status != models.ApplicationFormStatusOk && req.Status != models.ApplicationFormStatusReject {
		return &types.CommonResponse{Code: 400, Msg: "status值错误"}, nil
	}

	// 查询
	filter := bson.D{{"admin_id", admin.Id}, {"status", req.Status}}
	cur, err := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.ApplicationFormTable).Find(l.ctx, filter, options.Find().SetProjection(bson.D{{"user_id", 1}}))
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "can not get application_form data err"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	defer cur.Close(l.ctx)
	var ids []string
	for cur.Next(l.ctx) {
		// 单个字段也不可以直接使用基础类型来进行解析
		type Tmp struct {
			Id string `bson:"user_id"`
		}
		tmp := new(Tmp)
		if err = cur.Decode(tmp); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode mongo data error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
		ids = append(ids, tmp.Id)
	}
	if len(ids) == 0 {
		return &types.CommonResponse{Code: 200, Msg: "获取成功"}, nil
	}

	// 获取user数据
	filter = bson.D{{"_id", bson.M{"$in": ids}}}
	cur, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserTable).Find(l.ctx, filter, options.Find().SetProjection(bson.D{{"password", 0}, {"auth", 0}}))
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "can not get users data err"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	defer cur.Close(l.ctx)
	type UserData struct {
		Id    string `bson:"_id" json:"id,omitempty"`
		Email string `bson:"email" json:"email,omitempty"`
		Name  string `bson:"name" json:"name,omitempty"`
	}
	var users []*UserData
	for cur.Next(l.ctx) {
		user := new(UserData)
		if err = cur.Decode(user); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode mongo data error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
		users = append(users, user)
	}
	return &types.CommonResponse{Code: 200, Msg: "获取成功", Data: map[string]any{"users": users}}, nil
}
