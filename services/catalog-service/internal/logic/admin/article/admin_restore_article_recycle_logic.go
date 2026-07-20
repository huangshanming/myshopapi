package article

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminRestoreArticleRecycleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminRestoreArticleRecycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminRestoreArticleRecycleLogic {
	return &AdminRestoreArticleRecycleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminRestoreArticleRecycleLogic) AdminRestoreArticleRecycle(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).RecycleRestore(w, r)
}
