package product

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type MerchantPurgeProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantPurgeProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantPurgeProductsLogic {
	return &MerchantPurgeProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantPurgeProductsLogic) MerchantPurgeProducts(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).RecycleDelete(w, r)
}
