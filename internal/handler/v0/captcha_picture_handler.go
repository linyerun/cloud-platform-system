package v0

import (
	"cloud-platform-system/internal/types"
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
			httpx.OkJsonCtx(r.Context(), w, &types.CommonResponse{Code: 500, Msg: err.Error()})
		} else {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			w.Header().Set("Content-Type", "image/png")
			w.Write(resp.PicData)
		}
	}
}
