package instance

import (
	"net/http"

	"cdp-admin-service/internal/logic/instance"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func StrategyInstNoScheduingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StrategyInstNoScheduingReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := instance.NewStrategyInstNoScheduingLogic(r.Context(), svcCtx)
		resp, err := l.StrategyInstNoScheduing(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
