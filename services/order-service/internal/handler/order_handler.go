package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"
)

type OrderHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.OrderLogic
}

func NewOrderHandler(svcCtx *svc.ServiceContext) *OrderHandler {
	return &OrderHandler{
		svcCtx: svcCtx,
		logic:  logic.NewOrderLogic(svcCtx),
	}
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
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	var req types.CreateOrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	order, err := h.logic.CreateOrder(r.Context(), userID, req.Items)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, order, "创建成功")
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
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	p, ps := middleware.ParsePage(r)
	orders, total, err := h.logic.ListOrders(userID, p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: orders}, "查询成功")
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
func (h *OrderHandler) Detail(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单 ID 无效", http.StatusBadRequest)
		return
	}
	order, err := h.logic.GetOrder(userID, orderID)
	if err != nil {
		response.Error(w, "订单不存在", http.StatusNotFound)
		return
	}
	response.Success(w, order, "查询成功")
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
func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单 ID 无效", http.StatusBadRequest)
		return
	}
	if err := h.logic.CancelOrder(r.Context(), userID, orderID); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "取消成功")
}
