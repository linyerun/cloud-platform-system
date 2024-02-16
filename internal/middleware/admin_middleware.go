package middleware

import (
	"cloud-platform-system/internal/models"
	"cloud-platform-system/internal/types"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type AdminMiddleware struct {
}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (m *AdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*models.User)
		if !ok {
			panic(errors.Errorf("token  parse user change to *models.User error"))
		}
		if user.Auth != models.AdminAuth {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(types.CommonResponse{Code: 400, Msg: "权限不足，无法访问！"})
			return
		}
		next(w, r)
	}
}
