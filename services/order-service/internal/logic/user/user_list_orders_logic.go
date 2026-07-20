package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type UserListOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListOrdersLogic {
	return &UserListOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListOrdersLogic) UserListOrders(w http.ResponseWriter, r *http.Request) {
	huser.NewOrderHandler(l.svcCtx).List(w, r)
}
