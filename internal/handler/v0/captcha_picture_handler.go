package v0

import (
	"net/http"

	"cloud-platform-system/internal/logic/v0"
	"cloud-platform-system/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CaptchaPictureHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := v0.NewCaptchaPictureLogic(r.Context(), svcCtx)
		resp, err := l.CaptchaPicture()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
