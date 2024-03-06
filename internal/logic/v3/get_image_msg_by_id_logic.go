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

type GetImageMsgByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetImageMsgByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImageMsgByIdLogic {
	return &GetImageMsgByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetImageMsgByIdLogic) GetImageMsgById(req *types.GetImageMsgByIdReq) (resp *types.GetImageMsgByIdResp, err error) {
	res := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).FindOne(l.ctx, bson.D{{"_id", req.Id}})
	if err = res.Err(); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "获取image_msg数据失败")
	}

	// 解析image
	image := new(models.LinuxImage)
	if err = res.Decode(image); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(400, "解析Image数据失败")
	}

	res = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).FindOne(l.ctx, bson.D{{"_id", image.CreatorId}})
	if err = res.Err(); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(500, "获取user_msg数据失败")
	}

	// 解析user
	user := new(models.User)
	if err = res.Decode(user); err != nil {
		l.Logger.Error(err)
		return nil, errorx.NewCodeError(400, "解析User数据失败")
	}

	// 拼接resp
	resp = &types.GetImageMsgByIdResp{
		CreateorName:    user.Name,
		CreatorEmail:    user.Email,
		ImageName:       image.Name,
		ImageTag:        image.Tag,
		ImageSize:       image.Size,
		EnableCommands:  image.EnableCommands,
		MustExportPorts: image.MustExportPorts,
	}

	return
}
