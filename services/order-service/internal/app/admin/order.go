package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/app/shared"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/types"
	"net/http"
	"strconv"
)

func (h *OrderHandler) AdminList(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	orders, total, err := h.logic.ListAll(ctx, shopID, p, ps, in.QueryGet("status"), in.QueryGet("order_no"))
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: orders}, nil
}

func (h *OrderHandler) AdminDetail(ctx context.Context, in appinput.CallInput) (any, error) {
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	order, err := h.logic.GetOrderAdmin(ctx, orderID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "订单不存在")
	}
	as, _ := h.logic.ListAfterSalesByOrder(ctx, orderID)
	return map[string]interface{}{"order": order, "after_sales": as}, nil
}

func (h *OrderHandler) AdminShip(ctx context.Context, in appinput.CallInput) (any, error) {
	return shared.Ship(ctx, in, h.logic, 0)
}

func (h *OrderHandler) AdminComplete(ctx context.Context, in appinput.CallInput) (any, error) {
	return shared.Complete(ctx, in, h.logic, 0)
}

func (h *OrderHandler) AdminRemark(ctx context.Context, in appinput.CallInput) (any, error) {
	return shared.Remark(ctx, in, h.logic, 0)
}

func (h *OrderHandler) AdminAfterSales(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	list, total, err := h.logic.ListAfterSales(ctx, repository.AfterSaleListFilter{
		ShopID: shopID, Status: in.QueryGet("status"), OrderNo: in.QueryGet("order_no"),
		Page: p, PageSize: ps,
	})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *OrderHandler) AdminHandleAfterSale(ctx context.Context, in appinput.CallInput) (any, error) {
	uid, _ := middleware.GetUserID(ctx)
	return shared.HandleAfterSale(ctx, in, h.logic, 0, uid)
}
