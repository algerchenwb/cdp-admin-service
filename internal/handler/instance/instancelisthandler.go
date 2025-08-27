package instance

import (
	"context"
	"net/http"

	"cdp-admin-service/internal/logic/instance"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func InstanceListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommonPageRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(context.Background()).Errorf("httpx.Parse  err: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := instance.NewInstanceListLogic(r.Context(), svcCtx)
		resp, err := l.InstanceList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
