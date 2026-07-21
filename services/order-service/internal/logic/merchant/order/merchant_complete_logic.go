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

type MerchantCompleteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCompleteLogic(svcCtx *svc.ServiceContext) *MerchantCompleteLogic {
	return &MerchantCompleteLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCompleteLogic) MerchantComplete(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/merchant/orders/:id/complete", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewOrderHandler(l.svcCtx).MerchantComplete)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
