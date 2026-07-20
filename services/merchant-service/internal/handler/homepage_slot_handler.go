package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/pkg/xerr"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/types"
)

func (h *HomepageSlotAdminHandler) AdminListSlotPackages(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListSlotPackages(r.URL.Query().Get("slot_type"), false)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *HomepageSlotAdminHandler) AdminCreateSlotPackage(w http.ResponseWriter, r *http.Request) {
	var p model.HomepageSlotPackage
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.CreateSlotPackage(&p); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *HomepageSlotAdminHandler) AdminUpdateSlotPackage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "ID无效"))
		return
	}
	var p model.HomepageSlotPackage
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateSlotPackage(id, &p); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *HomepageSlotAdminHandler) AdminListSlotSettings(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListSlotSettings()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *HomepageSlotAdminHandler) AdminUpdateSlotSettings(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []model.HomepageSlotSetting `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateSlotSettings(req.Items); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *HomepageSlotAdminHandler) AdminListSlotOrders(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	list, total, err := h.logic.ListSlotOrders(shopID, r.URL.Query().Get("slot_type"), r.URL.Query().Get("status"), p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *HomepageSlotAdminHandler) AdminGrantSlot(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())
	var req logic.GrantSlotReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	order, err := h.logic.GrantSlot(adminID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, order)
}

func (h *HomepageSlotMerchantHandler) MerchantListSlotPackages(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.ListSlotPackages(r.URL.Query().Get("slot_type"), true)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *HomepageSlotMerchantHandler) MerchantBuySlot(w http.ResponseWriter, r *http.Request) {
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

func (h *HomepageSlotMerchantHandler) MerchantListSlotOrders(w http.ResponseWriter, r *http.Request) {
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

func (h *HomepageSlotPublicHandler) PublicHomeSlots(w http.ResponseWriter, r *http.Request) {
	slotType := r.URL.Query().Get("slot_type")
	list, err := h.logic.HomeSlots(slotType)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}
