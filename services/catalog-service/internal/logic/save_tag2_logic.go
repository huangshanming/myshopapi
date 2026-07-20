package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type SaveTag2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveTag2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveTag2Logic {
	return &SaveTag2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveTag2Logic) SaveTag2(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).SaveTag(w, r)
}
