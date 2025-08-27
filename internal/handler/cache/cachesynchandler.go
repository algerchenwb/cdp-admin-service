package cache

import (
	"net/http"

	"cdp-admin-service/internal/logic/cache"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CacheSyncHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cache.NewCacheSyncLogic(r.Context(), svcCtx)
		err := l.CacheSync()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.CommonNilJson{})
		}
	}
}
