package v4

import (
	"net/http"

	"cloud-platform-system/internal/logic/v4"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DelExceptionByIdxHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelExceptionByIdxReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := v4.NewDelExceptionByIdxLogic(r.Context(), svcCtx)
		resp, err := l.DelExceptionByIdx(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
