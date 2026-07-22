package wallet

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

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

func (l *MerchantGetWalletLogic) MerchantGetWallet(ctx context.Context) (resp *types.WalletResp, err error) {

	shopID := middleware.GetShopID(ctx)
	wallet, err := biz.NewMerchantLogic(l.svcCtx).GetWallet(shopID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.WalletResp{Data: wallet}, nil
}
