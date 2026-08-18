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
		rowChange := &pbe.RowChange{}
		if err := proto.Unmarshal(entry.GetStoreValue(), rowChange); err != nil {
			h.logger.Warn("unmarshal row change", zap.Error(err))
			continue
		}
		eventType := rowChange.GetEventType()
		for _, row := range rowChange.GetRowDatas() {
			switch table {
			case "product_skus":
				h.handleSku(ctx, eventType, row)
			case "lottery_prizes":
				h.handleLotteryPrize(ctx, eventType, row)
			}
		}
	}
}

func (h *Handler) handleSku(ctx context.Context, eventType pbe.EventType, row *pbe.RowData) {
	switch eventType {
	case pbe.EventType_INSERT:
		h.applySkuInsert(ctx, row.GetAfterColumns())
	case pbe.EventType_UPDATE:
		h.applySkuUpdate(ctx, row.GetAfterColumns())
	case pbe.EventType_DELETE:
		h.applySkuDelete(ctx, row.GetBeforeColumns())
	}
}

func (h *Handler) handleLotteryPrize(ctx context.Context, eventType pbe.EventType, row *pbe.RowData) {
	switch eventType {
	case pbe.EventType_INSERT:
		h.applyLotteryInsert(ctx, row.GetAfterColumns())
	case pbe.EventType_UPDATE:
		h.applyLotteryUpdate(ctx, row.GetAfterColumns())
	case pbe.EventType_DELETE:
		h.applyLotteryDelete(ctx, row.GetBeforeColumns())
	}
}

func (h *Handler) applySkuInsert(ctx context.Context, cols []*pbe.Column) {
	skuID, stock, ok := parseIDStock(cols)
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

// applySkuUpdate always SETs Redis to binlog after-image (MySQL is source of truth for incremental sync).
func (h *Handler) applySkuUpdate(ctx context.Context, after []*pbe.Column) {
	skuID, afterStock, ok := parseIDStock(after)
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

func (h *Handler) applySkuDelete(ctx context.Context, before []*pbe.Column) {
	skuID, _, ok := parseIDStock(before)
	if !ok {
		return
	}
	if err := cache.StockDel(ctx, h.rdb, skuID); err != nil {
		h.logger.Warn("stock del", zap.Uint64("sku_id", skuID), zap.Error(err))
		return
	}
	h.logger.Info("stock delete from binlog", zap.Uint64("sku_id", skuID))
}

func (h *Handler) applyLotteryInsert(ctx context.Context, cols []*pbe.Column) {
	prizeID, stock, ok := parseIDStock(cols)
	if !ok {
		return
	}
	if err := cache.LotteryStockSet(ctx, h.rdb, prizeID, stock); err != nil {
		h.logger.Warn("lottery stock set on insert", zap.Uint64("prize_id", prizeID), zap.Error(err))
		return
	}
	h.logger.Info("lottery stock insert from binlog", zap.Uint64("prize_id", prizeID), zap.Int("stock", stock))
}

func (h *Handler) applyLotteryUpdate(ctx context.Context, after []*pbe.Column) {
	prizeID, afterStock, ok := parseIDStock(after)
	if !ok {
		return
	}
	if err := cache.LotteryStockSet(ctx, h.rdb, prizeID, afterStock); err != nil {
		h.logger.Warn("lottery stock set on update", zap.Uint64("prize_id", prizeID), zap.Error(err))
		return
	}
	h.logger.Info("lottery stock sync from binlog",
		zap.Uint64("prize_id", prizeID),
		zap.Int("stock", afterStock),
	)
}

func (h *Handler) applyLotteryDelete(ctx context.Context, before []*pbe.Column) {
	prizeID, _, ok := parseIDStock(before)
	if !ok {
		return
	}
	if err := cache.LotteryStockDel(ctx, h.rdb, prizeID); err != nil {
		h.logger.Warn("lottery stock del", zap.Uint64("prize_id", prizeID), zap.Error(err))
		return
	}
	h.logger.Info("lottery stock delete from binlog", zap.Uint64("prize_id", prizeID))
}

func parseIDStock(cols []*pbe.Column) (id uint64, stock int, ok bool) {
	var hasID, hasStock bool
	for _, c := range cols {
		switch strings.ToLower(c.GetName()) {
		case "id":
			v, err := strconv.ParseUint(c.GetValue(), 10, 64)
			if err != nil {
				return 0, 0, false
			}
			id = v
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
	return id, stock, hasID && hasStock
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
