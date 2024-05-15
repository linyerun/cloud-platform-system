package v5

import (
	"context"

	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
)

type ChangeUserMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeUserMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUserMsgLogic {
	return &ChangeUserMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeUserMsgLogic) ChangeUserMsg(req *types.ChangeUserMsgReq) (resp *types.CommonResponse, err error) {
	id := l.ctx.Value("user").(*models.User).Id

	// 校验验证码
	//err = utils.IsValidEmailCaptcha(l.svcCtx.RedisClient, l.svcCtx.CAPTCHA, req.Captcha, req.Email)
	//if err != nil {
	//	return nil, err
	//}

	update := bson.M{}

	// 名称
	if len(req.Name) != 0 {
		update["name"] = req.Name
	}

	// 密码
	if len(req.Password) >= 6 {
		update["password"] = utils.DoHashAndBase64(l.svcCtx.Config.Salt, req.Password)
	} else if len(req.Password) != 0 {
		return nil, errorx.NewBaseError(400, "password长度不可小于6")
	}

	// 没有可修改的东西就直接返回就行了
	if len(update) == 0 {
		return &types.CommonResponse{Code: 200, Msg: "无可修改数据"}, nil
	}

	// 执行修改操作
	_, err = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).UpdateByID(l.ctx, id, bson.D{{"$set", update}})
	if err != nil {
		return nil, errors.Wrapf(errorx.NewBaseError(500, "update user error"), "err: %v, req: %+v", err, req)
	}

	return &types.CommonResponse{Code: 200, Msg: "成功"}, nil
}
