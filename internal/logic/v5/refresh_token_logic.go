package v5

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"context"
	"github.com/pkg/errors"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken() (resp *types.CommonResponse, err error) {
	user, ok := l.ctx.Value("user").(*models.User)
	if !ok {
		panic(errors.New("can not get user in context value"))
	}

	// 把信息封装到token中并返回给用户
	return common.GetToken(user, l.Logger, l.svcCtx)
}
