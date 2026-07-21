package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/types"
	"net/http"
	"strconv"
)

func (h *WalletHandler) AdminGetWallet(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "用户ID无效")
	}
	wallet, err := h.logic.GetWallet(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return wallet, nil
}

func (h *WalletHandler) AdminAdjustWallet(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "用户ID无效")
	}
	adminID, _ := middleware.GetUserID(ctx)
	var req types.WalletAdjustReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	wallet, err := h.logic.AdjustWallet(ctx, userID, req.Field, req.Amount, req.Remark, adminID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return wallet, nil
}

func (h *WalletHandler) AdminWalletLogs(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "用户ID无效")
	}
	p, ps := in.Page()
	list, total, err := h.logic.ListWalletLogs(ctx, userID, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}
