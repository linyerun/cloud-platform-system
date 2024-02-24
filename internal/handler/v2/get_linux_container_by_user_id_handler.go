package v2

import (
	"net/http"

	"cloud-platform-system/internal/logic/v2"
	"cloud-platform-system/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetLinuxContainerByUserIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := v2.NewGetLinuxContainerByUserIdLogic(r.Context(), svcCtx)
		resp, err := l.GetLinuxContainerByUserId()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
