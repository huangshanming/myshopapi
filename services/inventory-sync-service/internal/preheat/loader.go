package preheat

import (
	"context"
	"fmt"

	"mymall/pkg/cache"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type skuRow struct {
	ID    uint64 `gorm:"column:id"`
	Stock int    `gorm:"column:stock"`
}

// Stats summarizes a safe preheat pass.
type Stats struct {
	Scanned    int // rows read from MySQL
	Filled     int // Redis key missing → SET mysql
	PulledDown int // Redis > mysql → SET mysql (防偏高超卖)
	Kept       int // Redis <= mysql → 保留（保护在途预扣）
}

// LoadAllSkuStock safely warms Redis from MySQL without overwriting lower in-flight stock.
// Hot restart: keeps redis when redis <= mysql; only fills holes or pulls down stale-high keys.
func LoadAllSkuStock(ctx context.Context, db *gorm.DB, rdb *redis.Client, logger *zap.Logger) (Stats, error) {
	var st Stats
	if db == nil || rdb == nil {
		return st, fmt.Errorf("db/redis required")
	}
	const batch = 500
	var lastID uint64
	for {
		var rows []skuRow
		err := db.WithContext(ctx).
			Table("product_skus").
			Select("id, stock").
			Where("deleted_at IS NULL AND id > ?", lastID).
			Order("id ASC").
			Limit(batch).
			Find(&rows).Error
		if err != nil {
			return st, err
		}
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			lastID = row.ID
			st.Scanned++
			res, err := cache.StockPreheatSet(ctx, rdb, row.ID, row.Stock)
			if err != nil {
				return st, err
			}
			switch res {
			case cache.PreheatFilled:
				st.Filled++
			case cache.PreheatPulledDown:
				st.PulledDown++
			default:
				st.Kept++
			}
		}
		if logger != nil {
			logger.Info("preheat batch",
				zap.Int("batch", len(rows)),
				zap.Uint64("last_id", lastID),
				zap.Int("filled", st.Filled),
				zap.Int("pulled_down", st.PulledDown),
				zap.Int("kept", st.Kept),
			)
		}
	}
	return st, nil
}
