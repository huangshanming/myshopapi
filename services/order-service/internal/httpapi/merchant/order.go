package merchant

import (
	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/httpapi/shared"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/types"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

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
	shared.Ship(w, r, h.logic, shopID)
}

func (h *OrderHandler) MerchantComplete(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	shared.Complete(w, r, h.logic, shopID)
}

func (h *OrderHandler) MerchantRemark(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	shared.Remark(w, r, h.logic, shopID)
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
	shared.HandleAfterSale(w, r, h.logic, shopID, uid)
}
