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

type GetLinuxContainerByUserIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLinuxContainerByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLinuxContainerByUserIdLogic {
	return &GetLinuxContainerByUserIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLinuxContainerByUserIdLogic) GetLinuxContainerByUserId() (resp *types.CommonResponse, err error) {
	doc := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxContainerDocument)
	cursor, err := doc.Find(l.ctx, bson.D{{"user_id", l.ctx.Value("user").(*models.User).Id}})
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get LinuxContainerDocument data err"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	var containers []*models.LinuxContainer
	for cursor.Next(l.ctx) {
		container := new(models.LinuxContainer)
		if err = cursor.Decode(container); err != nil {
			l.Logger.Error(errors.Wrap(err, "code LinuxContainerDocument data err"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
		containers = append(containers, container)
	}
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"linux_containers": containers}}, nil
}
