package product

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type MerchantRestoreProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantRestoreProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantRestoreProductsLogic {
	return &MerchantRestoreProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantRestoreProductsLogic) MerchantRestoreProducts(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).RecycleRestore(w, r)
}
