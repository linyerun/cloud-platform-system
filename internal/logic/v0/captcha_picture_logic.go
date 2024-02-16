package v0

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaPictureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaPictureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaPictureLogic {
	return &CaptchaPictureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaPictureLogic) CaptchaPicture() (resp *types.CaptchaPictureResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
