package utils

import (
	"cloud-platform-system/internal/models"
	"crypto"
	"crypto/hmac"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	headerStr             string
	mySha256              = &signingMethod{crypto.SHA256}
	slat                  = "my jwt for generating token"
	ErrSignatureInvalid   = errors.New("signature is invalid")
	ErrTokenFormatInvalid = errors.New("token format is invalid")
	ErrTokenClaimInvalid  = errors.New("token claim is invalid")
)

func init() {
	header := map[string]interface{}{
		"typ": "JWT",
		"alg": "LYR",
	}
	headerBytes, _ := json.Marshal(header)
	headerStr = base64.RawURLEncoding.EncodeToString(headerBytes)
}

type signingMethod struct {
	Hash crypto.Hash
}

func (m *signingMethod) Verify(signingString, signature string, keyBytes []byte) error {
	sig, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	hasher := hmac.New(m.Hash.New, keyBytes)
	hasher.Write([]byte(signingString))
	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return ErrSignatureInvalid
	}

	return nil
}

func (m *signingMethod) Sign(signingString string, keyBytes []byte) (signature string) {
	hasher := hmac.New(m.Hash.New, keyBytes)
	hasher.Write([]byte(signingString))
	signature = base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
	return
}

type TokenBinaryFormat interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type DefaultTokenObject struct {
	Claims               *models.User `json:"claims"`
	ExpireAtMilliseconds int64        `json:"expire_at_milliseconds"`
}

func NewDefaultTokenObject(claims *models.User, expireTime time.Duration) (obj *DefaultTokenObject, err error) {
	// claim必须是指针，不然解码的时候就变成map了
	if !IsNotNilPointer(claims) {
		return nil, ErrTokenClaimInvalid
	}
	obj = new(DefaultTokenObject)
	obj.Claims = claims
	obj.ExpireAtMilliseconds = expireTime.Milliseconds() + time.Now().UnixMilli()
	return
}

func (d *DefaultTokenObject) MarshalBinary() (data []byte, err error) {
	// 转成json格式
	data, err = json.Marshal(d)
	// 对数据进行对称加密
	data = AESEncrypt(data)
	return
}

func (d *DefaultTokenObject) UnmarshalBinary(data []byte) (err error) {
	// 对数据进行解密
	data = AESDecrypt(data)
	// json转对象
	err = json.Unmarshal(data, d)
	return
}

func (d *DefaultTokenObject) IsValid() bool {
	return time.Now().UnixMilli() <= d.ExpireAtMilliseconds
}

func (d *DefaultTokenObject) GenerateToken() (token string, err error) {
	token, err = GenerateToken(d)
	return
}

func GenerateToken(val TokenBinaryFormat) (token string, err error) {
	bodyBytes, err := val.MarshalBinary()
	if err != nil {
		return "", err
	}
	bodyStr := base64.RawURLEncoding.EncodeToString(bodyBytes)
	part1 := headerStr + "." + bodyStr
	part2 := mySha256.Sign(part1, []byte(slat))
	return part1 + "." + part2, nil
}

func ParseToken(token string, val TokenBinaryFormat) (err error) {
	s := strings.Split(token, ".")
	if len(s) != 3 {
		return ErrTokenFormatInvalid
	}
	err = mySha256.Verify(s[0]+"."+s[1], s[2], []byte(slat))
	if err != nil {
		return err
	}
	bytes, err := base64.RawURLEncoding.DecodeString(s[1])
	if err != nil {
		return err
	}
	err = val.UnmarshalBinary(bytes)
	return err
}
