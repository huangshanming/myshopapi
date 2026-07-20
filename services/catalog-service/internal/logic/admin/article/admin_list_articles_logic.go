package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminListArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListArticlesLogic {
	return &AdminListArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListArticlesLogic) AdminListArticles(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).List(w, r)
}
