package v0

import (
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaEmailLogic {
	return &CaptchaEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaEmailLogic) CaptchaEmail(req *types.CaptchaEmailRequest) (resp *types.CommonResponse, err error) {
	// 校验邮箱
	if !utils.IsNormalEmail(req.Email) {
		return &types.CommonResponse{Code: 400, Msg: "邮箱地址有问题"}, nil
	}

	// 生成验证码图片
	buff, err := utils.GenerateCaptchaImgBuffer(l.Logger, l.svcCtx.RedisClient, l.svcCtx.CAPTCHA, l.svcCtx.Config.Captcha.Width, l.svcCtx.Config.Captcha.Height)
	if err != nil {
		l.Logger.Error(errors.Wrap(err, "can not generate captcha"))
		return &types.CommonResponse{Code: 500, Msg: "系统错误，生成验证码失败！"}, nil
	}

	// 发送
	if err = utils.SendCaptchaByEmail(req.Email, buff); err != nil {
		l.Logger.Error(errors.Wrap(err, "send email error"))
		return &types.CommonResponse{Code: 500, Msg: "系统错误，发送验证码失败！"}, nil
	}
	return &types.CommonResponse{Code: 200, Msg: "发送成功"}, nil
}
