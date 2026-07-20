package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type ConfirmReceiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmReceiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmReceiveLogic {
	return &ConfirmReceiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmReceiveLogic) ConfirmReceive(w http.ResponseWriter, r *http.Request) {
	huser.NewOrderHandler(l.svcCtx).ConfirmReceive(w, r)
}
