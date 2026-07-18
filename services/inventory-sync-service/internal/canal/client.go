package canal

import (
	"context"
	"fmt"
	"time"

	"mymall/pkg/config"
	"mymall/services/inventory-sync-service/internal/sync"

	"github.com/withlin/canal-go/client"
	"go.uber.org/zap"
)

type Consumer struct {
	cfg     config.CanalConfig
	handler *sync.Handler
	logger  *zap.Logger
}

func NewConsumer(cfg config.CanalConfig, handler *sync.Handler, logger *zap.Logger) *Consumer {
	return &Consumer{cfg: cfg, handler: handler, logger: logger}
}

// Run connects to Canal and blocks until ctx is cancelled.
func (c *Consumer) Run(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := c.loop(ctx); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			c.logger.Warn("canal loop exited, reconnecting", zap.Error(err))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(3 * time.Second):
			}
		}
	}
}

func (c *Consumer) loop(ctx context.Context) error {
	connector := client.NewSimpleCanalConnector(
		c.cfg.Host,
		c.cfg.Port,
		c.cfg.Username,
		c.cfg.Password,
		c.cfg.Destination,
		int32(c.cfg.SoTimeout),
		int32(c.cfg.IdleTimeout),
	)
	if err := connector.Connect(); err != nil {
		return fmt.Errorf("canal connect: %w", err)
	}
	defer func() {
		_ = connector.DisConnection()
	}()

	if err := connector.Subscribe(c.cfg.Filter); err != nil {
		return fmt.Errorf("canal subscribe: %w", err)
	}
	c.logger.Info("canal subscribed",
		zap.String("host", c.cfg.Host),
		zap.Int("port", c.cfg.Port),
		zap.String("destination", c.cfg.Destination),
		zap.String("filter", c.cfg.Filter),
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		message, err := connector.Get(100, nil, nil)
		if err != nil {
			return fmt.Errorf("canal get: %w", err)
		}
		if message.Id == -1 || len(message.Entries) == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(300 * time.Millisecond):
			}
			continue
		}
		c.handler.HandleEntries(ctx, message.Entries)
	}
}
