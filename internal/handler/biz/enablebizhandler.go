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

func EnableBizHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnableBizReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(context.Background()).Errorf("httpx.Parse  err: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := biz.NewEnableBizLogic(r.Context(), svcCtx)
		resp, err := l.EnableBiz(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
