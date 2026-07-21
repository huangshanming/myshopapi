package user

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"
)

func (h *FavoriteHandler) Add(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	var req struct {
		ProductID uint64 `json:"product_id"`
	}
	if err := appinput.BindBody(in, &req); err != nil || req.ProductID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.Add(ctx, userID, req.ProductID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *FavoriteHandler) Remove(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	productID, err := strconv.ParseUint(in.Path("product_id"), 10, 64)
	if err != nil || productID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	if err := h.logic.Remove(ctx, userID, productID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *FavoriteHandler) RemoveBatch(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	var req struct {
		ProductIDs []uint64 `json:"product_ids"`
	}
	if err := appinput.BindBody(in, &req); err != nil || len(req.ProductIDs) == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.RemoveBatch(ctx, userID, req.ProductIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *FavoriteHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.List(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *FavoriteHandler) Status(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	productID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || productID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	okFav, err := h.logic.IsFavorited(ctx, userID, productID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]bool{"favorited": okFav}, nil
}

func (h *FavoriteHandler) Count(ctx context.Context, in appinput.CallInput) (any, error) {
	productID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || productID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	n, err := h.logic.FavoriteCount(productID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "商品不存在")
	}
	return map[string]int64{"count": n}, nil
}
