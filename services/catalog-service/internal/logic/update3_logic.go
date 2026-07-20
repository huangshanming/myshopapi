package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type Update3Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdate3Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Update3Logic {
	return &Update3Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Update3Logic) Update3(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).Update(w, r)
}
