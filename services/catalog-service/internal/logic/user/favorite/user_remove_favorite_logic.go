package favorite

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	puser "mymall/services/catalog-service/internal/product/httpapi/user"
	"mymall/services/catalog-service/internal/svc"
)

type UserRemoveFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRemoveFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRemoveFavoriteLogic {
	return &UserRemoveFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRemoveFavoriteLogic) UserRemoveFavorite(w http.ResponseWriter, r *http.Request) {
	puser.NewFavoriteHandler(l.svcCtx).Remove(w, r)
}
