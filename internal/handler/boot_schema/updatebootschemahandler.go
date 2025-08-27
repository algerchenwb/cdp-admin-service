package boot_schema

import (
	"net/http"

	"cdp-admin-service/internal/logic/boot_schema"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateBootSchemaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateBootSchemaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := boot_schema.NewUpdateBootSchemaLogic(r.Context(), svcCtx)
		resp, err := l.UpdateBootSchema(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
