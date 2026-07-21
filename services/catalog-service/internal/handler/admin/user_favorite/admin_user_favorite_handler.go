package user_favorite

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/services/catalog-service/internal/logic/admin/user_favorite"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"
)

func AdminListUserFavoritesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user_favorite.NewAdminListUserFavoritesLogic(svcCtx)
		resp, err := l.AdminListUserFavorites(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
