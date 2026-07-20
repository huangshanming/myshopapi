package wallet

import (
	"net/http"

	"mymall/services/user-service/internal/logic/user/wallet"
	"mymall/services/user-service/internal/svc"
)

func UserGetWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewUserGetWalletLogic(r.Context(), svcCtx)
		l.UserGetWallet(w, r)
	}
}

func UserWalletLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewUserWalletLogsLogic(r.Context(), svcCtx)
		l.UserWalletLogs(w, r)
	}
}
