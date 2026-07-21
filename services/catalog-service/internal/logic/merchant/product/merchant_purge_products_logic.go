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

type MerchantPurgeProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantPurgeProductsLogic(svcCtx *svc.ServiceContext) *MerchantPurgeProductsLogic {
	return &MerchantPurgeProductsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantPurgeProductsLogic) MerchantPurgeProducts(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "DELETE", "/api/v1/merchant/products/recycle", nil, nil, nil, hmerchant.NewProductHandler(l.svcCtx).RecycleDelete)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
