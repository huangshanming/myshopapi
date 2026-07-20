package handler

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/svc"
)

func SaveTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewSaveTagLogic(r.Context(), svcCtx)
		l.SaveTag(w, r)
	}
}
