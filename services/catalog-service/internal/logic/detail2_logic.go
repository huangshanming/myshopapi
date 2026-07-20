package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type Detail2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetail2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Detail2Logic {
	return &Detail2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Detail2Logic) Detail2(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).Detail(w, r)
}
