package merchant

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MerchantGetWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantGetWalletLogic {
	return &MerchantGetWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantGetWalletLogic) MerchantGetWallet(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewWalletHandler(l.svcCtx).MerchantGetWallet(w, r)
}
