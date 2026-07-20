package product

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type BatchStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchStockLogic {
	return &BatchStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchStockLogic) BatchStock(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).BatchStock(w, r)
}
