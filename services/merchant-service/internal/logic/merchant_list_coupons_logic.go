package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListCouponsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListCouponsLogic {
	return &MerchantListCouponsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListCouponsLogic) MerchantListCoupons(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewCouponHandler(l.svcCtx).MerchantListCoupons(w, r)
}
