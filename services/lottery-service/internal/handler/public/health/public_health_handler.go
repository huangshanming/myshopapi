package health

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/xerr"
	"mymall/services/lottery-service/internal/svc"
	"mymall/services/lottery-service/internal/types"
)

func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.OkJsonCtx(r.Context(), w, types.EmptyResp{})
	}
}

func ReadyzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svcCtx.Health != nil {
			ok, _ := svcCtx.Health.CheckAll(r.Context())
			if !ok {
				httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusServiceUnavailable, "not ready"))
				return
			}
		}
		httpx.OkJsonCtx(r.Context(), w, types.EmptyResp{})
	}
}
