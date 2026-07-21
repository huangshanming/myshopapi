package merchant

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *SeckillHandler) MerchantSeckillSessions(w http.ResponseWriter, r *http.Request) {
	data, err := h.logic.MerchantSeckillSessions()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *SeckillHandler) MerchantApplySeckill(w http.ResponseWriter, r *http.Request) {
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

func (h *SeckillHandler) MerchantSetSeckillAutoRenew(w http.ResponseWriter, r *http.Request) {
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

func (h *SeckillHandler) MerchantListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListShopSeckillEntries(shopID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}
