package v2

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetDbImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserGetDbImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetDbImagesLogic {
	return &UserGetDbImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetDbImagesLogic) UserGetDbImages() (resp *types.CommonResponse, err error) {
	return common.GetDbImageList(l.svcCtx, l.Logger, l.ctx)
}
