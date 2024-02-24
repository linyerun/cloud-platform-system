package v2

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
