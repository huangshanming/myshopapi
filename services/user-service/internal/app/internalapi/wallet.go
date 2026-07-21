package internalapi

import (
	"encoding/json"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *WalletHandler) Freeze(w http.ResponseWriter, r *http.Request) {
	var req types.WalletOrderOpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.FreezeForOrder(r.Context(), req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *WalletHandler) Unfreeze(w http.ResponseWriter, r *http.Request) {
	var req types.WalletOrderOpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UnfreezeOrder(r.Context(), req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *WalletHandler) Settle(w http.ResponseWriter, r *http.Request) {
	var req types.WalletOrderOpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.SettleOrder(r.Context(), req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
