package handler

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/svc"
)

func SaveAttrTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewSaveAttrTemplateLogic(r.Context(), svcCtx)
		l.SaveAttrTemplate(w, r)
	}
}
