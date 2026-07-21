package admin

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *ShopHandler) AdminListApplications(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	list, total, err := h.logic.ListApplications(ctx, in.QueryGet("status"), p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *ShopHandler) AdminApprove(ctx context.Context, in appinput.CallInput) (any, error) {
	adminID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	appID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "申请ID无效")
	}
	shop, err := h.logic.Approve(ctx, appID, adminID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return shop, nil
}

func (h *ShopHandler) AdminReject(ctx context.Context, in appinput.CallInput) (any, error) {
	adminID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	appID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "申请ID无效")
	}
	var req types.RejectReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.Reject(ctx, appID, adminID, req.Reason); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ShopHandler) AdminListShops(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	list, total, err := h.logic.ListShops(ctx, in.QueryGet("status"), in.QueryGet("name"), p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *ShopHandler) AdminCreateShop(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.AdminCreateShopReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	shop, err := h.logic.CreateShop(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return shop, nil
}

func (h *ShopHandler) AdminUpdateShop(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	var req types.AdminUpdateShopReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.AdminUpdateShop(ctx, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ShopHandler) AdminResetOwnerPassword(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	var req types.OwnerPasswordReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.ResetOwnerPassword(ctx, id, req.Password); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ShopHandler) AdminGetShop(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	shop, err := h.logic.GetShop(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "店铺不存在")
	}
	return shop, nil
}

func (h *ShopHandler) AdminDisableShop(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	var req types.RejectReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.DisableShop(ctx, id, req.Reason); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ShopHandler) AdminEnableShop(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	if err := h.logic.EnableShop(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
