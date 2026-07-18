package mq

import (
	"context"
	"encoding/json"

	"mymall/pkg/cache"
	"mymall/pkg/mq"
	"mymall/services/order-service/internal/model"
	"mymall/services/order-service/internal/repository"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Publisher struct {
	client *mq.Client
}

func NewPublisher(client *mq.Client) *Publisher {
	return &Publisher{client: client}
}

type OrderCreatedEvent struct {
	OrderNo string            `json:"order_no"`
	Items   []model.StockItem `json:"items"`
}

type OrderCancelledEvent struct {
	OrderNo string            `json:"order_no"`
	Items   []model.StockItem `json:"items"`
}

func (p *Publisher) PublishOrderCreated(ctx context.Context, orderNo string, items []model.StockItem) error {
	return p.client.Publish(ctx, mq.RoutingOrderCreated, OrderCreatedEvent{OrderNo: orderNo, Items: items})
}

func (p *Publisher) PublishOrderCancelled(ctx context.Context, orderNo string, items []model.StockItem) error {
	return p.client.Publish(ctx, mq.RoutingOrderCancelled, OrderCancelledEvent{OrderNo: orderNo, Items: items})
}

type InventoryResultEvent struct {
	OrderNo string `json:"order_no"`
	Message string `json:"message,omitempty"`
}

type Consumer struct {
	client *mq.Client
	repo   *repository.OrderRepository
	rdb    *redis.Client
	logger *zap.Logger
}

func NewConsumer(client *mq.Client, repo *repository.OrderRepository, rdb *redis.Client, logger *zap.Logger) *Consumer {
	return &Consumer{client: client, repo: repo, rdb: rdb, logger: logger}
}

func (c *Consumer) Start() error {
	if err := c.client.Consume("order.inventory.reserved", mq.RoutingInventoryReserved, c.handleReserved); err != nil {
		return err
	}
	return c.client.Consume("order.inventory.failed", mq.RoutingInventoryFailed, c.handleFailed)
}

func (c *Consumer) handleReserved(ctx context.Context, _ string, body []byte) error {
	var evt InventoryResultEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}
	return c.repo.UpdateStatus(evt.OrderNo, model.OrderStatusConfirmed)
}

func (c *Consumer) handleFailed(ctx context.Context, _ string, body []byte) error {
	var evt InventoryResultEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}
	c.logger.Warn("inventory failed", zap.String("order_no", evt.OrderNo), zap.String("message", evt.Message))
	if err := c.repo.UpdateStatus(evt.OrderNo, model.OrderStatusFailed); err != nil {
		return err
	}
	if c.rdb == nil {
		return nil
	}
	order, err := c.repo.FindByOrderNo(evt.OrderNo)
	if err != nil {
		c.logger.Warn("load order for redis restore failed", zap.String("order_no", evt.OrderNo), zap.Error(err))
		return nil
	}
	items := make([]cache.StockItem, 0, len(order.Items))
	for _, it := range order.Items {
		if it.SkuID > 0 && it.Quantity > 0 {
			items = append(items, cache.StockItem{SkuID: it.SkuID, Quantity: it.Quantity})
		}
	}
	if len(items) == 0 {
		return nil
	}
	if err := cache.StockRestore(ctx, c.rdb, items); err != nil {
		c.logger.Warn("redis stock restore failed", zap.String("order_no", evt.OrderNo), zap.Error(err))
	}
	return nil
}
