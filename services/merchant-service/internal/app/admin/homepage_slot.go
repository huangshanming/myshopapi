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
	"mymall/services/merchant-service/internal/types"
)

func (h *HomepageSlotHandler) AdminListSlotPackages(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.ListSlotPackages(in.QueryGet("slot_type"), false)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *HomepageSlotHandler) AdminCreateSlotPackage(ctx context.Context, in appinput.CallInput) (any, error) {
	var p model.HomepageSlotPackage
	if err := appinput.BindBody(in, &p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.CreateSlotPackage(&p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *HomepageSlotHandler) AdminUpdateSlotPackage(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	var p model.HomepageSlotPackage
	if err := appinput.BindBody(in, &p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateSlotPackage(id, &p); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *HomepageSlotHandler) AdminListSlotSettings(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.ListSlotSettings()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *HomepageSlotHandler) AdminUpdateSlotSettings(ctx context.Context, in appinput.CallInput) (any, error) {
	var req struct {
		Items []model.HomepageSlotSetting `json:"items"`
	}
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateSlotSettings(req.Items); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *HomepageSlotHandler) AdminListSlotOrders(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	list, total, err := h.logic.ListSlotOrders(shopID, in.QueryGet("slot_type"), in.QueryGet("status"), p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *HomepageSlotHandler) AdminGrantSlot(ctx context.Context, in appinput.CallInput) (any, error) {
	adminID, _ := middleware.GetUserID(ctx)
	var req biz.GrantSlotReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	order, err := h.logic.GrantSlot(adminID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return order, nil
}
