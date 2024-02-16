package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strings"
)

const key = "2338244917lyr666"

var myCipher cipher.Block

func init() {
	// 创建对称加密
	var err error
	myCipher, err = aes.NewCipher([]byte(key)) // 只要key的字节数为: 16, 24, or 32就不会报错
	if err != nil {                            // 报错就不加密了
		panic(err)
	}
}

func AESEncrypt(data []byte) (encData []byte) {
	// 填充够16B的整数倍
	encodeToString := base64.URLEncoding.EncodeToString(data)
	for i, n := 0, len(encodeToString)%16; i < n; i++ {
		encodeToString += "!"
	}
	// 使用AES加密
	encData = make([]byte, len(encodeToString))
	myCipher.Encrypt(encData, []byte(encodeToString))
	return
}

func AESDecrypt(encData []byte) []byte {
	// 使用AES解密
	data := make([]byte, len(encData))
	myCipher.Decrypt(data, encData)
	// 恢复原始信息
	s := strings.TrimRight(string(data), "!")
	bytes, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bytes
}
