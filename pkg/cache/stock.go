package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const SkuStockKeyPrefix = "catalog:sku:stock:"

var (
	ErrStockInsufficient = errors.New("库存不足")
	ErrRedisUnavailable  = errors.New("redis unavailable")
)

// StockItem is a SKU quantity pair for deduct/restore.
type StockItem struct {
	SkuID    uint64
	Quantity int
}

func SkuStockKey(skuID uint64) string {
	return fmt.Sprintf("%s%d", SkuStockKeyPrefix, skuID)
}

// stockDeductScript: all-or-nothing check then DECRBY.
// KEYS = sku keys, ARGV = quantities (same order).
// Returns 1 on success, 0 if any key missing or stock insufficient.
var stockDeductScript = redis.NewScript(`
for i = 1, #KEYS do
  local v = redis.call('GET', KEYS[i])
  if (not v) then
    return 0
  end
  local n = tonumber(v)
  local q = tonumber(ARGV[i])
  if (not n) or (not q) or n < q then
    return 0
  end
end
for i = 1, #KEYS do
  redis.call('DECRBY', KEYS[i], ARGV[i])
end
return 1
`)

// stockCASSetScript: SET after only if key missing or GET == before.
// KEYS[1]=key, ARGV[1]=before, ARGV[2]=after
// Returns 1 if set, 0 if skipped.
var stockCASSetScript = redis.NewScript(`
local cur = redis.call('GET', KEYS[1])
if (not cur) or cur == ARGV[1] then
  redis.call('SET', KEYS[1], ARGV[2])
  return 1
end
return 0
`)

// stockPreheatScript: safe warm-up against concurrent Redis-first deduct.
// KEYS[1]=key, ARGV[1]=mysql stock
// Returns: 1 filled (missing), 2 pulled down (redis > mysql), 0 kept (redis <= mysql).
var stockPreheatScript = redis.NewScript(`
local mysql = tonumber(ARGV[1])
if (not mysql) or mysql < 0 then
  mysql = 0
end
local cur = redis.call('GET', KEYS[1])
if (not cur) then
  redis.call('SET', KEYS[1], mysql)
  return 1
end
local n = tonumber(cur)
if (not n) or n > mysql then
  redis.call('SET', KEYS[1], mysql)
  return 2
end
return 0
`)

// PreheatResult is the outcome of StockPreheatSet.
type PreheatResult int

const (
	PreheatKept       PreheatResult = 0 // redis exists and <= mysql (protect in-flight deduct)
	PreheatFilled     PreheatResult = 1 // key was missing
	PreheatPulledDown PreheatResult = 2 // redis was higher than mysql
)

func StockSet(ctx context.Context, rdb *redis.Client, skuID uint64, stock int) error {
	if rdb == nil {
		return ErrRedisUnavailable
	}
	if stock < 0 {
		stock = 0
	}
	return rdb.Set(ctx, SkuStockKey(skuID), stock, 0).Err()
}

func StockDel(ctx context.Context, rdb *redis.Client, skuID uint64) error {
	if rdb == nil {
		return ErrRedisUnavailable
	}
	return rdb.Del(ctx, SkuStockKey(skuID)).Err()
}

// StockCASSet sets after when key is absent or equals before (string form of int).
func StockCASSet(ctx context.Context, rdb *redis.Client, skuID uint64, before, after int) (bool, error) {
	if rdb == nil {
		return false, ErrRedisUnavailable
	}
	if after < 0 {
		after = 0
	}
	n, err := stockCASSetScript.Run(ctx, rdb, []string{SkuStockKey(skuID)},
		strconv.Itoa(before), strconv.Itoa(after)).Int()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

// StockPreheatSet warms one SKU without inflating in-flight Redis deducts:
// missing → SET mysql; redis > mysql → pull down; redis <= mysql → keep.
func StockPreheatSet(ctx context.Context, rdb *redis.Client, skuID uint64, mysqlStock int) (PreheatResult, error) {
	if rdb == nil {
		return PreheatKept, ErrRedisUnavailable
	}
	if mysqlStock < 0 {
		mysqlStock = 0
	}
	n, err := stockPreheatScript.Run(ctx, rdb, []string{SkuStockKey(skuID)}, mysqlStock).Int()
	if err != nil {
		return PreheatKept, err
	}
	return PreheatResult(n), nil
}

// StockDeduct atomically decrements all SKUs or fails with ErrStockInsufficient.
func StockDeduct(ctx context.Context, rdb *redis.Client, items []StockItem) error {
	if rdb == nil {
		return ErrRedisUnavailable
	}
	if len(items) == 0 {
		return nil
	}
	keys := make([]string, 0, len(items))
	args := make([]interface{}, 0, len(items))
	for _, it := range items {
		if it.SkuID == 0 || it.Quantity <= 0 {
			return fmt.Errorf("invalid stock item sku=%d qty=%d", it.SkuID, it.Quantity)
		}
		keys = append(keys, SkuStockKey(it.SkuID))
		args = append(args, it.Quantity)
	}
	n, err := stockDeductScript.Run(ctx, rdb, keys, args...).Int()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrStockInsufficient
	}
	return nil
}

// StockRestore increments SKU stocks (compensation).
func StockRestore(ctx context.Context, rdb *redis.Client, items []StockItem) error {
	if rdb == nil {
		return ErrRedisUnavailable
	}
	pipe := rdb.Pipeline()
	for _, it := range items {
		if it.SkuID == 0 || it.Quantity <= 0 {
			continue
		}
		pipe.IncrBy(ctx, SkuStockKey(it.SkuID), int64(it.Quantity))
	}
	_, err := pipe.Exec(ctx)
	return err
}
