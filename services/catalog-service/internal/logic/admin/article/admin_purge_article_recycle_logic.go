package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminPurgeArticleRecycleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminPurgeArticleRecycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminPurgeArticleRecycleLogic {
	return &AdminPurgeArticleRecycleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminPurgeArticleRecycleLogic) AdminPurgeArticleRecycle(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).RecycleDelete(w, r)
}
