package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantCreateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCreateCouponLogic {
	return &MerchantCreateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateCouponLogic) MerchantCreateCoupon(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewCouponHandler(l.svcCtx).MerchantCreateCoupon(w, r)
}
