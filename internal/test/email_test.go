package test

import (
	"cloud-platform-system/internal/utils"
	"testing"
)

func TestSendText(t *testing.T) {
	err := utils.SendTextByEmail("linyerundgut@126.com", "注册通知", "注册成功!")
	if err != nil {
		t.Fatal(err)
	}
	//utils.SendTextByEmail("2338244917@qq.com", "哈哈哈哈")
}
