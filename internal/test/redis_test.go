package test

import (
	"cloud-platform-system/internal/common"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	svcCtx.RedisClient.RPush(context.Background(), svcCtx.ExceptionList, common.NewJsonMsgString(map[string]any{"admin_id": "111", "user_id": "222", "document": "application_forms"}, "status需要修改为0"))
}

func TestSetNx(t *testing.T) {
	ret := svcCtx.RedisClient.SetNX(context.Background(), fmt.Sprintf(svcCtx.ImagePrefix, "image:latest"), "1", time.Second*100)
	if ret.Err() != nil {
		fmt.Println(ret.Err())
		return
	}
	result, err := ret.Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("res:", result)
}
