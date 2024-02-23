package v2

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LinuxStartApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLinuxStartApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LinuxStartApplyLogic {
	return &LinuxStartApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LinuxStartApplyLogic) LinuxStartApply(req *types.LinuxStartApplyRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
