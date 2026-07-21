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

type MerchantCreateAttrTemplateLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateAttrTemplateLogic(svcCtx *svc.ServiceContext) *MerchantCreateAttrTemplateLogic {
	return &MerchantCreateAttrTemplateLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateAttrTemplateLogic) MerchantCreateAttrTemplate(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/attr-templates", nil, nil, req, hmerchant.NewProductHandler(l.svcCtx).SaveAttrTemplate)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
