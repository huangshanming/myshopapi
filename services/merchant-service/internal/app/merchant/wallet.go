package merchant

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *WalletHandler) MerchantGetWallet(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	wallet, err := h.logic.GetWallet(shopID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return wallet, nil
}

func (h *WalletHandler) MerchantWalletLogs(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	p, ps := in.Page()
	list, total, err := h.logic.ListWalletLogs(shopID, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}
