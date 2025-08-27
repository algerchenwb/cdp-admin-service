package biz

import (
	"context"
	"net/http"

	"cdp-admin-service/internal/logic/biz"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConfigBizStrategyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConfigBizStrategyReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(context.Background()).Errorf("httpx.Parse  err: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := biz.NewConfigBizStrategyLogic(r.Context(), svcCtx)
		resp, err := l.ConfigBizStrategy(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
