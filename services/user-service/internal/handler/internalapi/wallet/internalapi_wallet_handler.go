package wallet

import (
	"net/http"

	"mymall/services/user-service/internal/logic/internalapi/wallet"
	"mymall/services/user-service/internal/svc"
)

func InternalFreezeWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewInternalFreezeWalletLogic(r.Context(), svcCtx)
		l.InternalFreezeWallet(w, r)
	}
}

func InternalSettleWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewInternalSettleWalletLogic(r.Context(), svcCtx)
		l.InternalSettleWallet(w, r)
	}
}

func InternalUnfreezeWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewInternalUnfreezeWalletLogic(r.Context(), svcCtx)
		l.InternalUnfreezeWallet(w, r)
	}
}
