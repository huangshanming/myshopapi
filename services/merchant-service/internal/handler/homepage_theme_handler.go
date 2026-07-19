package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/model"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *MerchantHandler) PublicThemeTiles(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListThemeTiles()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list})
}

func (h *MerchantHandler) AdminListThemeSlots(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.AdminListThemeSlots()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *MerchantHandler) AdminUpdateThemeSlot(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "ID无效"))
		return
	}
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	allowed := map[string]bool{
		"name": true, "desc": true, "cover_url": true, "default_link_type": true,
		"default_link_id": true, "status": true, "sort": true, "position": true,
	}
	updates := map[string]interface{}{}
	for k, v := range body {
		if allowed[k] {
			updates[k] = v
		}
	}
	if len(updates) == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "无更新字段"))
		return
	}
	if err := h.logic.AdminUpdateThemeSlot(id, updates); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *MerchantHandler) AdminListThemePackages(w http.ResponseWriter, r *http.Request) {
	slotID, _ := strconv.ParseUint(r.URL.Query().Get("theme_slot_id"), 10, 64)
	list, err := h.logic.ListThemePackages(slotID, false)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *MerchantHandler) AdminCreateThemePackage(w http.ResponseWriter, r *http.Request) {
	var p model.HomepageThemePackage
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.AdminCreateThemePackage(&p); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *MerchantHandler) AdminUpdateThemePackage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "ID无效"))
		return
	}
	var p model.HomepageThemePackage
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	updates := map[string]interface{}{
		"theme_slot_id": p.ThemeSlotID, "name": p.Name, "price": p.Price,
		"duration_days": p.DurationDays, "status": p.Status, "sort": p.Sort, "remark": p.Remark,
	}
	if err := h.logic.AdminUpdateThemePackage(id, updates); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *MerchantHandler) AdminListThemeOrders(w http.ResponseWriter, r *http.Request) {
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	slotID, _ := strconv.ParseUint(r.URL.Query().Get("theme_slot_id"), 10, 64)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.ListThemeOrders(shopID, slotID, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *MerchantHandler) AdminGrantTheme(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())
	var req logic.ThemeGrantReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	o, err := h.logic.GrantTheme(adminID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, o)
}

func (h *MerchantHandler) MerchantListThemeSlots(w http.ResponseWriter, r *http.Request) {
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

func (h *MerchantHandler) MerchantListThemePackages(w http.ResponseWriter, r *http.Request) {
	slotID, _ := strconv.ParseUint(r.URL.Query().Get("theme_slot_id"), 10, 64)
	list, err := h.logic.ListThemePackages(slotID, true)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *MerchantHandler) MerchantBuyTheme(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺"))
		return
	}
	var req logic.ThemeBuyReq
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

func (h *MerchantHandler) MerchantListThemeOrders(w http.ResponseWriter, r *http.Request) {
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
