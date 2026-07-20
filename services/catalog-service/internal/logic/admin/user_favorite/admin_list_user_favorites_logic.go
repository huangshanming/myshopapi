package user_favorite

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminListUserFavoritesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListUserFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListUserFavoritesLogic {
	return &AdminListUserFavoritesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListUserFavoritesLogic) AdminListUserFavorites(w http.ResponseWriter, r *http.Request) {
	padmin.NewFavoriteHandler(l.svcCtx).AdminUserList(w, r)
}
