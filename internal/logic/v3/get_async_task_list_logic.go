package v3

import (
	"cloud-platform-system/internal/models"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAsyncTaskListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAsyncTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAsyncTaskListLogic {
	return &GetAsyncTaskListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAsyncTaskListLogic) GetAsyncTaskList() (resp *types.CommonResponse, err error) {
	userId := l.ctx.Value("user").(*models.User).Id
	// 查询数据
	cursor, err := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.AsyncTaskDocument).Find(l.ctx, bson.D{{"user_id", userId}})
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get async task list error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// decode data
	var tasks []*models.AsyncTask
	for cursor.Next(l.ctx) {
		task := new(models.AsyncTask)
		if err = cursor.Decode(task); err != nil {
			l.Logger.Error(errors.Wrap(err, "decode async task list error"))
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
		tasks = append(tasks, task)
	}
	return &types.CommonResponse{Code: 200, Msg: "成功", Data: map[string]any{"async_tasks": tasks}}, nil
}
