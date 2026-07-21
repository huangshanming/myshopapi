package product

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hmerchant "mymall/services/catalog-service/internal/product/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantCreateProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCreateProductLogic {
	return &MerchantCreateProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateProductLogic) MerchantCreateProduct(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hmerchant.NewProductHandler(l.svcCtx).Create(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
