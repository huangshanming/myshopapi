package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type List3Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList3Logic(ctx context.Context, svcCtx *svc.ServiceContext) *List3Logic {
	return &List3Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List3Logic) List3(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).List(w, r)
}
