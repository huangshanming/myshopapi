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

type MerchantUpdateCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUpdateCouponLogic(svcCtx *svc.ServiceContext) *MerchantUpdateCouponLogic {
	return &MerchantUpdateCouponLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUpdateCouponLogic) MerchantUpdateCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/merchant/coupons/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewCouponHandler(l.svcCtx).MerchantUpdateCoupon)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
