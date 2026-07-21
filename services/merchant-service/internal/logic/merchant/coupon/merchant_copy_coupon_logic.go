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

type MerchantCopyCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCopyCouponLogic(svcCtx *svc.ServiceContext) *MerchantCopyCouponLogic {
	return &MerchantCopyCouponLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCopyCouponLogic) MerchantCopyCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/coupons/:id/copy", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewCouponHandler(l.svcCtx).MerchantCopyCoupon)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
