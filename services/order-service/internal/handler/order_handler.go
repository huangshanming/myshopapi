package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/repository"
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
	order, err := h.logic.CreateOrder(r.Context(), userID, req.AddressID, req.Items)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, order, "创建成功")
}

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

func (h *OrderHandler) Detail(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}
	order, err := h.logic.GetOrder(userID, orderID)
	if err != nil {
		response.Error(w, "订单不存在", http.StatusNotFound)
		return
	}
	response.Success(w, order, "查询成功")
}

func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}
	if err := h.logic.CancelOrder(r.Context(), userID, orderID); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "取消成功")
}

func (h *OrderHandler) CreateAfterSale(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}
	var req types.CreateAfterSaleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	as, err := h.logic.CreateAfterSale(userID, orderID, req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, as, "申请成功")
}

func (h *OrderHandler) MerchantList(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少 shop_id", http.StatusForbidden)
		return
	}
	p, ps := middleware.ParsePage(r)
	orders, total, err := h.logic.ListByShop(shopID, p, ps, r.URL.Query().Get("status"), r.URL.Query().Get("order_no"))
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
	as, _ := h.logic.ListAfterSalesByOrder(orderID)
	response.Success(w, map[string]interface{}{"order": order, "after_sales": as}, "查询成功")
}

func (h *OrderHandler) MerchantShip(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	h.ship(w, r, shopID)
}

func (h *OrderHandler) MerchantComplete(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	h.complete(w, r, shopID)
}

func (h *OrderHandler) MerchantRemark(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	h.remark(w, r, shopID)
}

func (h *OrderHandler) MerchantAfterSales(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少 shop_id", http.StatusForbidden)
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListAfterSales(repository.AfterSaleListFilter{
		ShopID: shopID, Status: r.URL.Query().Get("status"), OrderNo: r.URL.Query().Get("order_no"),
		Page: p, PageSize: ps,
	})
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "查询成功")
}

func (h *OrderHandler) MerchantHandleAfterSale(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	uid, _ := middleware.GetUserID(r.Context())
	h.handleAfterSale(w, r, shopID, uid)
}

func (h *OrderHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	orders, total, err := h.logic.ListAll(shopID, p, ps, r.URL.Query().Get("status"), r.URL.Query().Get("order_no"))
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
	as, _ := h.logic.ListAfterSalesByOrder(orderID)
	response.Success(w, map[string]interface{}{"order": order, "after_sales": as}, "查询成功")
}

func (h *OrderHandler) AdminShip(w http.ResponseWriter, r *http.Request) {
	h.ship(w, r, 0)
}

func (h *OrderHandler) AdminComplete(w http.ResponseWriter, r *http.Request) {
	h.complete(w, r, 0)
}

func (h *OrderHandler) AdminRemark(w http.ResponseWriter, r *http.Request) {
	h.remark(w, r, 0)
}

func (h *OrderHandler) AdminAfterSales(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	list, total, err := h.logic.ListAfterSales(repository.AfterSaleListFilter{
		ShopID: shopID, Status: r.URL.Query().Get("status"), OrderNo: r.URL.Query().Get("order_no"),
		Page: p, PageSize: ps,
	})
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "查询成功")
}

func (h *OrderHandler) AdminHandleAfterSale(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r.Context())
	h.handleAfterSale(w, r, 0, uid)
}

func (h *OrderHandler) ship(w http.ResponseWriter, r *http.Request, shopID uint64) {
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}
	var req types.ShipReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.Ship(orderID, shopID, req.ShipCompany, req.ShipNo); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "发货成功")
}

func (h *OrderHandler) complete(w http.ResponseWriter, r *http.Request, shopID uint64) {
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}
	if err := h.logic.Complete(orderID, shopID); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已完成")
}

func (h *OrderHandler) remark(w http.ResponseWriter, r *http.Request, shopID uint64) {
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "订单ID无效", http.StatusBadRequest)
		return
	}
	var req types.RemarkReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.UpdateRemark(orderID, shopID, req.Remark); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已更新")
}

func (h *OrderHandler) handleAfterSale(w http.ResponseWriter, r *http.Request, shopID, handledBy uint64) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "ID无效", http.StatusBadRequest)
		return
	}
	var req types.HandleAfterSaleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.HandleAfterSale(r.Context(), id, shopID, handledBy, req.Action, req.AdminRemark); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "处理成功")
}
