package utils

import (
	"bytes"
	"cloud-platform-system/internal/common/errorx"
	"cloud-platform-system/internal/config"
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
func WriteCaptchaImage(digits []byte, width, height int) (buff *bytes.Buffer, err error) {
	buff = bytes.NewBuffer(nil)
	_, err = captcha.NewImage("", digits, width, height).WriteTo(buff)
	if err != nil {
		return nil, errors.Wrapf(err, "generate captcha error")
	}
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

func GenerateCaptchaImgBuffer(log logx.Logger, redisClient *redis.Client, CAPTCHA string, width, height int, timeoutSec uint) (*bytes.Buffer, error) {
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
	buff, err := WriteCaptchaImage(digits, width, height)
	if err != nil {
		return nil, err
	}

	// 把信息保存到redis
	err = redisClient.Set(context.Background(), fmt.Sprintf(CAPTCHA, BytesNumsToStringNums(digits)), "", time.Second*time.Duration(timeoutSec)).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "save data to redis error")
	}

	return buff, nil
}

func IsValidEmailCaptcha(redisClient *redis.Client, CAPTCHA, captchaStr, email string) error {
	err := redisClient.GetDel(context.Background(), fmt.Sprintf(CAPTCHA, captchaStr+"-"+email)).Err()

	if err == redis.Nil {
		return errorx.NewBaseError(400, "no this captcha")
	} else if err != nil {
		return errorx.NewBaseError(500, "delete redis data error")
	}

	return nil
}

func GenerateCaptchaImgBufferByEmail(redisClient *redis.Client, c config.Config, CAPTCHA, email string) (*bytes.Buffer, error) {
	// 避免重复
	var digits []byte
	for {
		digits = GetRandomDigits(6)
		err := redisClient.Get(context.Background(), fmt.Sprintf(CAPTCHA, BytesNumsToStringNums(digits)+"-"+email)).Err()
		if err == redis.Nil {
			break
		} else if err != nil {
			return nil, errors.Wrapf(err, "redis获取数据异常")
		}
	}

	// 生成验证码图片
	buff, err := WriteCaptchaImage(digits, c.Captcha.Width, c.Captcha.Height)
	if err != nil {
		return nil, err
	}

	// 把信息保存到redis
	err = redisClient.Set(context.Background(), fmt.Sprintf(CAPTCHA, BytesNumsToStringNums(digits)+"-"+email), "", time.Second*time.Duration(c.Captcha.TimeoutSec)).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "save data to redis error")
	}

	return buff, nil
}
