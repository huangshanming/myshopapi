package favorite

import (
	"net/http"

	"mymall/services/catalog-service/internal/logic/user/favorite"
	"mymall/services/catalog-service/internal/svc"
)

func RemoveBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := favorite.NewRemoveBatchLogic(r.Context(), svcCtx)
		l.RemoveBatch(w, r)
	}
}

func UserAddFavoriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := favorite.NewUserAddFavoriteLogic(r.Context(), svcCtx)
		l.UserAddFavorite(w, r)
	}
}

func UserFavoriteStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := favorite.NewUserFavoriteStatusLogic(r.Context(), svcCtx)
		l.UserFavoriteStatus(w, r)
	}
}

func UserListFavoritesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := favorite.NewUserListFavoritesLogic(r.Context(), svcCtx)
		l.UserListFavorites(w, r)
	}
}

func UserRemoveFavoriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := favorite.NewUserRemoveFavoriteLogic(r.Context(), svcCtx)
		l.UserRemoveFavorite(w, r)
	}
}
