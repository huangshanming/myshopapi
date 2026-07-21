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

type MerchantCouponClaimsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCouponClaimsLogic(svcCtx *svc.ServiceContext) *MerchantCouponClaimsLogic {
	return &MerchantCouponClaimsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCouponClaimsLogic) MerchantCouponClaims(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/merchant/coupons/:id/claims", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hmerchant.NewCouponHandler(l.svcCtx).MerchantCouponClaims)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
