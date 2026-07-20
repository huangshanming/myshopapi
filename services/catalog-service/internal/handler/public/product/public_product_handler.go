package product

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/public/product"
	"mymall/services/catalog-service/internal/svc"
)

func CountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewCountLogic(r.Context(), svcCtx)
		l.Count(w, r)
	}
}

func GetProductDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewGetProductDetailLogic(r.Context(), svcCtx)
		l.GetProductDetail(w, r)
	}
}

func GetProductListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewGetProductListLogic(r.Context(), svcCtx)
		l.GetProductList(w, r)
	}
}

func GetSalesRankHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewGetSalesRankLogic(r.Context(), svcCtx)
		l.GetSalesRank(w, r)
	}
}
