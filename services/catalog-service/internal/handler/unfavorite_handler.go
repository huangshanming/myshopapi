package handler

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/svc"
)

func UnfavoriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUnfavoriteLogic(r.Context(), svcCtx)
		l.Unfavorite(w, r)
	}
}
