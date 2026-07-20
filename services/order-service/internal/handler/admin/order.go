package admin

import (
	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/handler/shared"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/types"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

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
	shared.Ship(w, r, h.logic, 0)
}

func (h *OrderHandler) AdminComplete(w http.ResponseWriter, r *http.Request) {
	shared.Complete(w, r, h.logic, 0)
}

func (h *OrderHandler) AdminRemark(w http.ResponseWriter, r *http.Request) {
	shared.Remark(w, r, h.logic, 0)
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
	shared.HandleAfterSale(w, r, h.logic, 0, uid)
}
