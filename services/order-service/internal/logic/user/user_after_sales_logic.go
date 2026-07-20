package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type UserAfterSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterSalesLogic {
	return &UserAfterSalesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserAfterSalesLogic) UserAfterSales(w http.ResponseWriter, r *http.Request) {
	huser.NewOrderHandler(l.svcCtx).UserAfterSales(w, r)
}
