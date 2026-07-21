package admin

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

func (h *HomepageThemeHandler) AdminListThemeSlots(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.AdminListThemeSlots()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *HomepageThemeHandler) AdminUpdateThemeSlot(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	var body map[string]interface{}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
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
		return nil, xerr.New(http.StatusBadRequest, "无更新字段")
	}
	if err := h.logic.AdminUpdateThemeSlot(id, updates); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *HomepageThemeHandler) AdminListThemePackages(ctx context.Context, in appinput.CallInput) (any, error) {
	slotID, _ := strconv.ParseUint(in.QueryGet("theme_slot_id"), 10, 64)
	list, err := h.logic.ListThemePackages(slotID, false)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *HomepageThemeHandler) AdminCreateThemePackage(ctx context.Context, in appinput.CallInput) (any, error) {
	var p model.HomepageThemePackage
	if err := appinput.BindBody(in, &p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.AdminCreateThemePackage(&p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *HomepageThemeHandler) AdminUpdateThemePackage(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	var p model.HomepageThemePackage
	if err := appinput.BindBody(in, &p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	updates := map[string]interface{}{
		"theme_slot_id": p.ThemeSlotID, "name": p.Name, "price": p.Price,
		"duration_days": p.DurationDays, "status": p.Status, "sort": p.Sort, "remark": p.Remark,
	}
	if err := h.logic.AdminUpdateThemePackage(id, updates); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *HomepageThemeHandler) AdminListThemeOrders(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	slotID, _ := strconv.ParseUint(in.QueryGet("theme_slot_id"), 10, 64)
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.ListThemeOrders(shopID, slotID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *HomepageThemeHandler) AdminGrantTheme(ctx context.Context, in appinput.CallInput) (any, error) {
	adminID, _ := middleware.GetUserID(ctx)
	var req biz.ThemeGrantReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	o, err := h.logic.GrantTheme(adminID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return o, nil
}
