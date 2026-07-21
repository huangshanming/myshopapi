package order

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hmerchant "mymall/services/order-service/internal/app/merchant"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantAfterSalesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantAfterSalesLogic(svcCtx *svc.ServiceContext) *MerchantAfterSalesLogic {
	return &MerchantAfterSalesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantAfterSalesLogic) MerchantAfterSales(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/merchant/after-sales", nil, nil, nil, hmerchant.NewOrderHandler(l.svcCtx).MerchantAfterSales)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
