package middleware

import (
	"cloud-platform-system/internal/resp"
	"cloud-platform-system/internal/types"
	"cloud-platform-system/internal/utils"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
	"time"
)

type JwtAuthMiddleware struct {
	lock              sync.Mutex
	userUpdateLockMap map[string]chan struct{}
}

func NewJwtAuthMiddleware() *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		userUpdateLockMap: make(map[string]chan struct{}),
	}
}

func (m *JwtAuthMiddleware) getUserLock(id string) (ch chan struct{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	ch, ok := m.userUpdateLockMap[id]
	if !ok { // 获取不到锁就创建一个
		ch = make(chan struct{}, 1)
		m.userUpdateLockMap[id] = ch
	}
	return ch
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		obj := new(utils.DefaultTokenObject)
		err := utils.ParseToken(token, obj)
		if err != nil {
			logx.WithContext(context.Background()).Error(errors.Wrap(err, "parse token error"))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(types.CommonResponse{Code: resp.TokenParseError, Msg: resp.MsgMap[resp.TokenParseError]})
			return
		}
		if !obj.IsValid() {
			w.Header().Set("Content-Type", "application/json")
			// Token过期的必须返回一个特殊的Code来方便前端进行判断
			json.NewEncoder(w).Encode(types.CommonResponse{Code: resp.TokenInValidError, Msg: resp.MsgMap[resp.TokenInValidError]})
			return
		}

		// 讲数据保存到request中
		newReq := r.WithContext(context.WithValue(r.Context(), "user", obj.Claims))

		// 这里使用锁保证用户修改操作同时刻只能是被使用一次
		method := r.Method
		if method != http.MethodDelete && method != http.MethodPut && method != http.MethodPost { // 如果是修改请求，直接放行
			next(w, newReq)
			return
		}

		// 获取锁
		ch := m.getUserLock(obj.Claims.Id)

		// 5分钟后获取不到锁就放弃
		select {
		case ch <- struct{}{}:
		case <-time.After(5 * time.Minute):
			// 返回修改失败结果
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(types.CommonResponse{Code: 401, Msg: "执行改变数据的操作失败"})
			return
		}

		// 执行修改逻辑
		next(w, newReq)

		// 释放锁
		<-ch
	}
}
