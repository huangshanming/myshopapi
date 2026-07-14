package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mymall/services/order-service/internal/grpc/client"
	"mymall/services/order-service/internal/model"
	ordermq "mymall/services/order-service/internal/mq"
	"mymall/services/order-service/internal/repository"

	"github.com/google/uuid"
)

type CreateItemInput struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type OrderService struct {
	repo      *repository.OrderRepository
	catalog   *client.CatalogClient
	user      *client.UserClient
	publisher *ordermq.Publisher
}

func NewOrderService(repo *repository.OrderRepository, catalog *client.CatalogClient, user *client.UserClient, publisher *ordermq.Publisher) *OrderService {
	return &OrderService{repo: repo, catalog: catalog, user: user, publisher: publisher}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, items []CreateItemInput) (*model.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("订单商品不能为空")
	}

	if s.user != nil {
		u, err := s.user.GetUser(ctx, userID)
		if err != nil || u.GetStatus() != 1 {
			return nil, errors.New("用户无效")
		}
	}

	ids := make([]uint64, 0, len(items))
	qtyMap := make(map[uint64]int, len(items))
	for _, it := range items {
		ids = append(ids, it.ProductID)
		qtyMap[it.ProductID] = it.Quantity
	}

	resp, err := s.catalog.BatchGetProducts(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("查询商品失败: %w", err)
	}
	if len(resp.GetProducts()) != len(ids) {
		return nil, errors.New("部分商品不存在或已下架")
	}

	var total float64
	var shopID uint64
	orderItems := make([]model.OrderItem, 0, len(items))
	stockItems := make([]model.StockItem, 0, len(items))
	productMap := make(map[uint64]struct{}, len(resp.GetProducts()))
	for _, p := range resp.GetProducts() {
		productMap[p.GetId()] = struct{}{}
		if shopID == 0 {
			shopID = p.GetShopId()
		} else if p.GetShopId() != shopID {
			return nil, errors.New("一期不支持跨店下单，请分开结算")
		}
		qty := qtyMap[p.GetId()]
		if p.GetStock() < int32(qty) {
			return nil, fmt.Errorf("商品 %s 库存不足", p.GetName())
		}
		lineAmount := p.GetSalePrice() * float64(qty)
		total += lineAmount
		orderItems = append(orderItems, model.OrderItem{
			ProductID:   p.GetId(),
			ProductName: p.GetName(),
			Price:       p.GetSalePrice(),
			Quantity:    qty,
		})
		stockItems = append(stockItems, model.StockItem{ProductID: p.GetId(), Quantity: qty})
	}
	for _, id := range ids {
		if _, ok := productMap[id]; !ok {
			return nil, errors.New("部分商品不存在或已下架")
		}
	}

	orderNo := fmt.Sprintf("ORD%s", uuid.NewString()[:8]+time.Now().Format("150405"))
	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      userID,
		ShopID:      shopID,
		TotalAmount: total,
		Status:      model.OrderStatusPending,
	}
	if err := s.repo.Create(order, orderItems); err != nil {
		return nil, err
	}

	if s.publisher != nil {
		if err := s.publisher.PublishOrderCreated(ctx, orderNo, stockItems); err != nil {
			_ = s.repo.UpdateStatus(orderNo, model.OrderStatusFailed)
			return nil, fmt.Errorf("发布订单事件失败: %w", err)
		}
	}

	order.Items = orderItems
	return order, nil
}

func (s *OrderService) ListOrders(userID uint64, page, pageSize int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListByUser(userID, page, pageSize)
}

func (s *OrderService) GetOrder(userID, orderID uint64) (*model.Order, error) {
	return s.repo.FindByID(orderID, userID)
}

func (s *OrderService) CancelOrder(ctx context.Context, userID, orderID uint64) error {
	order, err := s.repo.FindByID(orderID, userID)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusConfirmed {
		return errors.New("当前状态不可取消")
	}
	if err := s.repo.Cancel(orderID, userID); err != nil {
		return err
	}
	if s.publisher != nil {
		items := make([]model.StockItem, 0, len(order.Items))
		for _, it := range order.Items {
			items = append(items, model.StockItem{ProductID: it.ProductID, Quantity: it.Quantity})
		}
		_ = s.publisher.PublishOrderCancelled(ctx, order.OrderNo, items)
	}
	return nil
}

func (s *OrderService) ListByShop(shopID uint64, page, pageSize int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListByShop(shopID, page, pageSize)
}

func (s *OrderService) ListAll(shopID uint64, page, pageSize int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListAll(shopID, page, pageSize)
}

func (s *OrderService) GetOrderByShop(shopID, orderID uint64) (*model.Order, error) {
	return s.repo.FindByIDAndShop(orderID, shopID)
}

func (s *OrderService) GetOrderAdmin(orderID uint64) (*model.Order, error) {
	return s.repo.FindByIDAdmin(orderID)
}
