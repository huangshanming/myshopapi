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

type MerchantHandleAfterSaleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantHandleAfterSaleLogic(svcCtx *svc.ServiceContext) *MerchantHandleAfterSaleLogic {
	return &MerchantHandleAfterSaleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantHandleAfterSaleLogic) MerchantHandleAfterSale(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/merchant/after-sales/:id/handle", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewOrderHandler(l.svcCtx).MerchantHandleAfterSale)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
