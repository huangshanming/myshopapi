package rpclogic

import (
	"context"

	catalogv1 "mymall/api/gen/catalog/v1"
	productlogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReleaseStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReleaseStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseStockLogic {
	return &ReleaseStockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReleaseStockLogic) ReleaseStock(in *catalogv1.ReleaseStockRequest) (*catalogv1.ReleaseStockResponse, error) {
	items := make([]repository.StockItem, 0, len(in.GetItems()))
	for _, it := range in.GetItems() {
		items = append(items, repository.StockItem{ProductID: it.GetProductId(), SkuID: it.GetSkuId(), Quantity: int(it.GetQuantity())})
	}
	if err := productlogic.NewCatalogLogic(l.svcCtx).ReleaseStock(l.ctx, items); err != nil {
		return &catalogv1.ReleaseStockResponse{Success: false, Message: err.Error()}, nil
	}
	return &catalogv1.ReleaseStockResponse{Success: true, Message: "ok"}, nil
}
