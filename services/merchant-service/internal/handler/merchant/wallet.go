package merchant

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *WalletHandler) MerchantGetWallet(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	wallet, err := h.logic.GetWallet(shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *WalletHandler) MerchantWalletLogs(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(shopID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}
