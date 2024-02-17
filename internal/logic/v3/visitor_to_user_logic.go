package v3

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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
	// 校验status值
	if req.Status != models.ApplicationFormStatusOk && req.Status != models.ApplicationFormStatusReject {
		return &types.CommonResponse{Code: 400, Msg: "status参数有误"}, nil
	}

	// 验证是否存在这个申请, 存在则改正
	admin := l.ctx.Value("user").(*models.User)
	filterForm := bson.D{{"admin_id", admin.Id}, {"user_id", req.VisitorId}, {"status", models.ApplicationFormStatusIng}}
	err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.ApplicationFormTable).FindOneAndUpdate(l.ctx, filterForm, bson.D{{"$set", bson.M{"status": req.Status}}}).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		l.Logger.Error(errors.Wrap(err, "find and modify application_forms err"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	} else if err == mongo.ErrNoDocuments {
		return &types.CommonResponse{Code: 400, Msg: "不存在该申请"}, nil
	}

	// 如果是拒绝则不需要执行修改user的auth操作, 同意才需要执行
	if req.Status == models.ApplicationFormStatusOk {
		// 改变user的auth
		filterUser := bson.D{{"_id", req.VisitorId}, {"email", req.VisitorEmail}}
		_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserTable).UpdateOne(l.ctx, filterUser, bson.D{{"$set", bson.M{"auth": uint(models.UserAuth)}}})
		if err != nil {
			l.Logger.Error(errors.Wrap(err, "modify user err"))
			// 最大程度自救, 大概率不会用上的
			filterForm = bson.D{{"admin_id", admin.Id}, {"user_id", req.VisitorId}, {"status", req.Status}}
			for i := 0; i < 3; i++ {
				err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.ApplicationFormTable).FindOneAndUpdate(l.ctx, filterForm, bson.D{{"$set", bson.M{"status": models.ApplicationFormStatusIng}}}).Err()
				if err != nil {
					l.Logger.Error(errors.Wrap(err, "find and modify application_forms err"))
					time.Sleep(time.Millisecond * 50)
				} else {
					break
				}
			}
			if err != nil {
				// 把这个记录保存到redis中，等待后面再进行处理
				l.svcCtx.RedisClient.RPush(l.ctx, l.svcCtx.ExceptionList, common.NewJsonMsgString(map[string]any{"admin_id": admin.Id, "user_id": req.VisitorId, "document": "application_forms"}, "status需要修改为0"))
			}
			return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
		}
	}

	// 发送通知给user
	var ret string
	if req.Status == models.ApplicationFormStatusOk {
		ret = "管理员通过了您的申请，您现在的身份由游客转为团队用户了。"
	} else {
		ret = "管理员拒绝了您的申请，您现在的身份依旧是游客。"
	}
	if err = utils.SendTextByEmail(req.VisitorEmail, "管理员审核结果通知", ret); err != nil {
		l.Logger.Error(errors.Wrap(err, "send email err"))
		return &types.CommonResponse{Code: 501, Msg: "邮箱发送系统异常，但审核已经通过，请另行通知您的团队成员。"}, nil
	}

	//	返回审核结果
	return &types.CommonResponse{Code: 200, Msg: "审核请求处理完成"}, nil
}
