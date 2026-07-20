package shop

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/admin/shop"
	"mymall/services/catalog-service/internal/svc"
)

func AdminUploadShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewAdminUploadShopLogic(r.Context(), svcCtx)
		l.AdminUploadShop(w, r)
	}
}
