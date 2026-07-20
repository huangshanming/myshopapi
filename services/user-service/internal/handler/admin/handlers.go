package admin

import (
	"context"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

type AddressHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.AddressLogic
}

func NewAddressHandler(svcCtx *svc.ServiceContext) *AddressHandler {
	return &AddressHandler{
		svcCtx: svcCtx,
		logic:  logic.NewAddressLogic(context.Background(), svcCtx),
	}
}

type AdminHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.RBACLogic
}

func NewAdminHandler(svcCtx *svc.ServiceContext) *AdminHandler {
	return &AdminHandler{
		svcCtx: svcCtx,
		logic:  logic.NewRBACLogic(context.Background(), svcCtx),
	}
}

type TaskHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.TaskLogic
}

func NewTaskHandler(svcCtx *svc.ServiceContext) *TaskHandler {
	return &TaskHandler{
		svcCtx: svcCtx,
		logic:  logic.NewTaskLogic(context.Background(), svcCtx),
	}
}

type WalletHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.WalletLogic
}

func NewWalletHandler(svcCtx *svc.ServiceContext) *WalletHandler {
	return &WalletHandler{
		svcCtx: svcCtx,
		logic:  logic.NewWalletLogic(context.Background(), svcCtx),
	}
}
