package favorite

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	puser "mymall/services/catalog-service/internal/product/httpapi/user"
	"mymall/services/catalog-service/internal/svc"
)

type UserAddFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserAddFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAddFavoriteLogic {
	return &UserAddFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserAddFavoriteLogic) UserAddFavorite(w http.ResponseWriter, r *http.Request) {
	puser.NewFavoriteHandler(l.svcCtx).Add(w, r)
}
