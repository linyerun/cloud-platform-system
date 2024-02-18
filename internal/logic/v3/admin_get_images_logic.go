package v3

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetImagesLogic {
	return &AdminGetImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetImagesLogic) AdminGetImages() (resp *types.CommonResponse, err error) {
	return common.GetAllImage(l.svcCtx, l.Logger, l.ctx)
}
