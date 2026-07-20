package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	puser "mymall/services/catalog-service/internal/product/httpapi/user"
	"mymall/services/catalog-service/internal/svc"
)

type List2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *List2Logic {
	return &List2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List2Logic) List2(w http.ResponseWriter, r *http.Request) {
	puser.NewFavoriteHandler(l.svcCtx).List(w, r)
}
