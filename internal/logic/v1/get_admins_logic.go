package v1

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

type GetAdminsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminsLogic {
	return &GetAdminsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminsLogic) GetAdmins() (resp *types.CommonResponse, err error) {
	// 查询数据
	filter := bson.D{{"auth", models.AdminAuth}}
	cur, err := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserTable).Find(l.ctx, filter, options.Find().SetProjection(bson.M{"_id": 1, "email": 1, "name": 1}))
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get admins from mongo error"))
		return &types.CommonResponse{Code: 500, Msg: "系统错误，获取数据失败"}, nil
	}
	defer cur.Close(l.ctx)

	// 解码数据
	type Admin struct {
		Id    string `bson:"_id" json:"id"`
		Email string `bson:"email" json:"email"`
		Name  string `bson:"name" json:"name"`
	}
	var admins []*Admin
	for cur.Next(l.ctx) {
		admin := new(Admin)
		if err = cur.Decode(admin); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode admin from mongo data error"))
			return &types.CommonResponse{Code: 500, Msg: "系统错误，获取数据失败"}, nil
		}
		admins = append(admins, admin)
	}

	// 返回数据
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"admins": admins}}, nil
}
