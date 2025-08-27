package user

import (
	"net/http"

	"cdp-admin-service/internal/logic/user"
	"cdp-admin-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserRegionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserRegionsLogic(r.Context(), svcCtx)
		resp, err := l.GetUserRegions()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
