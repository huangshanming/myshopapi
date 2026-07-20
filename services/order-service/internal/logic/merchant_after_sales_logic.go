package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/order-service/internal/httpapi/merchant"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantAfterSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantAfterSalesLogic {
	return &MerchantAfterSalesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantAfterSalesLogic) MerchantAfterSales(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewOrderHandler(l.svcCtx).MerchantAfterSales(w, r)
}
