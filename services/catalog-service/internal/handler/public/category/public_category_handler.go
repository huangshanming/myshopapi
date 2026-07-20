package category

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/public/category"
	"mymall/services/catalog-service/internal/svc"
)

func GetCategoryDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewGetCategoryDetailLogic(r.Context(), svcCtx)
		l.GetCategoryDetail(w, r)
	}
}

func GetCategoryListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewGetCategoryListLogic(r.Context(), svcCtx)
		l.GetCategoryList(w, r)
	}
}
