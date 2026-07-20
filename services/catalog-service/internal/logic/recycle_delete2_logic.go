package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type RecycleDelete2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecycleDelete2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleDelete2Logic {
	return &RecycleDelete2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleDelete2Logic) RecycleDelete2(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).RecycleDelete(w, r)
}
