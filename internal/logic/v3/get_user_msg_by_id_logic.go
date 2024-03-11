package v3

import (
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserMsgByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserMsgByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMsgByIdLogic {
	return &GetUserMsgByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserMsgByIdLogic) GetUserMsgById(req *types.GetUserMsgByIdReq) (resp *types.GetUserMsgByIdResp, err error) {
	res := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).FindOne(l.ctx, bson.D{{"_id", req.Id}})
	if err = res.Err(); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewBaseError(500, "获取user_msg数据失败")
	}

	user := new(models.User)
	if err = res.Decode(user); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewBaseError(400, "解析User数据失败")
	}

	resp = &types.GetUserMsgByIdResp{Email: user.Email, Name: user.Name}

	return
}
