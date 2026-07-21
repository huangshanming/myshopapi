package product

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	hmerchant "mymall/services/catalog-service/internal/product/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantImportProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantImportProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantImportProductsLogic {
	return &MerchantImportProductsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantImportProductsLogic) MerchantImportProducts(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	data, err := hmerchant.NewProductHandler(l.svcCtx).Import(ctx, appinput.CallInput{Request: r})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
