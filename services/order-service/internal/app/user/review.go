package user

import (
	"context"
	"io"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/model"
	"net/http"
	"strconv"
)

func (h *ReviewHandler) Eligible(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	data, err := h.logic.ReviewEligible(ctx, userID, orderID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return data, nil
}

func (h *ReviewHandler) Create(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	var req model.CreateReviewReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	rev, err := h.logic.Create(ctx, userID, orderID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return rev, nil
}

func (h *ReviewHandler) GetByOrder(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	rev, err := h.logic.GetByOrder(ctx, userID, orderID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return rev, nil
}

func (h *ReviewHandler) ProductList(ctx context.Context, in appinput.CallInput) (any, error) {
	productID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || productID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.ListByProduct(ctx, productID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *ReviewHandler) Upload(ctx context.Context, in appinput.CallInput) (any, error) {
	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	file, hdr, err := in.Request.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "请上传文件")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取文件失败")
	}
	url, err := h.logic.SaveUpload(userID, hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return map[string]string{"url": url}, nil
}
