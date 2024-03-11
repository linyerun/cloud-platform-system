package v3

import (
	"cloud-platform-system/internal/asynctask"
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type HandleUserLinuxApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleUserLinuxApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleUserLinuxApplicationLogic {
	return &HandleUserLinuxApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandleUserLinuxApplicationLogic) HandleUserLinuxApplication(req *types.HandleUserLinuxApplicationReq) (resp *types.CommonResponse, err error) {
	// 判断status是否合法
	if req.Status != models.LinuxApplicationFormStatusReject && req.Status != models.LinuxApplicationFormStatusOk {
		return nil, errorx.NewBaseError(400, "status参数有误, 不能设置为除成功/失败外的其他状态")
	}

	// 判断这个form是否存在
	filter := bson.D{{"_id", req.FormId}}
	result := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).FindOne(l.ctx, filter)
	if err = result.Err(); err != nil {
		logx.Error(errors.Wrap(err, "get LinuxApplicationFormDocument error"))
		return nil, errorx.NewBaseError(500, "查找申请单失败")
	}
	// decode
	form := new(models.LinuxApplicationForm)
	if err = result.Decode(form); err != nil {
		logx.Error(errors.Wrap(err, "decode LinuxApplicationFrom error"))
		return nil, errorx.NewBaseError(500, "解码申请单数据失败")
	}
	// 判断这个申请单是否被处理过了
	if form.Status != models.LinuxApplicationFormStatusIng {
		return nil, errorx.NewBaseError(400, "该Linux服务器开启申请单已被处理")
	}

	// 修改订单状态
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxApplicationFormDocument).UpdateOne(l.ctx, filter, bson.D{{"$set", bson.M{"status": models.LinuxApplicationFormStatusAsync}}})
	if err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewBaseError(500, "更新叮当状态失败")
	}

	// 把异步任务保存到mongo中等待被执行
	asyncTask := &models.AsyncTask{
		Id:       utils.GetSnowFlakeIdAndBase64(),
		UserId:   l.ctx.Value("user").(*models.User).Id,
		Type:     asynctask.ContainerRunType,
		Args:     (&asynctask.ContainerRunArgs{FormId: req.FormId, Status: req.Status}).JsonMarshal(),
		Status:   models.AsyncTaskIng,
		CreateAt: time.Now().UnixMilli(),
	}
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.AsyncTaskDocument).InsertOne(l.ctx, asyncTask)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "insert async task error"))
		err = l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]any{"table_name": models.LinuxApplicationFormDocument, "_id": req.FormId, "status": form.Status}, "赶紧改成这个status")).Err()
		if err != nil {
			l.Logger.Error(common.NewJsonMsgString(map[string]any{"table_name": models.LinuxApplicationFormDocument, "_id": req.FormId, "status": form.Status}, "赶紧改成这个status"))
		}
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	return &types.CommonResponse{Code: 200, Msg: "已提交审核，等待异步处理"}, nil
}
