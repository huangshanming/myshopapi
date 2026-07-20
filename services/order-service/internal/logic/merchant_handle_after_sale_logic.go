package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/order-service/internal/httpapi/merchant"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantHandleAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantHandleAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantHandleAfterSaleLogic {
	return &MerchantHandleAfterSaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantHandleAfterSaleLogic) MerchantHandleAfterSale(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewOrderHandler(l.svcCtx).MerchantHandleAfterSale(w, r)
}
