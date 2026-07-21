package order

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hmerchant "mymall/services/order-service/internal/app/merchant"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantHandleAfterSaleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantHandleAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantHandleAfterSaleLogic {
	return &MerchantHandleAfterSaleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantHandleAfterSaleLogic) MerchantHandleAfterSale(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hmerchant.NewOrderHandler(l.svcCtx).MerchantHandleAfterSale(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
