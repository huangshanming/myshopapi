package health

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/public/health"
	"mymall/services/catalog-service/internal/svc"
)

func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := health.NewHealthzLogic(r.Context(), svcCtx)
		l.Healthz(w, r)
	}
}

func MetricsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := health.NewMetricsLogic(r.Context(), svcCtx)
		l.Metrics(w, r)
	}
}

func ReadyzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := health.NewReadyzLogic(r.Context(), svcCtx)
		l.Readyz(w, r)
	}
}
