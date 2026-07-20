package handler

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/svc"
)

func UploadMineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUploadMineLogic(r.Context(), svcCtx)
		l.UploadMine(w, r)
	}
}
