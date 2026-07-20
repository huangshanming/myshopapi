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

func (h *WalletAdminHandler) AdminGetWallet(w http.ResponseWriter, r *http.Request) {
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

func (h *WalletAdminHandler) AdminAdjustWallet(w http.ResponseWriter, r *http.Request) {
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

func (h *WalletAdminHandler) AdminWalletLogs(w http.ResponseWriter, r *http.Request) {
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

func (h *WalletMerchantHandler) MerchantGetWallet(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	wallet, err := h.logic.GetWallet(shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, wallet)
}

func (h *WalletMerchantHandler) MerchantWalletLogs(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(shopID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *SeckillAdminHandler) AdminGetSeckillRule(w http.ResponseWriter, r *http.Request) {
	rule, err := h.logic.GetSeckillRule()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, rule)
}

func (h *SeckillAdminHandler) AdminUpdateSeckillRule(w http.ResponseWriter, r *http.Request) {
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

func (h *SeckillAdminHandler) AdminListSeckillSessions(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListSeckillSessions(p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *SeckillAdminHandler) AdminListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	sid, _ := strconv.ParseUint(r.URL.Query().Get("session_id"), 10, 64)
	list, total, err := h.logic.ListAdminSeckillEntries(sid, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *SeckillMerchantHandler) MerchantSeckillSessions(w http.ResponseWriter, r *http.Request) {
	data, err := h.logic.MerchantSeckillSessions()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *SeckillMerchantHandler) MerchantApplySeckill(w http.ResponseWriter, r *http.Request) {
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

func (h *SeckillMerchantHandler) MerchantSetSeckillAutoRenew(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "报名ID无效"))
		return
	}
	var req types.SeckillAutoRenewReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	entry, err := h.logic.SetSeckillAutoRenew(shopID, id, req.AutoRenew)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, entry)
}

func (h *SeckillMerchantHandler) MerchantListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListShopSeckillEntries(shopID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *SeckillPublicHandler) PublicSeckillCurrent(w http.ResponseWriter, r *http.Request) {
	data, err := h.logic.PublicSeckillCurrent()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *SeckillPublicHandler) PublicSeckillList(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	data, err := h.logic.PublicSeckillList(p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *SeckillPublicHandler) PublicSeckillEntry(w http.ResponseWriter, r *http.Request) {
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

func (h *SeckillInternalHandler) SeckillConsume(w http.ResponseWriter, r *http.Request) {
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

func (h *SeckillInternalHandler) SeckillRestore(w http.ResponseWriter, r *http.Request) {
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
