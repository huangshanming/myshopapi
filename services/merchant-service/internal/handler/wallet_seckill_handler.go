package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/pkg/xerr")

func (h *MerchantHandler) AdminGetWallet(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	wallet, err := h.logic.GetWallet(shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *MerchantHandler) AdminAdjustWallet(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	adminID, _ := middleware.GetUserID(r.Context())
	var req types.WalletAdjustReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	wallet, err := h.logic.AdjustWallet(shopID, req.Field, req.Amount, req.Remark, adminID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *MerchantHandler) AdminWalletLogs(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(shopID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *MerchantHandler) MerchantGetWallet(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	wallet, err := h.logic.GetWallet(shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *MerchantHandler) MerchantWalletLogs(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(shopID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *MerchantHandler) AdminGetSeckillRule(w http.ResponseWriter, r *http.Request) {
	rule, err := h.logic.GetSeckillRule()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, rule)
}

func (h *MerchantHandler) AdminUpdateSeckillRule(w http.ResponseWriter, r *http.Request) {
	var req types.SeckillRuleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	rule, err := h.logic.UpdateSeckillRule(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, rule)
}

func (h *MerchantHandler) AdminListSeckillSessions(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListSeckillSessions(p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *MerchantHandler) AdminListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	sid, _ := strconv.ParseUint(r.URL.Query().Get("session_id"), 10, 64)
	list, total, err := h.logic.ListAdminSeckillEntries(sid, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *MerchantHandler) MerchantSeckillSessions(w http.ResponseWriter, r *http.Request) {
	data, err := h.logic.MerchantSeckillSessions()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *MerchantHandler) MerchantApplySeckill(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	var req types.SeckillApplyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	entry, err := h.logic.ApplySeckill(shopID, userID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, entry)
}

func (h *MerchantHandler) MerchantListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListShopSeckillEntries(shopID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *MerchantHandler) PublicSeckillCurrent(w http.ResponseWriter, r *http.Request) {
	data, err := h.logic.PublicSeckillCurrent()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *MerchantHandler) PublicSeckillList(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	data, err := h.logic.PublicSeckillList(p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *MerchantHandler) PublicSeckillEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "ID无效"))
		return
	}
	data, err := h.logic.PublicSeckillEntry(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *MerchantHandler) SeckillConsume(w http.ResponseWriter, r *http.Request) {
	var req types.SeckillConsumeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	data, err := h.logic.ConsumeSeckill(req.EntryID, req.ProductID, req.Quantity)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *MerchantHandler) SeckillRestore(w http.ResponseWriter, r *http.Request) {
	var req types.SeckillRestoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.RestoreSeckill(req.EntryID, req.Quantity); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
