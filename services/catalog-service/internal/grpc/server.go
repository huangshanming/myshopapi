package server

import (
	"context"
	"fmt"
	"net"

	catalogv1 "mymall/api/gen/catalog/v1"
	"mymall/pkg/grpc/interceptor"
	"mymall/services/catalog-service/internal/repository"
	"mymall/services/catalog-service/internal/service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type CatalogGRPCServer struct {
	catalogv1.UnimplementedCatalogServiceServer
	svc *service.CatalogService
}

func NewCatalogGRPCServer(svc *service.CatalogService) *CatalogGRPCServer {
	return &CatalogGRPCServer{svc: svc}
}

func (s *CatalogGRPCServer) BatchGetProducts(ctx context.Context, req *catalogv1.BatchGetProductsRequest) (*catalogv1.BatchGetProductsResponse, error) {
	products, err := s.svc.BatchGetProducts(req.GetProductIds())
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
	if err := s.svc.ReserveStock(items); err != nil {
		return &catalogv1.ReserveStockResponse{Success: false, Message: err.Error()}, nil
	}
	return &catalogv1.ReserveStockResponse{Success: true, Message: "ok"}, nil
}

func (s *CatalogGRPCServer) ReleaseStock(ctx context.Context, req *catalogv1.ReleaseStockRequest) (*catalogv1.ReleaseStockResponse, error) {
	items := make([]repository.StockItem, 0, len(req.GetItems()))
	for _, it := range req.GetItems() {
		items = append(items, repository.StockItem{ProductID: it.GetProductId(), Quantity: int(it.GetQuantity())})
	}
	if err := s.svc.ReleaseStock(items); err != nil {
		return &catalogv1.ReleaseStockResponse{Success: false, Message: err.Error()}, nil
	}
	return &catalogv1.ReleaseStockResponse{Success: true, Message: "ok"}, nil
}

func Listen(port int, svc *service.CatalogService, logger *zap.Logger) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, err
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Logging(logger)))
	catalogv1.RegisterCatalogServiceServer(s, NewCatalogGRPCServer(svc))
	return s, lis, nil
}
