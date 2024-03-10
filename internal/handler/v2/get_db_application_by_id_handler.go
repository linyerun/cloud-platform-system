package v2

import (
	"net/http"

	"cloud-platform-system/internal/logic/v2"
	"cloud-platform-system/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDbApplicationByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := v2.NewGetDbApplicationByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetDbApplicationById()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
