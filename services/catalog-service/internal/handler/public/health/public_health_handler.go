package health

import (
	"net/http"

	"mymall/pkg/metrics"

	"mymall/pkg/httpserver"

	"mymall/services/catalog-service/internal/svc"
)

func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpserver.Healthz("catalog-service")(w, r)
	}
}

func MetricsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.Handler()(w, r)
	}
}

func ReadyzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svcCtx.Health.ReadyHandler()(w, r)
	}
}
