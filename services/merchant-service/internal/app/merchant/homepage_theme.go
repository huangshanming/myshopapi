package merchant

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/model"
)

func (h *HomepageThemeHandler) MerchantListThemeSlots(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.AdminListThemeSlots()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	on := make([]model.HomepageThemeSlot, 0)
	for _, s := range list {
		if s.Status == model.ThemeSlotOn {
			on = append(on, s)
		}
	}
	return on, nil
}

func (h *HomepageThemeHandler) MerchantListThemePackages(ctx context.Context, in appinput.CallInput) (any, error) {
	slotID, _ := strconv.ParseUint(in.QueryGet("theme_slot_id"), 10, 64)
	list, err := h.logic.ListThemePackages(slotID, true)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *HomepageThemeHandler) MerchantBuyTheme(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺")
	}
	var req biz.ThemeBuyReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	o, err := h.logic.BuyTheme(shopID, userID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return o, nil
}

func (h *HomepageThemeHandler) MerchantListThemeOrders(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.ListThemeOrders(shopID, 0, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}
