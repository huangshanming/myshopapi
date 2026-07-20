package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type RecycleRestore2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecycleRestore2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleRestore2Logic {
	return &RecycleRestore2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleRestore2Logic) RecycleRestore2(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).RecycleRestore(w, r)
}
