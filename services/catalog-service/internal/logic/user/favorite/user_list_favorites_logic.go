package favorite

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	puser "mymall/services/catalog-service/internal/product/httpapi/user"
	"mymall/services/catalog-service/internal/svc"
)

type UserListFavoritesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListFavoritesLogic {
	return &UserListFavoritesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListFavoritesLogic) UserListFavorites(w http.ResponseWriter, r *http.Request) {
	puser.NewFavoriteHandler(l.svcCtx).List(w, r)
}
