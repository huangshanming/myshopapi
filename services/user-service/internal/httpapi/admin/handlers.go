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
