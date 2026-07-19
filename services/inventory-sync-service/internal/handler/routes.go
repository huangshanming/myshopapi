package handler

import (
	"net/http"

	"mymall/pkg/metrics"
	"mymall/services/inventory-sync-service/internal/logic"
	"mymall/services/inventory-sync-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext, wrap func(http.HandlerFunc) http.HandlerFunc) {
	if wrap == nil {
		wrap = func(h http.HandlerFunc) http.HandlerFunc { return h }
	}
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: wrap(HealthzHandler(serverCtx))},
		{Method: http.MethodGet, Path: "/readyz", Handler: wrap(ReadyzHandler(serverCtx))},
		{Method: http.MethodGet, Path: "/metrics", Handler: wrap(metrics.Handler())},
	})
}

func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewHealthzLogic(r.Context(), svcCtx)
		resp, err := l.Healthz()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func ReadyzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewReadyzLogic(r.Context(), svcCtx)
		resp, err := l.Readyz()
		if err != nil {
			// 仍返回检查详情
			if resp != nil {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusServiceUnavailable, resp)
				return
			}
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
