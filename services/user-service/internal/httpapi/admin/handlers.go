package admin

import (
	"context"

	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
)

type AddressHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.AddressLogic
}

func NewAddressHandler(svcCtx *svc.ServiceContext) *AddressHandler {
	return &AddressHandler{
		svcCtx: svcCtx,
		logic:  biz.NewAddressLogic(context.Background(), svcCtx),
	}
}

type AdminHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.RBACLogic
}

func NewAdminHandler(svcCtx *svc.ServiceContext) *AdminHandler {
	return &AdminHandler{
		svcCtx: svcCtx,
		logic:  biz.NewRBACLogic(context.Background(), svcCtx),
	}
}

type TaskHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.TaskLogic
}

func NewTaskHandler(svcCtx *svc.ServiceContext) *TaskHandler {
	return &TaskHandler{
		svcCtx: svcCtx,
		logic:  biz.NewTaskLogic(context.Background(), svcCtx),
	}
}

type WalletHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.WalletLogic
}

func NewWalletHandler(svcCtx *svc.ServiceContext) *WalletHandler {
	return &WalletHandler{
		svcCtx: svcCtx,
		logic:  biz.NewWalletLogic(context.Background(), svcCtx),
	}
}

type PointsOrderHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.PointsOrderLogic
}

func NewPointsOrderHandler(svcCtx *svc.ServiceContext) *PointsOrderHandler {
	return &PointsOrderHandler{
		svcCtx: svcCtx,
		logic:  biz.NewPointsOrderLogic(context.Background(), svcCtx),
	}
}

type PointsProductHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.PointsProductLogic
}

func NewPointsProductHandler(svcCtx *svc.ServiceContext) *PointsProductHandler {
	return &PointsProductHandler{
		svcCtx: svcCtx,
		logic:  biz.NewPointsProductLogic(context.Background(), svcCtx),
	}
}
