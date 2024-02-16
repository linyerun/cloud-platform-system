package v0

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaEmailLogic {
	return &CaptchaEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaEmailLogic) CaptchaEmail(req *types.CaptchaEmailRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
