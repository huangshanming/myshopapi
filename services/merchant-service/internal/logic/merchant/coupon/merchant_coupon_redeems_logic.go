package coupon

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MerchantCouponRedeemsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantCouponRedeemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCouponRedeemsLogic {
	return &MerchantCouponRedeemsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantCouponRedeemsLogic) MerchantCouponRedeems(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewCouponHandler(l.svcCtx).MerchantCouponRedeems(w, r)
}
