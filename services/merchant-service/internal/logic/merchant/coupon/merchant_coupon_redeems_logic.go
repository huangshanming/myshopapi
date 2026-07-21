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

type MerchantCouponRedeemsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCouponRedeemsLogic(svcCtx *svc.ServiceContext) *MerchantCouponRedeemsLogic {
	return &MerchantCouponRedeemsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCouponRedeemsLogic) MerchantCouponRedeems(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/merchant/coupons/:id/redeems", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hmerchant.NewCouponHandler(l.svcCtx).MerchantCouponRedeems)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
