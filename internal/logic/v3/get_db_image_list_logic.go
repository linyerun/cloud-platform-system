package v3

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDbImageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDbImageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDbImageListLogic {
	return &GetDbImageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDbImageListLogic) GetDbImageList() (resp *types.CommonResponse, err error) {
	return common.GetDbImageList(l.svcCtx, l.Logger, l.ctx)
}
