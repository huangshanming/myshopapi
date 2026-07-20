package handler

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

func PublicThemeTilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewPublicThemeTilesLogic(r.Context(), svcCtx)
		l.PublicThemeTiles(w, r)
	}
}
