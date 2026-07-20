package region

import (
	"net/http"

	"mymall/services/user-service/internal/logic/public/region"
	"mymall/services/user-service/internal/svc"
)

func ListRegionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := region.NewListRegionsLogic(r.Context(), svcCtx)
		l.ListRegions(w, r)
	}
}

func RegionTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := region.NewRegionTreeLogic(r.Context(), svcCtx)
		l.RegionTree(w, r)
	}
}
