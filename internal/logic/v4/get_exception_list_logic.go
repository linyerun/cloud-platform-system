package v4

import (
	"cloud-platform-system/internal/common"
	"context"
	"github.com/pkg/errors"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExceptionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExceptionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExceptionListLogic {
	return &GetExceptionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExceptionListLogic) GetExceptionList() (resp *types.CommonResponse, err error) {
	// 从redis中获取数据
	exceptionList, err := l.svcCtx.RedisClient.LRange(l.ctx, l.svcCtx.ExceptionList, 0, -1).Result()
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get redis list ExceptionList data error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// decode json data to object
	var list []*common.JsonMsg
	for _, exception := range exceptionList {
		list = append(list, common.NewJsonMsg([]byte(exception)))
	}

	// return data to super admin
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"exception_list": list}}, nil
}
