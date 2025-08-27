package user

import (
	"net/http"

	"cdp-admin-service/internal/logic/user"
	"cdp-admin-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserPermMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserPermMenuLogic(r.Context(), svcCtx)
		resp, err := l.GetUserPermMenu()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
