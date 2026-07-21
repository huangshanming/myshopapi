package merchant

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

func (h *OrderHandler) MerchantList(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少 shop_id")
	}
	p, ps := in.Page()
	orders, total, err := h.logic.ListByShop(ctx, shopID, p, ps, in.QueryGet("status"), in.QueryGet("order_no"))
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: orders}, nil
}

func (h *OrderHandler) MerchantDetail(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少 shop_id")
	}
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	order, err := h.logic.GetOrderByShop(ctx, shopID, orderID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "订单不存在")
	}
	as, _ := h.logic.ListAfterSalesByOrder(ctx, orderID)
	return map[string]interface{}{"order": order, "after_sales": as}, nil
}

func (h *OrderHandler) MerchantShip(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	return shared.Ship(ctx, in, h.logic, shopID)
}

func (h *OrderHandler) MerchantComplete(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	return shared.Complete(ctx, in, h.logic, shopID)
}

func (h *OrderHandler) MerchantRemark(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	return shared.Remark(ctx, in, h.logic, shopID)
}

func (h *OrderHandler) MerchantAfterSales(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少 shop_id")
	}
	p, ps := in.Page()
	list, total, err := h.logic.ListAfterSales(ctx, repository.AfterSaleListFilter{
		ShopID: shopID, Status: in.QueryGet("status"), OrderNo: in.QueryGet("order_no"),
		Page: p, PageSize: ps,
	})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *OrderHandler) MerchantHandleAfterSale(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	uid, _ := middleware.GetUserID(ctx)
	return shared.HandleAfterSale(ctx, in, h.logic, shopID, uid)
}
