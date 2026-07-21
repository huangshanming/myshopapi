package coupon

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hmerchant "mymall/services/merchant-service/internal/app/merchant"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantGrantCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantGrantCouponLogic(svcCtx *svc.ServiceContext) *MerchantGrantCouponLogic {
	return &MerchantGrantCouponLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantGrantCouponLogic) MerchantGrantCoupon(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/coupons/grant", nil, nil, req, hmerchant.NewCouponHandler(l.svcCtx).MerchantGrantCoupon)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
