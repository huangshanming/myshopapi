package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// Lottery prize stock mirrors lottery_prizes.stock (MySQL is source of truth via Canal).
// stock < 0 (unlimited) → key deleted / absent; finite stock → integer value.
const LotteryPrizeStockKeyPrefix = "lottery:prize:stock:"

func LotteryPrizeStockKey(prizeID uint64) string {
	return fmt.Sprintf("%s%d", LotteryPrizeStockKeyPrefix, prizeID)
}

func LotteryStockSet(ctx context.Context, rdb *redis.Client, prizeID uint64, stock int) error {
	if rdb == nil {
		return ErrRedisUnavailable
	}
	if stock < 0 {
		return LotteryStockDel(ctx, rdb, prizeID)
	}
	return rdb.Set(ctx, LotteryPrizeStockKey(prizeID), stock, 0).Err()
}

func LotteryStockDel(ctx context.Context, rdb *redis.Client, prizeID uint64) error {
	if rdb == nil {
		return ErrRedisUnavailable
	}
	return rdb.Del(ctx, LotteryPrizeStockKey(prizeID)).Err()
}

// LotteryStockPreheatSet warms one prize stock without inflating in-flight Redis deducts.
// mysqlStock < 0 (unlimited) → delete key.
func LotteryStockPreheatSet(ctx context.Context, rdb *redis.Client, prizeID uint64, mysqlStock int) (PreheatResult, error) {
	if rdb == nil {
		return PreheatKept, ErrRedisUnavailable
	}
	if mysqlStock < 0 {
		if err := LotteryStockDel(ctx, rdb, prizeID); err != nil {
			return PreheatKept, err
		}
		return PreheatFilled, nil
	}
	n, err := stockPreheatScript.Run(ctx, rdb, []string{LotteryPrizeStockKey(prizeID)}, mysqlStock).Int()
	if err != nil {
		return PreheatKept, err
	}
	return PreheatResult(n), nil
}

// LotteryStockDeduct atomically decrements prize stock by quantity (finite prizes only).
func LotteryStockDeduct(ctx context.Context, rdb *redis.Client, prizeID uint64, quantity int) error {
	if rdb == nil {
		return ErrRedisUnavailable
	}
	if prizeID == 0 || quantity <= 0 {
		return fmt.Errorf("invalid lottery stock deduct prize=%d qty=%d", prizeID, quantity)
	}
	n, err := stockDeductScript.Run(ctx, rdb, []string{LotteryPrizeStockKey(prizeID)}, quantity).Int()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrStockInsufficient
	}
	return nil
}

// LotteryStockRestore increments prize stock (compensation after MySQL finalize failure).
func LotteryStockRestore(ctx context.Context, rdb *redis.Client, prizeID uint64, quantity int) error {
	if rdb == nil {
		return ErrRedisUnavailable
	}
	if prizeID == 0 || quantity <= 0 {
		return nil
	}
	return rdb.IncrBy(ctx, LotteryPrizeStockKey(prizeID), int64(quantity)).Err()
}

// LotteryStockGet returns redis stock; ok=false when key missing (unlimited or not preheated).
func LotteryStockGet(ctx context.Context, rdb *redis.Client, prizeID uint64) (stock int, ok bool, err error) {
	if rdb == nil {
		return 0, false, ErrRedisUnavailable
	}
	v, err := rdb.Get(ctx, LotteryPrizeStockKey(prizeID)).Result()
	if err == redis.Nil {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false, err
	}
	return n, true, nil
}
