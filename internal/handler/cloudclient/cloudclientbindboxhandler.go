package cloudclient

import (
	"net/http"

	"cdp-admin-service/internal/logic/cloudclient"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CloudClientBindBoxHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClientBindBoxReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cloudclient.NewCloudClientBindBoxLogic(r.Context(), svcCtx)
		resp, err := l.CloudClientBindBox(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
