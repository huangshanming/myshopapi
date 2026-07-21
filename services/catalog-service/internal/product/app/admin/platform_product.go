package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"
	"time"
)

func (h *PlatformProductHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	catID, _ := strconv.ParseUint(in.QueryGet("category_id"), 10, 64)
	f := repository.ProductListFilter{
		ShopID: shopID, Name: in.QueryGet("name"), ProductNo: in.QueryGet("product_no"),
		CategoryID: catID, Status: in.QueryGet("status"), ProductType: in.QueryGet("product_type"),
		Page: page, PageSize: pageSize, OrderBy: in.QueryGet("order_by"),
		PlatformScope: true,
	}
	if s := in.QueryGet("created_from"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if s := in.QueryGet("created_to"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	if s := in.QueryGet("publish_from"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.PublishFrom = &t
		}
	}
	if s := in.QueryGet("publish_to"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.PublishTo = &end
		}
	}
	data, err := h.logic.List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *PlatformProductHandler) OffSale(ctx context.Context, in appinput.CallInput) (any, error) {
	uid, _ := middleware.GetUserID(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.PlatformProductRemarkReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.ForceOffSale(ctx, id, uid, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *PlatformProductHandler) Delete(ctx context.Context, in appinput.CallInput) (any, error) {
	uid, _ := middleware.GetUserID(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.PlatformProductRemarkReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.SoftDelete(ctx, id, uid, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
