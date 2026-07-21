package order

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	huser "mymall/services/order-service/internal/app/user"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponPreviewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCouponPreviewLogic(svcCtx *svc.ServiceContext) *CouponPreviewLogic {
	return &CouponPreviewLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *CouponPreviewLogic) CouponPreview(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/orders/coupon-preview", nil, nil, req, huser.NewOrderHandler(l.svcCtx).CouponPreview)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
