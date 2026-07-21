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

type MerchantCopyProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCopyProductLogic(svcCtx *svc.ServiceContext) *MerchantCopyProductLogic {
	return &MerchantCopyProductLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCopyProductLogic) MerchantCopyProduct(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/products/:id/copy", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewProductHandler(l.svcCtx).Copy)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
