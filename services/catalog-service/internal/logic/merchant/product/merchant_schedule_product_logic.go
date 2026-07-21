package product

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hmerchant "mymall/services/catalog-service/internal/product/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantScheduleProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantScheduleProductLogic(svcCtx *svc.ServiceContext) *MerchantScheduleProductLogic {
	return &MerchantScheduleProductLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantScheduleProductLogic) MerchantScheduleProduct(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/products/:id/schedules", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewProductHandler(l.svcCtx).Schedule)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
