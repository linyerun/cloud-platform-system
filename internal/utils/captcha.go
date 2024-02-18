package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"time"
)

var nums = [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

func GetRandomDigits(arrLen uint) []byte {
	if arrLen == 0 {
		return nil
	}
	res := make([]byte, arrLen)
	for i := uint(0); i < arrLen; i++ {
		res[i] = nums[rand.Intn(len(nums))]
	}
	return res
}

func BytesNumsToStringNums(digits []byte) string {
	dst := make([]byte, len(digits))
	for i := 0; i < len(digits); i++ {
		dst[i] = digits[i] + '0'
	}
	return string(dst)
}

// WriteCaptchaImage given digits, where each digit must be in range 0-9. image suffix is png.
func WriteCaptchaImage(redisClient *redis.Client, CAPTCHA string, digits []byte, width, height int) (buff *bytes.Buffer, err error) {
	buff = bytes.NewBuffer(nil)
	_, err = captcha.NewImage("", digits, width, height).WriteTo(buff)
	if err != nil {
		return nil, err
	}
	err = redisClient.Set(context.Background(), fmt.Sprintf(CAPTCHA, BytesNumsToStringNums(digits)), 1, time.Minute*10).Err()
	return
}

func IsValidCaptcha(log logx.Logger, redisClient *redis.Client, CAPTCHA string, captchaStr string) bool {
	err := redisClient.GetDel(context.Background(), fmt.Sprintf(CAPTCHA, captchaStr)).Err()
	if err == nil {
		return true
	} else if err == redis.Nil {
		return false
	} else {
		log.Error(err, errors.New("verify captcha error"))
	}
	return false
}

func GenerateCaptchaImgBuffer(log logx.Logger, redisClient *redis.Client, CAPTCHA string, width, height int) (*bytes.Buffer, error) {
	// 避免重复
	var digits []byte
	for {
		digits = GetRandomDigits(6)
		err := redisClient.Get(context.Background(), fmt.Sprintf(CAPTCHA, BytesNumsToStringNums(digits))).Err()
		if err == redis.Nil {
			break
		} else if err != nil {
			log.Error(errors.Wrap(err, "redis获取数据异常"))
			return nil, errors.New("redis获取数据异常")
		}
	}
	// 生成验证码图片缓存
	return WriteCaptchaImage(redisClient, CAPTCHA, digits, width, height)
}
