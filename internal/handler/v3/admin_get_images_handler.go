package v3

import (
	"net/http"

	"cloud-platform-system/internal/logic/v3"
	"cloud-platform-system/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminGetImagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := v3.NewAdminGetImagesLogic(r.Context(), svcCtx)
		resp, err := l.AdminGetImages()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
