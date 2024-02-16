package v0

import (
	"net/http"

	"cloud-platform-system/internal/logic/v0"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CaptchaEmailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CaptchaEmailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := v0.NewCaptchaEmailLogic(r.Context(), svcCtx)
		resp, err := l.CaptchaEmail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
