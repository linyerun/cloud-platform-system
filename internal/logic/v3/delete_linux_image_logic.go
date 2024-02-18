package v3

import (
	"cloud-platform-system/internal/common"
	"cloud-platform-system/internal/models"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLinuxImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLinuxImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLinuxImageLogic {
	return &DeleteLinuxImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLinuxImageLogic) DeleteLinuxImage(req *types.ImageDelRequest) (resp *types.CommonResponse, err error) {
	// 根据Id获取镜像信息
	res := l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).FindOne(l.ctx, bson.D{{"_id", req.Id}}, options.FindOne().SetProjection(bson.D{{"image_id", 1}}))
	if err = res.Err(); err == mongo.ErrNoDocuments {
		return &types.CommonResponse{Code: 400, Msg: "不存在此镜像"}, nil
	} else if err != nil {
		l.Logger.Error(errors.Wrap(err, "get image msg in mongo error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	var dockerImage = new(models.LinuxImage) // 只有ImageId有数据
	err = res.Decode(dockerImage)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "decode mongo data error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 获取分布式锁，然后执行删除操作
	// 使用锁来保证同一个镜像只能有一个管理员删除
	ret := l.svcCtx.RedisClient.SetNX(l.ctx, fmt.Sprintf(l.svcCtx.ImagePrefix, dockerImage.ImageId), "1", time.Second*30)
	if err = res.Err(); err != nil {
		l.Logger.Error(errors.Wrap(err, "error to apply for a distributed lock"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	ok, err := ret.Result()
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "get result from redis error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}
	if !ok {
		return &types.CommonResponse{Code: 401, Msg: "已有管理员删除此镜像，无所重复删除！"}, nil
	}
	// 归还分布式锁
	defer l.svcCtx.RedisClient.Del(l.ctx, fmt.Sprintf(l.svcCtx.ImagePrefix, dockerImage.ImageId))

	// 执行删除操作
	err = l.svcCtx.DockerClient.RemoveImage(dockerImage.ImageId)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "remove image error"))
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 删除镜像在数据库数据，然后保存于别的文档
	err = common.DelData(l.Logger, l.svcCtx, func() (any, error) {
		// 真出现异常了，我们可以根据删除日志恢复镜像
		filter := bson.D{{"_id", req.Id}}
		res = l.svcCtx.MongoClient.Database(l.svcCtx.Config.Mongo.DbName).Collection(models.LinuxImageDocument).FindOneAndDelete(l.ctx, filter)
		if err = res.Err(); err != nil {
			l.Logger.Error(errors.Wrap(err, "get and del mongo data error"))
			return nil, err
		}
		dockerImage = new(models.LinuxImage)
		if err = res.Decode(dockerImage); err != nil {
			l.Logger.Error(errors.Wrap(err, "get and del mongo data error"))
			return nil, err
		}
		return dockerImage, nil
	})
	if err != nil && err != common.SaveMongoDelDataError {
		return &types.CommonResponse{Code: 500, Msg: "系统异常"}, nil
	}

	// 返回成功信息
	return &types.CommonResponse{Code: 200, Msg: "删除成功"}, nil
}
