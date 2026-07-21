package wallet

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

type MerchantGetWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantGetWalletLogic(svcCtx *svc.ServiceContext) *MerchantGetWalletLogic {
	return &MerchantGetWalletLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantGetWalletLogic) MerchantGetWallet(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/merchant/wallet", nil, nil, nil, hmerchant.NewWalletHandler(l.svcCtx).MerchantGetWallet)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
