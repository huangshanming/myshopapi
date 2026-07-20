package handler

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/svc"
)

func SaveAttrTemplate2Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewSaveAttrTemplate2Logic(r.Context(), svcCtx)
		l.SaveAttrTemplate2(w, r)
	}
}
