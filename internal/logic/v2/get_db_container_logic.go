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

type GetDbContainerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDbContainerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDbContainerLogic {
	return &GetDbContainerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDbContainerLogic) GetDbContainer() (resp *types.CommonResponse, err error) {
	// search data
	doc := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.DbContainerDocument)
	cursor, err := doc.Find(l.ctx, bson.D{{"user_id", l.ctx.Value("user").(*models.User).Id}})
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get DbContainerDocument data err"))
		return &types.CommonResponse{Code: 500, Msg: "get DbContainerDocument data err"}, nil
	}

	// decode data
	var containers []*models.DbContainer
	for cursor.Next(l.ctx) {
		container := new(models.DbContainer)
		if err = cursor.Decode(container); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode DbContainerDocument data err"))
			return &types.CommonResponse{Code: 500, Msg: "decode DbContainerDocument data err"}, nil
		}
		containers = append(containers, container)
	}

	// return data
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"db_containers": containers}}, nil
}
