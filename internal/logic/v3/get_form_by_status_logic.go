package v3

import (
	"cloud-platform-system/internal/common/errorx"
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
	if req.Status != models.UserApplicationFormStatusIng && req.Status != models.UserApplicationFormStatusOk && req.Status != models.UserApplicationFormStatusReject {
		return &types.CommonResponse{Code: 400, Msg: "status值错误"}, nil
	}

	// 查询管理员旗下所有订单
	filter := bson.D{{"admin_id", admin.Id}}
	cur, err := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserApplicationFormDocument).Find(l.ctx, filter, options.Find().SetProjection(bson.D{{"admin_id", 0}}))
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "can not get user_application_form data err"))
		return nil, errorx.NewCodeError(500, "can not get user_application_form data err")
	}
	defer cur.Close(l.ctx)

	// 单个字段也不可以直接使用基础类型来进行解析
	type Tmp struct {
		Id          string `bson:"_id"`
		UserId      string `bson:"user_id"`
		Explanation string `bson:"explanation"`
		Status      uint   `bson:"status"`
	}

	var idToTmp = make(map[string]*Tmp)
	var ids []string

	for cur.Next(l.ctx) {
		tmp := new(Tmp)
		if err = cur.Decode(tmp); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode mongo data error"))
			return &types.CommonResponse{Code: 500, Msg: "decode mongo data error"}, nil
		}
		ids = append(ids, tmp.UserId)
		idToTmp[tmp.UserId] = tmp
	}

	if len(ids) == 0 {
		return &types.CommonResponse{Code: 200, Msg: "获取成功"}, nil
	}

	// 获取user数据
	filter = bson.D{{"_id", bson.M{"$in": ids}}}
	cur, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).Find(l.ctx, filter, options.Find().SetProjection(bson.D{{"password", 0}, {"auth", 0}}))
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "can not get userApplications data err"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	defer cur.Close(l.ctx)

	type UserApplicationData struct {
		Id          string `bson:"_id" json:"id"`
		Email       string `bson:"email" json:"email"`
		Name        string `bson:"name" json:"name"`
		Explanation string `json:"explanation"`
		Status      uint   `json:"status"`
	}

	var userApplications []*UserApplicationData
	for cur.Next(l.ctx) {
		userApplication := new(UserApplicationData)
		if err = cur.Decode(userApplication); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode mongo data error"))
			return &types.CommonResponse{Code: 500, Msg: "decode mongo data error"}, nil
		}
		tmp := idToTmp[userApplication.Id]
		userApplication.Id = tmp.Id
		userApplication.Explanation = tmp.Explanation
		userApplication.Status = tmp.Status
		userApplications = append(userApplications, userApplication)
	}
	return &types.CommonResponse{Code: 200, Msg: "获取成功", Data: map[string]any{"visitor_applications": userApplications}}, nil
}
