package merchant

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *ShopHandler) Apply(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	var in types.ApplyReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	app, err := h.logic.Apply(r.Context(), userID, in)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, app)
}

func (h *ShopHandler) MyShops(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	shops, err := h.logic.MyShops(r.Context(), userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, shops)
}

func (h *ShopHandler) UpdateMyShop(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	shopID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "店铺ID无效"))
		return
	}
	var req types.UpdateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateMyShop(r.Context(), shopID, userID, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
