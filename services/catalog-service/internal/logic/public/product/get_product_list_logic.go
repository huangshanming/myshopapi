package product

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	ppublic "mymall/services/catalog-service/internal/product/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type GetProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductListLogic {
	return &GetProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductListLogic) GetProductList(w http.ResponseWriter, r *http.Request) {
	ppublic.NewCatalogHandler(l.svcCtx).GetProductList(w, r)
}
