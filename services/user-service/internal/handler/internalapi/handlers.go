package internalapi

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

type NotificationHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.UserLogic
}

func NewNotificationHandler(svcCtx *svc.ServiceContext) *NotificationHandler {
	return &NotificationHandler{
		svcCtx: svcCtx,
		logic:  logic.NewUserLogic(context.Background(), svcCtx),
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
