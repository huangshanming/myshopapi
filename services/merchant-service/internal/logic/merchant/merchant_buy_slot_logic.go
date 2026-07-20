package merchant

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MerchantBuySlotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantBuySlotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantBuySlotLogic {
	return &MerchantBuySlotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantBuySlotLogic) MerchantBuySlot(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewHomepageSlotHandler(l.svcCtx).MerchantBuySlot(w, r)
}
