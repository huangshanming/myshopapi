package merchant

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/types"
)

func (h *HomepageSlotHandler) MerchantListSlotPackages(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListSlotPackages(r.URL.Query().Get("slot_type"), true)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *HomepageSlotHandler) MerchantBuySlot(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	var req logic.BuySlotReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	order, err := h.logic.BuySlot(shopID, userID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, order)
}

func (h *HomepageSlotHandler) MerchantListSlotOrders(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺"))
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListSlotOrders(shopID, r.URL.Query().Get("slot_type"), r.URL.Query().Get("status"), p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}
