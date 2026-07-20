package handler

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/svc"
)

func GetCategoryDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetCategoryDetailLogic(r.Context(), svcCtx)
		l.GetCategoryDetail(w, r)
	}
}
