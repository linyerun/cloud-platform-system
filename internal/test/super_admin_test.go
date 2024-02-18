package test

import (
	"cloud-platform-system/internal/config"
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/utils"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"testing"
)

var svcCtx *svc.ServiceContext

func init() {
	var c config.Config
	conf.MustLoad("../../etc/cloud_platform.yaml", &c)
	svcCtx = svc.NewServiceContext(c)
}

func TestInsertSuperAdmin(t *testing.T) {
	admin := &models.User{
		// TODO: 为啥运行太快获取的是零值的Base64
		Id:       utils.GetSnowFlakeIdAndBase64(),
		Email:    "linyerun0620@qq.com",
		Password: utils.DoHashAndBase64(svcCtx.Config.Salt, "123456"),
		Name:     "超级管理员",
		Auth:     models.SuperAdminAuth,
	}
	result, err := svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.UserDocument).InsertOne(context.Background(), admin)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
