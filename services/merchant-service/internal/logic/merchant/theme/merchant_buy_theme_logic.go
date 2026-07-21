package theme

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hmerchant "mymall/services/merchant-service/internal/app/merchant"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantBuyThemeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantBuyThemeLogic(svcCtx *svc.ServiceContext) *MerchantBuyThemeLogic {
	return &MerchantBuyThemeLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantBuyThemeLogic) MerchantBuyTheme(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/theme-orders", nil, nil, req, hmerchant.NewHomepageThemeHandler(l.svcCtx).MerchantBuyTheme)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
