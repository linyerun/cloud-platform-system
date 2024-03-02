package common

import (
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func GetToken(user *models.User, log logx.Logger, svcCtx *svc.ServiceContext) (*types.CommonResponse, error) {
	// 把信息封装到token中并返回给用户
	obj, err := utils.NewDefaultTokenObject(user, time.Duration(svcCtx.Config.Jwt.ExpireSec)*time.Second)
	if err != nil {
		log.Error(errors.Wrap(err, "utils.NewDefaultTokenObject error"))
		return &types.CommonResponse{Code: 500, Msg: err.Error()}, nil
	}
	token, err := obj.GenerateToken()
	if err != nil {
		log.Error(errors.Wrap(err, "generate token error"))
		return &types.CommonResponse{Code: 500, Msg: err.Error()}, nil
	}

	// 生成refresh token
	obj, err = utils.NewDefaultTokenObject(user, time.Duration(svcCtx.Config.Jwt.ExpireSec*2)*time.Second)
	if err != nil {
		log.Error(errors.Wrap(err, "utils.NewDefaultTokenObject error"))
		return &types.CommonResponse{Code: 500, Msg: err.Error()}, nil
	}
	refreshToken, err := obj.GenerateToken()
	if err != nil {
		log.Error(errors.Wrap(err, "generate token error"))
		return &types.CommonResponse{Code: 500, Msg: err.Error()}, nil
	}

	return &types.CommonResponse{
		Code: 200,
		Msg:  "成功",
		Data: map[string]any{
			"token":         token,
			"refresh_token": refreshToken,
			"user": map[string]any{
				"auth":  user.Auth,
				"email": user.Email,
				"name":  user.Name,
			},
		},
	}, nil
}
