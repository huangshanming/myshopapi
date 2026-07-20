package public

import (
	"net/http"

	"mymall/services/user-service/internal/logic/public"
	"mymall/services/user-service/internal/svc"
)

func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewHealthzLogic(r.Context(), svcCtx)
		l.Healthz(w, r)
	}
}

func ListRegionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewListRegionsLogic(r.Context(), svcCtx)
		l.ListRegions(w, r)
	}
}

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewLoginLogic(r.Context(), svcCtx)
		l.Login(w, r)
	}
}

func MetricsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewMetricsLogic(r.Context(), svcCtx)
		l.Metrics(w, r)
	}
}

func ReadyzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewReadyzLogic(r.Context(), svcCtx)
		l.Readyz(w, r)
	}
}

func RegionTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewRegionTreeLogic(r.Context(), svcCtx)
		l.RegionTree(w, r)
	}
}

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewRegisterLogic(r.Context(), svcCtx)
		l.Register(w, r)
	}
}
