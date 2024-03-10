package v2

import (
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDbApplicationByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDbApplicationByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDbApplicationByIdLogic {
	return &GetDbApplicationByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDbApplicationByIdLogic) GetDbApplicationById() (resp *types.CommonResponse, err error) {
	// 查询user名下Db申请
	filter := bson.D{{"user_id", l.ctx.Value("user").(*models.User).Id}}
	cursor, err := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbApplicationFormDocument).Find(l.ctx, filter)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get models.DbApplicationFormDocument error"))
		return &types.CommonResponse{Code: 500, Msg: "get models.DbApplicationFormDocument error"}, nil
	}

	// decode msg
	var forms []*models.DbApplicationForm
	for cursor.Next(l.ctx) {
		form := new(models.DbApplicationForm)
		if err = cursor.Decode(form); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode models.DbApplicationFormDocument error"))
			return nil, errorx.NewCodeError(500, "decode models.DbApplicationFormDocument error")
		}
		forms = append(forms, form)
	}

	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"db_application_forms": forms}}, nil
}
