package mq

import (
	"context"
	"encoding/json"

	"mymall/pkg/mq"
	"mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"

	"go.uber.org/zap"
)

type stockItem struct {
	ProductID uint64 `json:"product_id"`
	SkuID     uint64 `json:"sku_id"`
	Quantity  int    `json:"quantity"`
}

type OrderCreatedEvent struct {
	OrderNo string      `json:"order_no"`
	Items   []stockItem `json:"items"`
}

type OrderCancelledEvent struct {
	OrderNo string      `json:"order_no"`
	Items   []stockItem `json:"items"`
}

type InventoryResultEvent struct {
	OrderNo string `json:"order_no"`
	Message string `json:"message,omitempty"`
}

type Consumer struct {
	client *mq.Client
	logic  *logic.CatalogLogic
	logger *zap.Logger
}

func NewConsumer(client *mq.Client, l *logic.CatalogLogic, logger *zap.Logger) *Consumer {
	return &Consumer{client: client, logic: l, logger: logger}
}

func (c *Consumer) Start() error {
	if err := c.client.Consume("catalog.order.created", mq.RoutingOrderCreated, c.handleOrderCreated); err != nil {
		return err
	}
	return c.client.Consume("catalog.order.cancelled", mq.RoutingOrderCancelled, c.handleOrderCancelled)
}

func (c *Consumer) handleOrderCreated(ctx context.Context, _ string, body []byte) error {
	var evt OrderCreatedEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}
	items := toRepoItems(evt.Items)
	result := mq.RoutingInventoryReserved
	payload := InventoryResultEvent{OrderNo: evt.OrderNo}
	if err := c.logic.ReserveStock(items); err != nil {
		c.logger.Warn("reserve stock failed", zap.String("order_no", evt.OrderNo), zap.Error(err))
		result = mq.RoutingInventoryFailed
		payload.Message = err.Error()
	}
	return c.client.Publish(ctx, result, payload)
}

func (c *Consumer) handleOrderCancelled(ctx context.Context, _ string, body []byte) error {
	var evt OrderCancelledEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}
	if err := c.logic.ReleaseStock(toRepoItems(evt.Items)); err != nil {
		c.logger.Warn("release stock failed", zap.String("order_no", evt.OrderNo), zap.Error(err))
		return err
	}
	return nil
}

func toRepoItems(items []stockItem) []repository.StockItem {
	out := make([]repository.StockItem, 0, len(items))
	for _, it := range items {
		out = append(out, repository.StockItem{ProductID: it.ProductID, SkuID: it.SkuID, Quantity: it.Quantity})
	}
	return out
}
