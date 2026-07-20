package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type UserGetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetOrderLogic {
	return &UserGetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetOrderLogic) UserGetOrder(w http.ResponseWriter, r *http.Request) {
	huser.NewOrderHandler(l.svcCtx).Detail(w, r)
}
