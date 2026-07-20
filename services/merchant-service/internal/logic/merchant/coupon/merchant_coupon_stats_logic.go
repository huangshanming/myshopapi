package coupon

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MerchantCouponStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantCouponStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCouponStatsLogic {
	return &MerchantCouponStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantCouponStatsLogic) MerchantCouponStats(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewCouponHandler(l.svcCtx).MerchantCouponStats(w, r)
}
