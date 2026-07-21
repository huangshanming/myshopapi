package user

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/types"
	"net/http"
)

func (h *WalletHandler) UserGetWallet(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	wallet, err := h.logic.GetWallet(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return wallet, nil
}

func (h *WalletHandler) UserWalletLogs(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	p, ps := in.Page()
	list, total, err := h.logic.ListWalletLogs(ctx, userID, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}
