package v3

import (
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFormByStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFormByStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFormByStatusLogic {
	return &GetFormByStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFormByStatusLogic) GetFormByStatus(req *types.ApplicationFormPostRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
