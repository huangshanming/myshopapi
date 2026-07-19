package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/pkg/xerr")

type OrderHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.OrderLogic
}

func NewOrderHandler(svcCtx *svc.ServiceContext) *OrderHandler {
	return &OrderHandler{
		svcCtx: svcCtx,
		logic:  logic.NewOrderLogic(context.Background(), svcCtx),
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	var req types.CreateOrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	order, err := h.logic.CreateOrder(r.Context(), userID, req.AddressID, req.Items, req.UserCouponID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, order)
}

func (h *OrderHandler) CouponPreview(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	var req types.CouponPreviewReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	resp, err := h.logic.CouponPreview(r.Context(), userID, req.Items, req.UserCouponID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, resp)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	p, ps := middleware.ParsePage(r)
	status := r.URL.Query().Get("status")
	orders, total, err := h.logic.ListOrders(userID, p, ps, status)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: orders})
}

func (h *OrderHandler) StatusCounts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	counts, err := h.logic.UserOrderStatusCounts(userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, counts)
}

func (h *OrderHandler) UserAfterSales(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListUserAfterSales(userID, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *OrderHandler) Detail(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	order, err := h.logic.GetOrder(userID, orderID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "订单不存在"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, order)
}

func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	if err := h.logic.CancelOrder(r.Context(), userID, orderID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *OrderHandler) ConfirmReceive(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	if err := h.logic.ConfirmReceive(r.Context(), userID, orderID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *OrderHandler) CreateAfterSale(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	var req types.CreateAfterSaleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	as, err := h.logic.CreateAfterSale(userID, orderID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, as)
}

func (h *OrderHandler) MerchantList(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少 shop_id"))
		return
	}
	p, ps := middleware.ParsePage(r)
	orders, total, err := h.logic.ListByShop(shopID, p, ps, r.URL.Query().Get("status"), r.URL.Query().Get("order_no"))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: orders})
}

func (h *OrderHandler) MerchantDetail(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少 shop_id"))
		return
	}
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	order, err := h.logic.GetOrderByShop(shopID, orderID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "订单不存在"))
		return
	}
	as, _ := h.logic.ListAfterSalesByOrder(orderID)
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"order": order, "after_sales": as})
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
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少 shop_id"))
		return
	}
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListAfterSales(repository.AfterSaleListFilter{
		ShopID: shopID, Status: r.URL.Query().Get("status"), OrderNo: r.URL.Query().Get("order_no"),
		Page: p, PageSize: ps,
	})
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
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
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: orders})
}

func (h *OrderHandler) AdminDetail(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	order, err := h.logic.GetOrderAdmin(orderID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "订单不存在"))
		return
	}
	as, _ := h.logic.ListAfterSalesByOrder(orderID)
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"order": order, "after_sales": as})
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
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *OrderHandler) AdminHandleAfterSale(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r.Context())
	h.handleAfterSale(w, r, 0, uid)
}

func (h *OrderHandler) ship(w http.ResponseWriter, r *http.Request, shopID uint64) {
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	var req types.ShipReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.Ship(r.Context(), orderID, shopID, req.ShipCompany, req.ShipNo); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *OrderHandler) complete(w http.ResponseWriter, r *http.Request, shopID uint64) {
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	if err := h.logic.Complete(r.Context(), orderID, shopID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *OrderHandler) remark(w http.ResponseWriter, r *http.Request, shopID uint64) {
	orderID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	var req types.RemarkReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateRemark(orderID, shopID, req.Remark); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *OrderHandler) handleAfterSale(w http.ResponseWriter, r *http.Request, shopID, handledBy uint64) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "ID无效"))
		return
	}
	var req types.HandleAfterSaleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.HandleAfterSale(r.Context(), id, shopID, handledBy, req.Action, req.AdminRemark); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
