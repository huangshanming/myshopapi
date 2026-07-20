package internalapi

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

type NotificationHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.UserLogic
}

func NewNotificationHandler(svcCtx *svc.ServiceContext) *NotificationHandler {
	return &NotificationHandler{
		svcCtx: svcCtx,
		logic:  biz.NewUserLogic(context.Background(), svcCtx),
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
