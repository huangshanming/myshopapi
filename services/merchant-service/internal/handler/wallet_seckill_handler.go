package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/merchant-service/internal/types"
)

func (h *MerchantHandler) AdminGetWallet(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	wallet, err := h.logic.GetWallet(shopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, wallet, "ok")
}

func (h *MerchantHandler) AdminAdjustWallet(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	adminID, _ := middleware.GetUserID(r.Context())
	var req types.WalletAdjustReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	wallet, err := h.logic.AdjustWallet(shopID, req.Field, req.Amount, req.Remark, adminID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, wallet, "调账成功")
}

func (h *MerchantHandler) AdminWalletLogs(w http.ResponseWriter, r *http.Request) {
	shopID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(shopID, p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "ok")
}

func (h *MerchantHandler) MerchantGetWallet(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	wallet, err := h.logic.GetWallet(shopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, wallet, "ok")
}

func (h *MerchantHandler) MerchantWalletLogs(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListWalletLogs(shopID, p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "ok")
}

func (h *MerchantHandler) AdminGetSeckillRule(w http.ResponseWriter, r *http.Request) {
	rule, err := h.logic.GetSeckillRule()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, rule, "ok")
}

func (h *MerchantHandler) AdminUpdateSeckillRule(w http.ResponseWriter, r *http.Request) {
	var req types.SeckillRuleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	rule, err := h.logic.UpdateSeckillRule(req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, rule, "已保存")
}

func (h *MerchantHandler) AdminListSeckillSessions(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListSeckillSessions(p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "ok")
}

func (h *MerchantHandler) AdminListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	sid, _ := strconv.ParseUint(r.URL.Query().Get("session_id"), 10, 64)
	list, total, err := h.logic.ListAdminSeckillEntries(sid, p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "ok")
}

func (h *MerchantHandler) MerchantSeckillSessions(w http.ResponseWriter, r *http.Request) {
	data, err := h.logic.MerchantSeckillSessions()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *MerchantHandler) MerchantApplySeckill(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	var req types.SeckillApplyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	entry, err := h.logic.ApplySeckill(shopID, userID, req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, entry, "报名成功")
}

func (h *MerchantHandler) MerchantListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListShopSeckillEntries(shopID, p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "ok")
}

func (h *MerchantHandler) PublicSeckillCurrent(w http.ResponseWriter, r *http.Request) {
	data, err := h.logic.PublicSeckillCurrent()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *MerchantHandler) PublicSeckillList(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	data, err := h.logic.PublicSeckillList(p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *MerchantHandler) PublicSeckillEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "ID无效", http.StatusBadRequest)
		return
	}
	data, err := h.logic.PublicSeckillEntry(id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response.Success(w, data, "ok")
}

func (h *MerchantHandler) SeckillConsume(w http.ResponseWriter, r *http.Request) {
	var req types.SeckillConsumeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	data, err := h.logic.ConsumeSeckill(req.EntryID, req.ProductID, req.Quantity)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, data, "ok")
}

func (h *MerchantHandler) SeckillRestore(w http.ResponseWriter, r *http.Request) {
	var req types.SeckillRestoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.RestoreSeckill(req.EntryID, req.Quantity); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}
