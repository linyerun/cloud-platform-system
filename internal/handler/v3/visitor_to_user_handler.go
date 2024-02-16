package v3

import (
	"net/http"

	"cloud-platform-system/internal/logic/v3"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func VisitorToUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutVisitorToUserRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := v3.NewVisitorToUserLogic(r.Context(), svcCtx)
		resp, err := l.VisitorToUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
