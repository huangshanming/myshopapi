package admin

import (
	"encoding/json"
	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/types"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *WalletHandler) AdminGetWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "用户ID无效"))
		return
	}
	wallet, err := h.logic.GetWallet(r.Context(), userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *WalletHandler) AdminAdjustWallet(w http.ResponseWriter, r *http.Request) {
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
	wallet, err := h.logic.AdjustWallet(r.Context(), userID, req.Field, req.Amount, req.Remark, adminID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *WalletHandler) AdminWalletLogs(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "用户ID无效"))
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(r.Context(), userID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}
