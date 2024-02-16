package test

import (
	"cloud-platform-system/internal/utils"
	"testing"
)

func TestSendText(t *testing.T) {
	utils.SendTextByEmail("linyerundgut@126.com", "哈哈哈哈")
}
