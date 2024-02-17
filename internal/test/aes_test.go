package test

import (
	"cloud-platform-system/internal/utils"
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	data := []byte{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}
	encData := utils.AESEncrypt(data)

	decData := utils.AESDecrypt(encData)

	fmt.Println(data)
	fmt.Println(decData)
}
