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

type MerchantRemarkLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantRemarkLogic(svcCtx *svc.ServiceContext) *MerchantRemarkLogic {
	return &MerchantRemarkLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantRemarkLogic) MerchantRemark(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/merchant/orders/:id/remark", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewOrderHandler(l.svcCtx).MerchantRemark)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
