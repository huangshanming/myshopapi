package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantCouponClaimsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantCouponClaimsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCouponClaimsLogic {
	return &MerchantCouponClaimsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantCouponClaimsLogic) MerchantCouponClaims(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewCouponHandler(l.svcCtx).MerchantCouponClaims(w, r)
}
