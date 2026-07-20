package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cpublic "mymall/services/catalog-service/internal/content/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type PublicGetArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicGetArticleLogic {
	return &PublicGetArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicGetArticleLogic) PublicGetArticle(w http.ResponseWriter, r *http.Request) {
	cpublic.NewArticleHandler(l.svcCtx).Detail(w, r)
}
