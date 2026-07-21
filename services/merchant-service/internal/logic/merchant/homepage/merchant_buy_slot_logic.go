package homepage

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

type MerchantBuySlotLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantBuySlotLogic(svcCtx *svc.ServiceContext) *MerchantBuySlotLogic {
	return &MerchantBuySlotLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantBuySlotLogic) MerchantBuySlot(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/homepage-orders", nil, nil, req, hmerchant.NewHomepageSlotHandler(l.svcCtx).MerchantBuySlot)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
