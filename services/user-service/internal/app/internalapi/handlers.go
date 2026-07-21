package internalapi

import (

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
		logic:  biz.NewAddressLogic(svcCtx),
	}
}

type NotificationHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.UserLogic
}

func NewNotificationHandler(svcCtx *svc.ServiceContext) *NotificationHandler {
	return &NotificationHandler{
		svcCtx: svcCtx,
		logic:  biz.NewUserLogic(svcCtx),
	}
}

type TaskHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.TaskLogic
}

func NewTaskHandler(svcCtx *svc.ServiceContext) *TaskHandler {
	return &TaskHandler{
		svcCtx: svcCtx,
		logic:  biz.NewTaskLogic(svcCtx),
	}
}

type WalletHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.WalletLogic
}

func NewWalletHandler(svcCtx *svc.ServiceContext) *WalletHandler {
	return &WalletHandler{
		svcCtx: svcCtx,
		logic:  biz.NewWalletLogic(svcCtx),
	}
}
