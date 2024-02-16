package v1

import (
	"net/http"

	"cloud-platform-system/internal/logic/v1"
	"cloud-platform-system/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAdminsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := v1.NewGetAdminsLogic(r.Context(), svcCtx)
		resp, err := l.GetAdmins()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
