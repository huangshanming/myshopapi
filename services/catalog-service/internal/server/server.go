package server

import (
	"context"
	"fmt"

	catalogv1 "mymall/api/gen/catalog/v1"
	"mymall/pkg/grpc/interceptor"
	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/repository"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type CatalogGRPCServer struct {
	catalogv1.UnimplementedCatalogServiceServer
	logic *logic.CatalogLogic
}

func NewCatalogGRPCServer(l *logic.CatalogLogic) *CatalogGRPCServer {
	return &CatalogGRPCServer{logic: l}
}

func (s *CatalogGRPCServer) BatchGetProducts(ctx context.Context, req *catalogv1.BatchGetProductsRequest) (*catalogv1.BatchGetProductsResponse, error) {
	products, err := s.logic.BatchGetProducts(req.GetProductIds())
	if err != nil {
		return nil, err
	}
	resp := &catalogv1.BatchGetProductsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, &catalogv1.Product{
			Id:        p.ID,
			ProductNo: p.ProductNo,
			Name:      p.Name,
			SalePrice: p.SalePrice,
			Stock:     int32(p.Stock),
			Status:    p.Status,
			ShopId:    p.ShopID,
		})
	}
	return resp, nil
}

func (s *CatalogGRPCServer) ReserveStock(ctx context.Context, req *catalogv1.ReserveStockRequest) (*catalogv1.ReserveStockResponse, error) {
	items := make([]repository.StockItem, 0, len(req.GetItems()))
	for _, it := range req.GetItems() {
		items = append(items, repository.StockItem{ProductID: it.GetProductId(), Quantity: int(it.GetQuantity())})
	}
	if err := s.logic.ReserveStock(items); err != nil {
		return &catalogv1.ReserveStockResponse{Success: false, Message: err.Error()}, nil
	}
	return &catalogv1.ReserveStockResponse{Success: true, Message: "ok"}, nil
}

func (s *CatalogGRPCServer) ReleaseStock(ctx context.Context, req *catalogv1.ReleaseStockRequest) (*catalogv1.ReleaseStockResponse, error) {
	items := make([]repository.StockItem, 0, len(req.GetItems()))
	for _, it := range req.GetItems() {
		items = append(items, repository.StockItem{ProductID: it.GetProductId(), Quantity: int(it.GetQuantity())})
	}
	if err := s.logic.ReleaseStock(items); err != nil {
		return &catalogv1.ReleaseStockResponse{Success: false, Message: err.Error()}, nil
	}
	return &catalogv1.ReleaseStockResponse{Success: true, Message: "ok"}, nil
}

func StartZRPC(port int, l *logic.CatalogLogic, logger *zap.Logger) *zrpc.RpcServer {
	c := zrpc.RpcServerConf{
		ListenOn: fmt.Sprintf("0.0.0.0:%d", port),
	}
	c.Mode = service.DevMode
	c.Log.Mode = "console"
	c.Log.Encoding = "plain"

	s := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		catalogv1.RegisterCatalogServiceServer(grpcServer, NewCatalogGRPCServer(l))
	})
	s.AddUnaryInterceptors(interceptor.Logging(logger))
	return s
}
