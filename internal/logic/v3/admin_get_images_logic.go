package v3

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
