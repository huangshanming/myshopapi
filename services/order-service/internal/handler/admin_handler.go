package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) MerchantList(c *gin.Context) {
	shopID := middleware.GetShopID(c)
	if shopID == 0 {
		response.Error(c, "缺少 shop_id", http.StatusForbidden)
		return
	}
	var page pagination.PageReq
	_ = c.ShouldBindQuery(&page)
	p, ps, _ := pagination.Normalize(&page)
	orders, total, err := h.svc.ListByShop(shopID, p, ps)
	if err != nil {
		response.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(c, gin.H{"total": total, "list": orders}, "查询成功")
}

func (h *OrderHandler) MerchantDetail(c *gin.Context) {
	shopID := middleware.GetShopID(c)
	if shopID == 0 {
		response.Error(c, "缺少 shop_id", http.StatusForbidden)
		return
	}
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "订单ID无效", http.StatusBadRequest)
		return
	}
	order, err := h.svc.GetOrderByShop(shopID, orderID)
	if err != nil {
		response.Error(c, "订单不存在", http.StatusNotFound)
		return
	}
	response.Success(c, order, "查询成功")
}

func (h *OrderHandler) AdminList(c *gin.Context) {
	var page pagination.PageReq
	_ = c.ShouldBindQuery(&page)
	p, ps, _ := pagination.Normalize(&page)
	shopID, _ := strconv.ParseUint(c.Query("shop_id"), 10, 64)
	orders, total, err := h.svc.ListAll(shopID, p, ps)
	if err != nil {
		response.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(c, gin.H{"total": total, "list": orders}, "查询成功")
}

func (h *OrderHandler) AdminDetail(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "订单ID无效", http.StatusBadRequest)
		return
	}
	order, err := h.svc.GetOrderAdmin(orderID)
	if err != nil {
		response.Error(c, "订单不存在", http.StatusNotFound)
		return
	}
	response.Success(c, order, "查询成功")
}
