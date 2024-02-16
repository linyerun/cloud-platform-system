package v3

import (
	"net/http"

	"cloud-platform-system/internal/logic/v3"
	"cloud-platform-system/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := v3.NewGetUsersLogic(r.Context(), svcCtx)
		resp, err := l.GetUsers()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
