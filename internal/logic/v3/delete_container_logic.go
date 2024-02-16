package v3

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteContainerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteContainerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteContainerLogic {
	return &DeleteContainerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteContainerLogic) DeleteContainer(req *types.DeleteContainerRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
