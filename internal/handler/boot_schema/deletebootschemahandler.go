package boot_schema

import (
	"net/http"

	"cdp-admin-service/internal/logic/boot_schema"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteBootSchemaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteBootSchemaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := boot_schema.NewDeleteBootSchemaLogic(r.Context(), svcCtx)
		resp, err := l.DeleteBootSchema(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
