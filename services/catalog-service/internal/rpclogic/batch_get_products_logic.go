package rpclogic

import (
	"context"

	catalogv1 "mymall/api/gen/catalog/v1"
	productlogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetProductsLogic {
	return &BatchGetProductsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *BatchGetProductsLogic) BatchGetProducts(in *catalogv1.BatchGetProductsRequest) (*catalogv1.BatchGetProductsResponse, error) {
	cl := productlogic.NewCatalogLogic(l.svcCtx)
	products, err := cl.BatchGetProducts(l.ctx, in.GetProductIds())
	if err != nil {
		return nil, err
	}
	resp := &catalogv1.BatchGetProductsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, &catalogv1.Product{
			Id: p.ID, ProductNo: p.ProductNo, Name: p.Name, SalePrice: p.SalePrice,
			Stock: int32(p.Stock), Status: p.Status, ShopId: p.ShopID,
			DefaultSkuId: cl.DefaultSkuID(l.ctx, p.ID),
		})
	}
	return resp, nil
}
