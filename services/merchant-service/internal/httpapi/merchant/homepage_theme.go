package merchant

import (
	"encoding/json"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/model"
)

func (h *HomepageThemeHandler) MerchantListThemeSlots(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.AdminListThemeSlots()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	on := make([]model.HomepageThemeSlot, 0)
	for _, s := range list {
		if s.Status == model.ThemeSlotOn {
			on = append(on, s)
		}
	}
	httpx.OkJsonCtx(r.Context(), w, on)
}

func (h *HomepageThemeHandler) MerchantListThemePackages(w http.ResponseWriter, r *http.Request) {
	slotID, _ := strconv.ParseUint(r.URL.Query().Get("theme_slot_id"), 10, 64)
	list, err := h.logic.ListThemePackages(slotID, true)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *HomepageThemeHandler) MerchantBuyTheme(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺"))
		return
	}
	var req biz.ThemeBuyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	o, err := h.logic.BuyTheme(shopID, userID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, o)
}

func (h *HomepageThemeHandler) MerchantListThemeOrders(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺"))
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.ListThemeOrders(shopID, 0, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}
