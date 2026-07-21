package user

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/types"
	"net/http"
	"strconv"
)

func (h *OrderHandler) Create(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	var req types.CreateOrderReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	order, err := h.logic.CreateOrder(ctx, userID, req.AddressID, req.Items, req.UserCouponID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return order, nil
}

func (h *OrderHandler) CouponPreview(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	var req types.CouponPreviewReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	resp, err := h.logic.CouponPreview(ctx, userID, req.Items, req.UserCouponID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return resp, nil
}

func (h *OrderHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	p, ps := in.Page()
	status := in.QueryGet("status")
	orders, total, err := h.logic.ListOrders(ctx, userID, p, ps, status)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: orders}, nil
}

func (h *OrderHandler) StatusCounts(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	counts, err := h.logic.UserOrderStatusCounts(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return counts, nil
}

func (h *OrderHandler) UserAfterSales(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	p, ps := in.Page()
	list, total, err := h.logic.ListUserAfterSales(ctx, userID, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *OrderHandler) Detail(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	order, err := h.logic.GetOrder(ctx, userID, orderID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "订单不存在")
	}
	return order, nil
}

func (h *OrderHandler) Cancel(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	if err := h.logic.CancelOrder(ctx, userID, orderID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *OrderHandler) ConfirmReceive(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	if err := h.logic.ConfirmReceive(ctx, userID, orderID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *OrderHandler) CreateAfterSale(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	var req types.CreateAfterSaleReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	as, err := h.logic.CreateAfterSale(ctx, userID, orderID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return as, nil
}
