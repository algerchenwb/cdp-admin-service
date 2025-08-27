package instance

import (
	"net/http"

	"cdp-admin-service/internal/logic/instance"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func InstancePoolAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InstancePoolAddReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := instance.NewInstancePoolAddLogic(r.Context(), svcCtx)
		resp, err := l.InstancePoolAdd(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
