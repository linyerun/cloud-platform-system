package middleware

import (
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type JwtAuthMiddleware struct {
}

func NewJwtAuthMiddleware() *JwtAuthMiddleware {
	return &JwtAuthMiddleware{}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		obj := new(utils.DefaultTokenObject)
		err := utils.ParseToken(token, obj)
		if err != nil {
			logx.WithContext(context.Background()).Error(errors.Wrap(err, "parse token error"))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(types.CommonResponse{Code: 500, Msg: err.Error()})
			return
		}
		if !obj.IsValid() {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(types.CommonResponse{Code: 400, Msg: "token已过期"})
			return
		}

		// 讲数据保存到request中
		newReq := r.WithContext(context.WithValue(r.Context(), "user", obj.Claims))
		next(w, newReq)
	}
}
