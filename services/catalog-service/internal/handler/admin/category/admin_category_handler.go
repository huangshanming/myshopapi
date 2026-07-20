package category

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/admin/category"
	"mymall/services/catalog-service/internal/svc"
)

func AdminCreateCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewAdminCreateCategoryLogic(r.Context(), svcCtx)
		l.AdminCreateCategory(w, r)
	}
}

func AdminDeleteCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewAdminDeleteCategoryLogic(r.Context(), svcCtx)
		l.AdminDeleteCategory(w, r)
	}
}

func AdminListCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewAdminListCategoriesLogic(r.Context(), svcCtx)
		l.AdminListCategories(w, r)
	}
}

func AdminUpdateCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewAdminUpdateCategoryLogic(r.Context(), svcCtx)
		l.AdminUpdateCategory(w, r)
	}
}
