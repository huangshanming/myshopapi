package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"mymall/pkg/cache"
	"mymall/pkg/mq"
	"mymall/services/order-service/internal/client/merchantrpc"
	"mymall/services/order-service/internal/client/userrpc"
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
	client       *mq.Client
	repo         *repository.OrderRepository
	rdb          *redis.Client
	userHTTP     *userrpc.Client
	merchantHTTP *merchantrpc.Client
	logger       *zap.Logger
}

func NewConsumer(client *mq.Client, repo *repository.OrderRepository, rdb *redis.Client, userHTTP *userrpc.Client, merchantHTTP *merchantrpc.Client, logger *zap.Logger) *Consumer {
	return &Consumer{client: client, repo: repo, rdb: rdb, userHTTP: userHTTP, merchantHTTP: merchantHTTP, logger: logger}
}

func payAmount(o *model.Order) float64 {
	if o.PayAmount > 0 {
		return o.PayAmount
	}
	return o.TotalAmount
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
	order, err := c.repo.FindByOrderNo(ctx, evt.OrderNo)
	if err != nil {
		c.logger.Warn("load order for wallet settle failed", zap.String("order_no", evt.OrderNo), zap.Error(err))
		return err
	}
	// 先实扣再改状态：Settle 幂等；失败则返回 err 让 MQ 重试，避免「已确认却无实扣」
	amt := payAmount(order)
	if c.userHTTP != nil && amt > 0 {
		if err := c.userHTTP.Settle(ctx, order.UserID, amt, order.ID, order.OrderNo); err != nil {
			c.logger.Warn("wallet settle failed", zap.String("order_no", evt.OrderNo), zap.Error(err))
			return err
		}
	}
	if err := c.repo.UpdateStatus(ctx, evt.OrderNo, model.OrderStatusConfirmed); err != nil {
		return err
	}
	if c.merchantHTTP != nil {
		_ = c.merchantHTTP.OrderGiftCoupons(ctx, order.UserID, order.ShopID)
	}
	if c.userHTTP != nil {
		extra, _ := json.Marshal(map[string]interface{}{"order_no": order.OrderNo})
		_ = c.userHTTP.Notify(ctx, userrpc.NotifyReq{
			UserID: order.UserID, Title: "订单已支付成功",
			Content:  fmt.Sprintf("订单 %s 已支付成功，商家将尽快发货", order.OrderNo),
			MsgType:  "order",
			LinkType: "order",
			LinkID:   order.ID,
			Extra:    string(extra),
		})
		_ = c.userHTTP.TaskEvent(ctx, userrpc.TaskEventReq{
			UserID: order.UserID, TaskCode: "place_order", Delta: 1,
			RefType: "order", RefID: order.ID,
		})
	}
	return nil
}

func (c *Consumer) handleFailed(ctx context.Context, _ string, body []byte) error {
	var evt InventoryResultEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return err
	}
	c.logger.Warn("inventory failed", zap.String("order_no", evt.OrderNo), zap.String("message", evt.Message))
	if err := c.repo.UpdateStatus(ctx, evt.OrderNo, model.OrderStatusFailed); err != nil {
		return err
	}
	order, err := c.repo.FindByOrderNo(ctx, evt.OrderNo)
	if err != nil {
		c.logger.Warn("load order for restore failed", zap.String("order_no", evt.OrderNo), zap.Error(err))
		return nil
	}
	if c.userHTTP != nil {
		if err := c.userHTTP.Unfreeze(ctx, order.UserID, payAmount(order), order.ID, order.OrderNo); err != nil {
			c.logger.Warn("wallet unfreeze failed", zap.String("order_no", evt.OrderNo), zap.Error(err))
		}
	}
	if c.merchantHTTP != nil && order.UserCouponID > 0 {
		_ = c.merchantHTTP.UnlockCoupon(ctx, order.UserCouponID, order.ID)
	}
	if c.userHTTP != nil {
		extra, _ := json.Marshal(map[string]interface{}{"order_no": order.OrderNo})
		_ = c.userHTTP.Notify(ctx, userrpc.NotifyReq{
			UserID: order.UserID, Title: "订单已取消",
			Content:  fmt.Sprintf("订单 %s 因库存等原因未能完成，已取消", order.OrderNo),
			MsgType:  "order",
			LinkType: "order",
			LinkID:   order.ID,
			Extra:    string(extra),
		})
	}
	if c.rdb == nil {
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
