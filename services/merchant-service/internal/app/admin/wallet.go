package admin

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *WalletHandler) AdminGetWallet(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	wallet, err := h.logic.GetWallet(shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return wallet, nil
}

func (h *WalletHandler) AdminAdjustWallet(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	adminID, _ := middleware.GetUserID(ctx)
	var req types.WalletAdjustReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	wallet, err := h.logic.AdjustWallet(shopID, req.Field, req.Amount, req.Remark, adminID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return wallet, nil
}

func (h *WalletHandler) AdminWalletLogs(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	p, ps := in.Page()
	list, total, err := h.logic.ListWalletLogs(shopID, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}
