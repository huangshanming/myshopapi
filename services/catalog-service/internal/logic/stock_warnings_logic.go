package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type StockWarningsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStockWarningsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StockWarningsLogic {
	return &StockWarningsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StockWarningsLogic) StockWarnings(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).StockWarnings(w, r)
}
