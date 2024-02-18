package v3

import (
	"cloud-platform-system/internal/common"
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetLinuxImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetLinuxImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetLinuxImagesLogic {
	return &AdminGetLinuxImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetLinuxImagesLogic) AdminGetLinuxImages() (resp *types.CommonResponse, err error) {
	return common.GetAllImage(l.svcCtx, l.Logger, l.ctx)
}
