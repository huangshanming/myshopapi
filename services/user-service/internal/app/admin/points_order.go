package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/services/user-service/internal/biz"
	"net/http"
	"strconv"

	"mymall/pkg/xerr"
)

func (h *PointsOrderHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	var userID uint64
	if v := in.QueryGet("user_id"); v != "" {
		userID, _ = strconv.ParseUint(v, 10, 64)
	}
	list, total, err := h.logic.AdminList(ctx, page, pageSize,
		in.QueryGet("status"),
		in.QueryGet("order_no"),
		in.QueryGet("keyword"),
		userID,
	)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *PointsOrderHandler) Detail(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	o, err := h.logic.AdminGet(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return o, nil
}

func (h *PointsOrderHandler) Ship(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	var req biz.ShipReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	o, err := h.logic.AdminShip(ctx, id, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return o, nil
}

func (h *PointsOrderHandler) Complete(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	o, err := h.logic.AdminComplete(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return o, nil
}

func (h *PointsOrderHandler) Cancel(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	var req biz.RemarkReq
	_ = appinput.BindBody(in, &req)
	o, err := h.logic.AdminCancel(ctx, id, req.AdminRemark)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return o, nil
}

func (h *PointsOrderHandler) Remark(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	var req biz.RemarkReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	o, err := h.logic.AdminRemark(ctx, id, req.AdminRemark)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return o, nil
}
