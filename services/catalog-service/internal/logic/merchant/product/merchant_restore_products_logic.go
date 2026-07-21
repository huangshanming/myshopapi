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

type MerchantRestoreProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantRestoreProductsLogic(svcCtx *svc.ServiceContext) *MerchantRestoreProductsLogic {
	return &MerchantRestoreProductsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantRestoreProductsLogic) MerchantRestoreProducts(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/products/recycle/restore", nil, nil, req, hmerchant.NewProductHandler(l.svcCtx).RecycleRestore)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
