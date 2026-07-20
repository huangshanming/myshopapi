package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"mymall/pkg/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type userWalletDeps struct {
	svcCtx *svc.ServiceContext
	logic  *logic.WalletLogic
}

func newUserWalletDeps(svcCtx *svc.ServiceContext) userWalletDeps {
	return userWalletDeps{
		svcCtx: svcCtx,
		logic:  logic.NewWalletLogic(context.Background(), svcCtx),
	}
}

type WalletUserHandler struct{ userWalletDeps }
type WalletAdminHandler struct{ userWalletDeps }
type WalletInternalHandler struct{ userWalletDeps }

func NewWalletUserHandler(svcCtx *svc.ServiceContext) *WalletUserHandler {
	return &WalletUserHandler{userWalletDeps: newUserWalletDeps(svcCtx)}
}
func NewWalletAdminHandler(svcCtx *svc.ServiceContext) *WalletAdminHandler {
	return &WalletAdminHandler{userWalletDeps: newUserWalletDeps(svcCtx)}
}
func NewWalletInternalHandler(svcCtx *svc.ServiceContext) *WalletInternalHandler {
	return &WalletInternalHandler{userWalletDeps: newUserWalletDeps(svcCtx)}
}

func (h *WalletAdminHandler) AdminGetWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "用户ID无效"))
		return
	}
	wallet, err := h.logic.GetWallet(userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *WalletAdminHandler) AdminAdjustWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "用户ID无效"))
		return
	}
	adminID, _ := middleware.GetUserID(r.Context())
	var req types.WalletAdjustReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	wallet, err := h.logic.AdjustWallet(userID, req.Field, req.Amount, req.Remark, adminID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *WalletAdminHandler) AdminWalletLogs(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "用户ID无效"))
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(userID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *WalletUserHandler) UserGetWallet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	wallet, err := h.logic.GetWallet(userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *WalletUserHandler) UserWalletLogs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(userID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *WalletInternalHandler) Freeze(w http.ResponseWriter, r *http.Request) {
	var req types.WalletOrderOpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.FreezeForOrder(req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *WalletInternalHandler) Unfreeze(w http.ResponseWriter, r *http.Request) {
	var req types.WalletOrderOpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UnfreezeOrder(req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *WalletInternalHandler) Settle(w http.ResponseWriter, r *http.Request) {
	var req types.WalletOrderOpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.SettleOrder(req.UserID, req.Amount, req.OrderID, req.OrderNo); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
