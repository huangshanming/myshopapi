package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestLotteryStockPreheatAndDeduct(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	ctx := context.Background()

	res, err := LotteryStockPreheatSet(ctx, rdb, 9, 5)
	if err != nil || res != PreheatFilled {
		t.Fatalf("fill: res=%v err=%v", res, err)
	}
	if mustGet(t, mr, "lottery:prize:stock:9") != "5" {
		t.Fatalf("want 5 got %s", mustGet(t, mr, "lottery:prize:stock:9"))
	}

	if err := LotteryStockDeduct(ctx, rdb, 9, 1); err != nil {
		t.Fatal(err)
	}
	if mustGet(t, mr, "lottery:prize:stock:9") != "4" {
		t.Fatalf("want 4 got %s", mustGet(t, mr, "lottery:prize:stock:9"))
	}

	_ = mr.Set("lottery:prize:stock:9", "0")
	if err := LotteryStockDeduct(ctx, rdb, 9, 1); !errors.Is(err, ErrStockInsufficient) {
		t.Fatalf("want insufficient got %v", err)
	}

	// unlimited: delete key
	res, err = LotteryStockPreheatSet(ctx, rdb, 9, -1)
	if err != nil || res != PreheatFilled {
		t.Fatalf("unlimited: res=%v err=%v", res, err)
	}
	if mr.Exists("lottery:prize:stock:9") {
		t.Fatal("unlimited should delete key")
	}
}
