package theme

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hmerchant "mymall/services/merchant-service/internal/app/merchant"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantBuyThemeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantBuyThemeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantBuyThemeLogic {
	return &MerchantBuyThemeLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantBuyThemeLogic) MerchantBuyTheme(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hmerchant.NewHomepageThemeHandler(l.svcCtx).MerchantBuyTheme(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
