package order

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/order-service/internal/httpapi/merchant"
	"mymall/services/order-service/internal/svc"
)

type MerchantDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantDetailLogic {
	return &MerchantDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantDetailLogic) MerchantDetail(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewOrderHandler(l.svcCtx).MerchantDetail(w, r)
}
