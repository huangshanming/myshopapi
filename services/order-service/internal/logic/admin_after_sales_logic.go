package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminAfterSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAfterSalesLogic {
	return &AdminAfterSalesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminAfterSalesLogic) AdminAfterSales(w http.ResponseWriter, r *http.Request) {
	hadmin.NewOrderHandler(l.svcCtx).AdminAfterSales(w, r)
}
