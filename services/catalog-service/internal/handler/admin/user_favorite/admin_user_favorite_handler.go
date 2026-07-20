package user_favorite

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/admin/user_favorite"
	"mymall/services/catalog-service/internal/svc"
)

func AdminListUserFavoritesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user_favorite.NewAdminListUserFavoritesLogic(r.Context(), svcCtx)
		l.AdminListUserFavorites(w, r)
	}
}
