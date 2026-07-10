package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/response"
	"mymall/services/order-service/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

type createOrderReq struct {
	Items []service.CreateItemInput `json:"items" binding:"required,min=1,dive"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	var req createOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	order, err := h.svc.CreateOrder(c.Request.Context(), userID, req.Items)
	if err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": 200, "msg": "创建成功", "data": order})
}

func (h *OrderHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	var page pagination.PageReq
	_ = c.ShouldBindQuery(&page)
	p, ps, _ := pagination.Normalize(&page)
	orders, total, err := h.svc.ListOrders(userID, p, ps)
	if err != nil {
		response.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(c, gin.H{"total": total, "list": orders}, "查询成功")
}

func (h *OrderHandler) Detail(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "订单 ID 无效", http.StatusBadRequest)
		return
	}
	order, err := h.svc.GetOrder(userID, orderID)
	if err != nil {
		response.Error(c, "订单不存在", http.StatusNotFound)
		return
	}
	response.Success(c, order, "查询成功")
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "订单 ID 无效", http.StatusBadRequest)
		return
	}
	if err := h.svc.CancelOrder(c.Request.Context(), userID, orderID); err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, nil, "取消成功")
}
