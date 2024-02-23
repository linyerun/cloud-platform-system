package asynctask

import (
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type RespData struct {
	Code int    `json:"code" bson:"code"`
	Msg  string `json:"msg" bson:"msg"`
	Data any    `json:"data" bson:"data"`
}

func (r *RespData) JsonMarshal() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}

type MyAsyncTaskPoolServer struct {
	cancelFunc context.CancelFunc
	ctx        context.Context
	svcCtx     *svc.ServiceContext
}

func NewMyAsyncTaskPoolServer(ctx context.Context, svcCtx *svc.ServiceContext, cancelFunc context.CancelFunc) *MyAsyncTaskPoolServer {
	return &MyAsyncTaskPoolServer{
		cancelFunc: cancelFunc,
		ctx:        ctx,
		svcCtx:     svcCtx,
	}
}

func (i *MyAsyncTaskPoolServer) Start() {
	ctx, svcCtx := i.ctx, i.svcCtx
	defer logx.Info("AsyncTaskHandler Stop")
	taskHandlerFactory := NewTaskHandlerFactory(ctx, svcCtx)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 从MongoDB中拉取任务来进行处理
			result := svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.AsyncTaskDocument).FindOne(ctx, bson.D{{"status", models.AsyncTaskIng}}, options.FindOne().SetSort(bson.D{{"priority", -1}}))
			if err := result.Err(); err != nil && err != mongo.ErrNoDocuments { // 拉取异步任务失败
				logx.Error(errors.Wrap(err, "get AsyncTask data from mongo error"))
			} else if err == nil { // 拉取到一个异步任务
				asyncTask := new(models.AsyncTask)
				if err = result.Decode(asyncTask); err != nil { // 解码失败
					logx.Error(errors.Wrap(err, "decodeAsyncTask data error"))
				} else {
					// 处理任务(后面这个可以采用协程池来处理)
					go func() {
						// 异常捕获器，避免因异常而退出
						defer func() {
							if errMsg := recover(); err != nil {
								status := models.AsyncTaskFail
								respData := &RespData{Code: 500, Msg: fmt.Sprintf("%v", errMsg)}
								if _, err = svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.AsyncTaskDocument).UpdateOne(ctx, bson.D{{"_id", asyncTask.Id}}, bson.D{{"$set", bson.M{"status": status, "resp_data": respData, "finish_at": time.Now().UnixMilli()}}}); err != nil {
									logx.Error(errors.Wrap(err, "update AsyncTask error"))
								}
							}
							logx.Infof("[id=%s]的异步任务处理结束", asyncTask.Id)
						}()
						// 获取处理器
						handler, ok := taskHandlerFactory.NewTaskHandler(asyncTask.Type)
						var respData *RespData
						var status uint
						if !ok {
							respData = &RespData{Code: 500, Msg: "系统异常，无法找到处理该类型异步任务的处理器"}
							status = models.AsyncTaskFail
						} else {
							// 处理器获取成功，执行处理逻辑
							respData, status = handler.Execute(asyncTask.Args)
						}
						// 处理完成后把结果保存回mongo中
						_, err = svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.AsyncTaskDocument).UpdateOne(ctx, bson.D{{"_id", asyncTask.Id}}, bson.D{{"$set", bson.M{"status": status, "resp_data": respData.JsonMarshal(), "finish_at": time.Now().UnixMilli()}}})
						if err != nil {
							logx.Error(errors.Wrap(err, "update AsyncTask error"))
						}
					}()
				}
			}
			// 休眠一段时间后再拉取
			time.Sleep(time.Millisecond * time.Duration(svcCtx.Config.AsyncTask.PullTaskWaitMillSec))
		}
	}
}

func (i *MyAsyncTaskPoolServer) Stop() {
	i.cancelFunc()
}
