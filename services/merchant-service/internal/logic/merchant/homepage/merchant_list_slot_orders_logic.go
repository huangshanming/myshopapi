package homepage

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MerchantListSlotOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListSlotOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListSlotOrdersLogic {
	return &MerchantListSlotOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListSlotOrdersLogic) MerchantListSlotOrders(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewHomepageSlotHandler(l.svcCtx).MerchantListSlotOrders(w, r)
}
