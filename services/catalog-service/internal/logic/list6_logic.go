package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type List6Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList6Logic(ctx context.Context, svcCtx *svc.ServiceContext) *List6Logic {
	return &List6Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List6Logic) List6(w http.ResponseWriter, r *http.Request) {
	padmin.NewPlatformProductHandler(l.svcCtx).List(w, r)
}
