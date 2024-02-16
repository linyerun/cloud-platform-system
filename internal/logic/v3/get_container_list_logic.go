package v3

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetContainerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContainerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContainerListLogic {
	return &GetContainerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetContainerListLogic) GetContainerList() (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
