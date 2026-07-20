package wallet

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/merchant/wallet"
	"mymall/services/merchant-service/internal/svc"
)

func MerchantGetWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewMerchantGetWalletLogic(r.Context(), svcCtx)
		l.MerchantGetWallet(w, r)
	}
}

func MerchantWalletLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewMerchantWalletLogsLogic(r.Context(), svcCtx)
		l.MerchantWalletLogs(w, r)
	}
}
