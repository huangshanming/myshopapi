package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminBatchAuditArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminBatchAuditArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminBatchAuditArticlesLogic {
	return &AdminBatchAuditArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminBatchAuditArticlesLogic) AdminBatchAuditArticles(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).BatchAudit(w, r)
}
