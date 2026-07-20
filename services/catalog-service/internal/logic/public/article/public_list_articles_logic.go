package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type PublicListArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicListArticlesLogic {
	return &PublicListArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicListArticlesLogic) PublicListArticles(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).List(w, r)
}
