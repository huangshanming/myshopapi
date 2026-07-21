package internalapi

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/types"
	"net/http"
)

func (h *WalletHandler) Freeze(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.WalletOrderOpReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.FreezeForOrder(ctx, req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *WalletHandler) Unfreeze(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.WalletOrderOpReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UnfreezeOrder(ctx, req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *WalletHandler) Settle(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.WalletOrderOpReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.SettleOrder(ctx, req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
