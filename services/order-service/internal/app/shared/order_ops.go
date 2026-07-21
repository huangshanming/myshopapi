package shared

import (
	"context"
	"net/http"
	"strconv"

	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/types"
)

func Ship(ctx context.Context, in appinput.CallInput, l *biz.OrderLogic, shopID uint64) (any, error) {
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	var req types.ShipReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := l.Ship(ctx, orderID, shopID, req.ShipCompany, req.ShipNo); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func Complete(ctx context.Context, in appinput.CallInput, l *biz.OrderLogic, shopID uint64) (any, error) {
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	if err := l.Complete(ctx, orderID, shopID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func Remark(ctx context.Context, in appinput.CallInput, l *biz.OrderLogic, shopID uint64) (any, error) {
	orderID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	var req types.RemarkReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := l.UpdateRemark(ctx, orderID, shopID, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func HandleAfterSale(ctx context.Context, in appinput.CallInput, l *biz.OrderLogic, shopID, handledBy uint64) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	var req types.HandleAfterSaleReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := l.HandleAfterSale(ctx, id, shopID, handledBy, req.Action, req.AdminRemark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
