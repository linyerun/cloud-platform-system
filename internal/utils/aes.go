package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

const key = "2338244917lyr666"

var myCipher cipher.Block

func init() {
	// 创建对称加密
	var err error
	myCipher, err = aes.NewCipher([]byte(key)) // key的字节数为: 16B
	if err != nil {                            // 报错就不加密了
		panic(err)
	}
}

func AESEncrypt(data []byte) (encData []byte) {
	// 填充够myCipher.BlockSize()字节的整数倍
	var tmp []byte
	for i, n := 0, myCipher.BlockSize()-len(data)%myCipher.BlockSize(); i < n; i++ {
		tmp = append(tmp, byte(n))
	}
	data = append(data, tmp...)
	// 使用AES加密
	encData = make([]byte, len(data))
	for i, n := 0, len(data)/myCipher.BlockSize(); i < n; i++ {
		myCipher.Encrypt(encData[i*myCipher.BlockSize():], data[i*myCipher.BlockSize():])
	}
	return
}

func AESDecrypt(encData []byte) []byte {
	// 使用AES解密
	data := make([]byte, len(encData))
	for i, n := 0, len(data)/myCipher.BlockSize(); i < n; i++ {
		myCipher.Decrypt(data[i*myCipher.BlockSize():], encData[i*myCipher.BlockSize():])
	}
	// 去除填充的数据
	n := len(data) - int(data[len(data)-1])
	return data[:n]
}
