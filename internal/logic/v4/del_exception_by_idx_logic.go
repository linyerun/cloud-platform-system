package v4

import (
	"cloud-platform-system/internal/utils"
	"context"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelExceptionByIdxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelExceptionByIdxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelExceptionByIdxLogic {
	return &DelExceptionByIdxLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelExceptionByIdxLogic) DelExceptionByIdx(req *types.DelExceptionByIdxReq) (resp *types.CommonResponse, err error) {
	// 获取长度
	result := l.svcCtx.RedisClient.LLen(l.ctx, l.svcCtx.ExceptionList)

	// 判断idx是否在范围内
	listLen, err := result.Uint64()
	if err != nil {
		l.Logger.Error(err)
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	if req.Idx < 0 && uint64(req.Idx) >= listLen {
		return &types.CommonResponse{Code: 400, Msg: "超出范围"}, nil
	}

	// 设置特殊值
	val := utils.GetSnowFlakeIdAndBase64()
	if err = l.svcCtx.RedisClient.LSet(l.ctx, l.svcCtx.ExceptionList, int64(req.Idx), val).Err(); err != nil {
		l.Logger.Error(err)
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 删除
	err = l.svcCtx.RedisClient.LRem(l.ctx, l.svcCtx.ExceptionList, 0, val).Err()
	if err != nil {
		l.Logger.Error(err)
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 返回结果
	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}
