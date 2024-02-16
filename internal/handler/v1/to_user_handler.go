package v1

import (
	"net/http"

	"cloud-platform-system/internal/logic/v1"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ToUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplicationFormPostRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := v1.NewToUserLogic(r.Context(), svcCtx)
		resp, err := l.ToUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
