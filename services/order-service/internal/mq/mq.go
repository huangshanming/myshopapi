package mq

import (
	"context"
	"encoding/json"

	"mymall/pkg/mq"
	"mymall/services/order-service/internal/model"
	"mymall/services/order-service/internal/repository"

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
	logger *zap.Logger
}

func NewConsumer(client *mq.Client, repo *repository.OrderRepository, logger *zap.Logger) *Consumer {
	return &Consumer{client: client, repo: repo, logger: logger}
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
	return c.repo.UpdateStatus(evt.OrderNo, model.OrderStatusFailed)
}
