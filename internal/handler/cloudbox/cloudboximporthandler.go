package cloudbox

import (
	"net/http"

	"cdp-admin-service/internal/logic/cloudbox"
	"cdp-admin-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CloudBoxImportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := cloudbox.NewCloudBoxImportLogic(r.Context(), svcCtx)
		resp, err := l.CloudBoxImport(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
