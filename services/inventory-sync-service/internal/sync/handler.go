package sync

import (
	"context"
	"strconv"
	"strings"

	"mymall/pkg/cache"

	pbe "github.com/withlin/canal-go/protocol/entry"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	rdb    *redis.Client
	logger *zap.Logger
}

func NewHandler(rdb *redis.Client, logger *zap.Logger) *Handler {
	return &Handler{rdb: rdb, logger: logger}
}

func (h *Handler) HandleEntries(ctx context.Context, entries []pbe.Entry) {
	for _, entry := range entries {
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN ||
			entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		header := entry.GetHeader()
		table := strings.ToLower(header.GetTableName())
		if table != "product_skus" {
			continue
		}
		rowChange := &pbe.RowChange{}
		if err := proto.Unmarshal(entry.GetStoreValue(), rowChange); err != nil {
			h.logger.Warn("unmarshal row change", zap.Error(err))
			continue
		}
		eventType := rowChange.GetEventType()
		for _, row := range rowChange.GetRowDatas() {
			switch eventType {
			case pbe.EventType_INSERT:
				h.applyInsert(ctx, row.GetAfterColumns())
			case pbe.EventType_UPDATE:
				h.applyUpdate(ctx, row.GetAfterColumns())
			case pbe.EventType_DELETE:
				h.applyDelete(ctx, row.GetBeforeColumns())
			}
		}
	}
}

func (h *Handler) applyInsert(ctx context.Context, cols []*pbe.Column) {
	skuID, stock, ok := parseSkuStock(cols)
	if !ok {
		return
	}
	if softDeleted(cols) {
		_ = cache.StockDel(ctx, h.rdb, skuID)
		return
	}
	if err := cache.StockSet(ctx, h.rdb, skuID, stock); err != nil {
		h.logger.Warn("stock set on insert", zap.Uint64("sku_id", skuID), zap.Error(err))
		return
	}
	h.logger.Info("stock insert from binlog", zap.Uint64("sku_id", skuID), zap.Int("stock", stock))
}

// applyUpdate always SETs Redis to binlog after-image (MySQL is source of truth for incremental sync).
// Startup preheat stays conservative; Canal updates must not be skipped by CAS.
func (h *Handler) applyUpdate(ctx context.Context, after []*pbe.Column) {
	skuID, afterStock, ok := parseSkuStock(after)
	if !ok {
		return
	}
	if softDeleted(after) {
		_ = cache.StockDel(ctx, h.rdb, skuID)
		h.logger.Info("stock del soft-delete from binlog", zap.Uint64("sku_id", skuID))
		return
	}
	if err := cache.StockSet(ctx, h.rdb, skuID, afterStock); err != nil {
		h.logger.Warn("stock set on update", zap.Uint64("sku_id", skuID), zap.Error(err))
		return
	}
	h.logger.Info("stock sync from binlog",
		zap.Uint64("sku_id", skuID),
		zap.Int("stock", afterStock),
	)
}

func (h *Handler) applyDelete(ctx context.Context, before []*pbe.Column) {
	skuID, _, ok := parseSkuStock(before)
	if !ok {
		return
	}
	if err := cache.StockDel(ctx, h.rdb, skuID); err != nil {
		h.logger.Warn("stock del", zap.Uint64("sku_id", skuID), zap.Error(err))
		return
	}
	h.logger.Info("stock delete from binlog", zap.Uint64("sku_id", skuID))
}

func parseSkuStock(cols []*pbe.Column) (skuID uint64, stock int, ok bool) {
	var hasID, hasStock bool
	for _, c := range cols {
		switch strings.ToLower(c.GetName()) {
		case "id":
			v, err := strconv.ParseUint(c.GetValue(), 10, 64)
			if err != nil {
				return 0, 0, false
			}
			skuID = v
			hasID = true
		case "stock":
			v, err := strconv.Atoi(c.GetValue())
			if err != nil {
				return 0, 0, false
			}
			stock = v
			hasStock = true
		}
	}
	return skuID, stock, hasID && hasStock
}

func softDeleted(cols []*pbe.Column) bool {
	for _, c := range cols {
		if strings.ToLower(c.GetName()) == "deleted_at" {
			v := strings.TrimSpace(c.GetValue())
			return v != "" && !strings.EqualFold(v, "null")
		}
	}
	return false
}
