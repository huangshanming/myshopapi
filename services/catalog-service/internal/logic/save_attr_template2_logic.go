package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type SaveAttrTemplate2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveAttrTemplate2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAttrTemplate2Logic {
	return &SaveAttrTemplate2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveAttrTemplate2Logic) SaveAttrTemplate2(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).SaveAttrTemplate(w, r)
}
