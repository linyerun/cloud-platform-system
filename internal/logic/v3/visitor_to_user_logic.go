package v3

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VisitorToUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVisitorToUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VisitorToUserLogic {
	return &VisitorToUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VisitorToUserLogic) VisitorToUser(req *types.PutVisitorToUserRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
