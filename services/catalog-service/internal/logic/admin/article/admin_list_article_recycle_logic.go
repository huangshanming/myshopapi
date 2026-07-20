package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminListArticleRecycleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListArticleRecycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListArticleRecycleLogic {
	return &AdminListArticleRecycleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListArticleRecycleLogic) AdminListArticleRecycle(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).List(w, r)
}
