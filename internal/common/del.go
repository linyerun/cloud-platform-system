package common

import (
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

var SaveMongoDelDataError = errors.New("save mongo del data error")

func DelData(log logx.Logger, scvCtx *svc.ServiceContext, delFunc func() (any, error)) error {
	data, err := delFunc()
	if err != nil {
		return err
	}
	_, err = scvCtx.MongoClient.Database(scvCtx.Config.Mongo.DbName).Collection(models.DelMessageDocument).InsertOne(context.Background(), models.NewDelMessage(utils.GetSnowFlakeIdAndBase64(), data))
	log.Error(errors.Wrap(err, "save data to del_message error"))
	scvCtx.RedisClient.RPush(context.Background(), scvCtx.ExceptionList, NewJsonMsgString(data, "保存该models.Image数据到Mongo中"))
	return SaveMongoDelDataError
}
