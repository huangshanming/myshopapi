package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/order-service/internal/types"
)

func (h *OrderHandler) MerchantList(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少 shop_id", http.StatusForbidden)
		return
	}
	p, ps := middleware.ParsePage(r)
	orders, total, err := h.logic.ListByShop(shopID, p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: orders}, "查询成功")
}

func (h *OrderHandler) MerchantDetail(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少 shop_id", http.StatusForbidden)
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}
	order, err := h.logic.GetOrderByShop(shopID, orderID)
	if err != nil {
		response.Error(w, "订单不存在", http.StatusNotFound)
		return
	}
	response.Success(w, order, "查询成功")
}

func (h *OrderHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	orders, total, err := h.logic.ListAll(shopID, p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: orders}, "查询成功")
}

func (h *OrderHandler) AdminDetail(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}
	order, err := h.logic.GetOrderAdmin(orderID)
	if err != nil {
		response.Error(w, "订单不存在", http.StatusNotFound)
		return
	}
	response.Success(w, order, "查询成功")
}
