package v3

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserByIdLogic {
	return &DeleteUserByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserByIdLogic) DeleteUserById(req *types.DeleteUserRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
