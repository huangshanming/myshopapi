package wallet

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

type MerchantGetWalletLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantGetWalletLogic {
	return &MerchantGetWalletLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantGetWalletLogic) MerchantGetWallet(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hmerchant.NewWalletHandler(l.svcCtx).MerchantGetWallet(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
