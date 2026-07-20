package order

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type UserCreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateOrderLogic {
	return &UserCreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCreateOrderLogic) UserCreateOrder(w http.ResponseWriter, r *http.Request) {
	huser.NewOrderHandler(l.svcCtx).List(w, r)
}
