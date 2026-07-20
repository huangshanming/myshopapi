package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantGrantCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantGrantCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantGrantCouponLogic {
	return &MerchantGrantCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantGrantCouponLogic) MerchantGrantCoupon(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewCouponHandler(l.svcCtx).MerchantGrantCoupon(w, r)
}
