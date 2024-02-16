package v3

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullImageLogic {
	return &PullImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullImageLogic) PullImage(req *types.ImagePullRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
