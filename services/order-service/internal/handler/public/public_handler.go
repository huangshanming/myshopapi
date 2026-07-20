package public

import (
	"net/http"

	"mymall/services/order-service/internal/logic/public"
	"mymall/services/order-service/internal/svc"
)

func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewHealthzLogic(r.Context(), svcCtx)
		l.Healthz(w, r)
	}
}

func MetricsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewMetricsLogic(r.Context(), svcCtx)
		l.Metrics(w, r)
	}
}

func PublicListProductReviewsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicListProductReviewsLogic(r.Context(), svcCtx)
		l.PublicListProductReviews(w, r)
	}
}

func ReadyzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewReadyzLogic(r.Context(), svcCtx)
		l.Readyz(w, r)
	}
}
