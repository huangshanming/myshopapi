package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type List8Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList8Logic(ctx context.Context, svcCtx *svc.ServiceContext) *List8Logic {
	return &List8Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List8Logic) List8(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).List(w, r)
}
