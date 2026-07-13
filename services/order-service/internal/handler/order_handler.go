package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/apidoc/dto"
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

// Create 创建订单
// @Summary      创建订单
// @Description  需要 JWT；经 APISIX 网关时从 Authorization 注入 X-User-Id
// @Tags         订单
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-User-Id  header  int             false  "用户 ID（网关注入，直连调试时必填）"
// @Param        body       body    dto.CreateOrderReq  true  "下单商品"
// @Success      201  {object}  apidoc.Response{data=dto.OrderInfo}  "创建成功"
// @Router       /api/v1/orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	var req dto.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	items := make([]service.CreateItemInput, len(req.Items))
	for i, it := range req.Items {
		items[i] = service.CreateItemInput{ProductID: it.ProductID, Quantity: it.Quantity}
	}
	order, err := h.svc.CreateOrder(c.Request.Context(), userID, items)
	if err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": 200, "msg": "创建成功", "data": order})
}

// List 我的订单列表
// @Summary      我的订单列表
// @Description  分页查询当前用户订单
// @Tags         订单
// @Produce      json
// @Security     BearerAuth
// @Param        X-User-Id  header  int  false  "用户 ID（网关注入）"
// @Param        page       query   int  false  "页码"      default(1)
// @Param        page_size  query   int  false  "每页数量"  default(10)
// @Success      200  {object}  apidoc.Response{data=dto.OrderListResp}  "查询成功"
// @Router       /api/v1/orders [get]
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

// Detail 订单详情
// @Summary      订单详情
// @Tags         订单
// @Produce      json
// @Security     BearerAuth
// @Param        X-User-Id  header  int  false  "用户 ID（网关注入）"
// @Param        id         path    int  true   "订单 ID"
// @Success      200  {object}  apidoc.Response{data=dto.OrderInfo}  "查询成功"
// @Router       /api/v1/orders/{id} [get]
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

// Cancel 取消订单
// @Summary      取消订单
// @Description  pending/confirmed 状态可取消，会触发库存回滚
// @Tags         订单
// @Produce      json
// @Security     BearerAuth
// @Param        X-User-Id  header  int  false  "用户 ID（网关注入）"
// @Param        id         path    int  true   "订单 ID"
// @Success      200  {object}  apidoc.Response{data=apidoc.EmptyData}  "取消成功"
// @Router       /api/v1/orders/{id}/cancel [put]
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
