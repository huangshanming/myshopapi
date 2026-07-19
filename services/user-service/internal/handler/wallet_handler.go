package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type WalletHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.WalletLogic
}

func NewWalletHandler(svcCtx *svc.ServiceContext) *WalletHandler {
	return &WalletHandler{
		svcCtx: svcCtx,
		logic:  logic.NewWalletLogic(svcCtx),
	}
}

func (h *WalletHandler) AdminGetWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "用户ID无效", http.StatusBadRequest)
		return
	}
	wallet, err := h.logic.GetWallet(userID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, wallet, "ok")
}

func (h *WalletHandler) AdminAdjustWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "用户ID无效", http.StatusBadRequest)
		return
	}
	adminID, _ := middleware.GetUserID(r.Context())
	var req types.WalletAdjustReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	wallet, err := h.logic.AdjustWallet(userID, req.Field, req.Amount, req.Remark, adminID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, wallet, "调账成功")
}

func (h *WalletHandler) AdminWalletLogs(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "用户ID无效", http.StatusBadRequest)
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(userID, p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "ok")
}

func (h *WalletHandler) UserGetWallet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		response.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	wallet, err := h.logic.GetWallet(userID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, wallet, "ok")
}

func (h *WalletHandler) UserWalletLogs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		response.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(userID, p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "ok")
}

func (h *WalletHandler) Freeze(w http.ResponseWriter, r *http.Request) {
	var req types.WalletOrderOpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.FreezeForOrder(req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *WalletHandler) Unfreeze(w http.ResponseWriter, r *http.Request) {
	var req types.WalletOrderOpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.UnfreezeOrder(req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *WalletHandler) Settle(w http.ResponseWriter, r *http.Request) {
	var req types.WalletOrderOpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.SettleOrder(req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}
