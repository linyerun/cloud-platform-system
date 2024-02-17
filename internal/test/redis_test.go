package test

import (
	"cloud-platform-system/internal/common"
	"context"
	"testing"
)

func TestList(t *testing.T) {
	svcCtx.RedisClient.RPush(context.Background(), svcCtx.ExceptionList, common.NewJsonMsgString(map[string]any{"admin_id": "111", "user_id": "222", "document": "application_forms"}, "status需要修改为0"))
}
