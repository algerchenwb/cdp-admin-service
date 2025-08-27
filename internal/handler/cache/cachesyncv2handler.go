package cache

import (
	"net/http"

	"cdp-admin-service/internal/logic/cache"
	"cdp-admin-service/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CacheSyncV2Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cache.NewCacheSyncV2Logic(r.Context(), svcCtx)
		resp, err := l.CacheSyncV2()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
