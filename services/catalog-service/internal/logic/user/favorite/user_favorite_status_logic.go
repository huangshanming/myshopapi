package favorite

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	puser "mymall/services/catalog-service/internal/product/httpapi/user"
	"mymall/services/catalog-service/internal/svc"
)

type UserFavoriteStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFavoriteStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFavoriteStatusLogic {
	return &UserFavoriteStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFavoriteStatusLogic) UserFavoriteStatus(w http.ResponseWriter, r *http.Request) {
	puser.NewFavoriteHandler(l.svcCtx).Status(w, r)
}
