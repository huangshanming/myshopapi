package rpclogic

import (
	"context"

	catalogv1 "mymall/api/gen/catalog/v1"
	productlogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReserveStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReserveStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReserveStockLogic {
	return &ReserveStockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReserveStockLogic) ReserveStock(in *catalogv1.ReserveStockRequest) (*catalogv1.ReserveStockResponse, error) {
	items := make([]repository.StockItem, 0, len(in.GetItems()))
	for _, it := range in.GetItems() {
		items = append(items, repository.StockItem{ProductID: it.GetProductId(), SkuID: it.GetSkuId(), Quantity: int(it.GetQuantity())})
	}
	if err := productlogic.NewCatalogLogic(l.svcCtx).ReserveStock(l.ctx, items); err != nil {
		return &catalogv1.ReserveStockResponse{Success: false, Message: err.Error()}, nil
	}
	return &catalogv1.ReserveStockResponse{Success: true, Message: "ok"}, nil
}
