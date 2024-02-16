package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/redis/go-redis/v9"
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
	for i := 0; i < len(digits); i++ {
		digits[i] += '0'
	}
	return string(digits)
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

func IsValidCaptcha(redisClient *redis.Client, CAPTCHA string, captchaStr string) bool {
	err := redisClient.Del(context.Background(), fmt.Sprintf(CAPTCHA, captchaStr)).Err()
	if err == nil {
		return true
	} else if err == redis.Nil {
		return false
	}
	return false
}
