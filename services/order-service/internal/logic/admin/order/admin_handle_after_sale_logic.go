package order

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"
)

type AdminHandleAfterSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminHandleAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminHandleAfterSaleLogic {
	return &AdminHandleAfterSaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminHandleAfterSaleLogic) AdminHandleAfterSale(w http.ResponseWriter, r *http.Request) {
	hadmin.NewOrderHandler(l.svcCtx).AdminHandleAfterSale(w, r)
}
