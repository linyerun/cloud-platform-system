package asynctask

import (
	"cloud-platform-system/internal/svc"
	"context"
)

type IAsyncTaskHandler interface {
	Execute(args string) (respData *RespData, status uint)
}

type TaskHandlerFactory struct {
	ctx    context.Context
	srvCtx *svc.ServiceContext
}

func NewTaskHandlerFactory(ctx context.Context, srvCtx *svc.ServiceContext) *TaskHandlerFactory {
	return &TaskHandlerFactory{ctx: ctx, srvCtx: srvCtx}
}

func (f *TaskHandlerFactory) NewTaskHandler(taskType string) (handler IAsyncTaskHandler, ok bool) {
	ok = true
	switch taskType {
	case ImagePullType:
		handler = NewImagePullHandler(f.ctx, f.srvCtx)
	default:
		return nil, false
	}
	return
}
