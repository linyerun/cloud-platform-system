package v2

import (
	"cloud-platform-system/internal/common"
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetLinuxImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserGetLinuxImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetLinuxImagesLogic {
	return &UserGetLinuxImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetLinuxImagesLogic) UserGetLinuxImages() (resp *types.CommonResponse, err error) {
	return common.GetAllImage(l.svcCtx, l.Logger, l.ctx)
}
