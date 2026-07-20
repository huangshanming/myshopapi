package logic

import (
	"net/http"

	"context"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicThemeTilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicThemeTilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicThemeTilesLogic {
	return &PublicThemeTilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicThemeTilesLogic) PublicThemeTiles(w http.ResponseWriter, r *http.Request) {
	hpublic.NewHomepageThemeHandler(l.svcCtx).PublicThemeTiles(w, r)
}
