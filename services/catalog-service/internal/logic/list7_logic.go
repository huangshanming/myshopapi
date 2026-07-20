package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type List7Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList7Logic(ctx context.Context, svcCtx *svc.ServiceContext) *List7Logic {
	return &List7Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List7Logic) List7(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).List(w, r)
}
