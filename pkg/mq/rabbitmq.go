package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mymall/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RoutingOrderCreated    = "order.created"
	RoutingOrderCancelled  = "order.cancelled"
	RoutingInventoryReserved = "inventory.reserved"
	RoutingInventoryFailed   = "inventory.failed"
)

type Client struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func New(cfg config.RabbitMQConfig) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}
	if err := ch.ExchangeDeclare(cfg.Exchange, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}
	return &Client{conn: conn, channel: ch, exchange: cfg.Exchange}, nil
}

func (c *Client) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	if c.conn == nil || c.conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection closed")
	}
	return nil
}

func (c *Client) Publish(ctx context.Context, routingKey string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.channel.PublishWithContext(ctx, c.exchange, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	})
}

type Handler func(ctx context.Context, routingKey string, body []byte) error

func (c *Client) Consume(queue, routingKey string, handler Handler) error {
	q, err := c.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err := c.channel.QueueBind(q.Name, routingKey, c.exchange, false, nil); err != nil {
		return err
	}
	if err := c.channel.Qos(1, 0, false); err != nil {
		return err
	}
	deliveries, err := c.channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for d := range deliveries {
			ctx := context.Background()
			if err := handler(ctx, d.RoutingKey, d.Body); err != nil {
				d.Nack(false, true)
				continue
			}
			d.Ack(false)
		}
	}()
	return nil
}
