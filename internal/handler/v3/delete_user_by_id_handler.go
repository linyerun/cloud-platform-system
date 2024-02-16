package v3

import (
	"net/http"

	"cloud-platform-system/internal/logic/v3"
	"cloud-platform-system/internal/svc"
	"cloud-platform-system/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteUserByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteUserRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := v3.NewDeleteUserByIdLogic(r.Context(), svcCtx)
		resp, err := l.DeleteUserById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
