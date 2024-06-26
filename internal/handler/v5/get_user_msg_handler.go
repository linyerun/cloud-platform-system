package v5

import (
	"net/http"

	"cloud-platform-system/internal/logic/v5"
	"cloud-platform-system/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := v5.NewGetUserMsgLogic(r.Context(), svcCtx)
		resp, err := l.GetUserMsg()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
