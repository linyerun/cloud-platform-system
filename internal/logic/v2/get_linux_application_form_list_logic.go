package v2

import (
	"cloud-platform-system/internal/models"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLinuxApplicationFormListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLinuxApplicationFormListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLinuxApplicationFormListLogic {
	return &GetLinuxApplicationFormListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLinuxApplicationFormListLogic) GetLinuxApplicationFormList() (resp *types.CommonResponse, err error) {
	doc := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument)
	cursor, err := doc.Find(l.ctx, bson.D{{"user_id", l.ctx.Value("user").(*models.User).Id}})
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get LinuxApplicationFormDocument msg err"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	var forms []*models.LinuxApplicationForm
	for cursor.Next(l.ctx) {
		form := new(models.LinuxApplicationForm)
		if err = cursor.Decode(form); err != nil {
			l.Logger.Error(errors.Wrap(err, "code LinuxApplicationFormDocument data err"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
		forms = append(forms, form)
	}
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"linux_application_forms": forms}}, nil
}
