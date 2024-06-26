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

func DelData(log logx.Logger, scvCtx *svc.ServiceContext, delFunc func() (any, error), delDocument string) error {
	data, err := delFunc()
	if err != nil {
		return err
	}
	_, err = scvCtx.MongoClient.Database(scvCtx.Config.Mongo.DbName).Collection(models.DelMessageDocument).InsertOne(context.Background(), models.NewDelMessage(utils.GetSnowFlakeIdAndBase64(), data, delDocument))
	if err != nil {
		log.Error(errors.Wrap(err, "save data to del_message error"))
		scvCtx.RedisClient.RPush(context.Background(), scvCtx.ExceptionList, NewJsonMsgString(data, "保存该数据到Mongo中的"+delDocument))
		return SaveMongoDelDataError
	}
	return nil
}
