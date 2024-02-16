package v2

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserGetImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetImagesLogic {
	return &UserGetImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetImagesLogic) UserGetImages() (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
