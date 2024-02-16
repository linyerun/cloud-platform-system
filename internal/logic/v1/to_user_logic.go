package v1

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToUserLogic {
	return &ToUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToUserLogic) ToUser(req *types.ApplicationFormPostRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
