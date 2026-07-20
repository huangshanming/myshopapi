package product

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/admin/product"
	"mymall/services/catalog-service/internal/svc"
)

func AdminDeleteProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewAdminDeleteProductLogic(r.Context(), svcCtx)
		l.AdminDeleteProduct(w, r)
	}
}

func AdminListProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewAdminListProductsLogic(r.Context(), svcCtx)
		l.AdminListProducts(w, r)
	}
}

func AdminOffSaleProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewAdminOffSaleProductLogic(r.Context(), svcCtx)
		l.AdminOffSaleProduct(w, r)
	}
}
