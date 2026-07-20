package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type ListMyFavoritesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyFavoritesLogic {
	return &ListMyFavoritesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyFavoritesLogic) ListMyFavorites(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).ListMyFavorites(w, r)
}
