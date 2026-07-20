package wallet

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/admin/wallet"
	"mymall/services/merchant-service/internal/svc"
)

func AdminAdjustWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewAdminAdjustWalletLogic(r.Context(), svcCtx)
		l.AdminAdjustWallet(w, r)
	}
}

func AdminGetWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewAdminGetWalletLogic(r.Context(), svcCtx)
		l.AdminGetWallet(w, r)
	}
}

func AdminWalletLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewAdminWalletLogsLogic(r.Context(), svcCtx)
		l.AdminWalletLogs(w, r)
	}
}
