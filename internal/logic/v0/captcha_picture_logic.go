package v0

import (
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaPictureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaPictureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaPictureLogic {
	return &CaptchaPictureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaPictureLogic) CaptchaPicture() (resp *types.CaptchaPictureResponse, err error) {
	// 生成验证码图片
	buff, err := utils.GenerateCaptchaImgBuffer(l.Logger, l.svcCtx.RedisClient, l.svcCtx.CAPTCHA, l.svcCtx.Config.Captcha.Width, l.svcCtx.Config.Captcha.Height)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "can not generate captcha"))
		return nil, errors.New("系统错误，生成验证码失败！")
	}
	return &types.CaptchaPictureResponse{PicData: buff.Bytes()}, nil
}
