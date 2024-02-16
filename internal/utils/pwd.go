package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func DoHashAndBase64(salt, pwd string) string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s-%s", pwd, salt)))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func VerifyHash(salt, pwd, hash string) (valid bool) {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s-%s", pwd, salt)))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)) == hash
}
