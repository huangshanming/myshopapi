package cache

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func mustGet(t *testing.T, mr *miniredis.Miniredis, key string) string {
	t.Helper()
	v, err := mr.Get(key)
	if err != nil {
		t.Fatalf("get %s: %v", key, err)
	}
	return v
}

func TestStockPreheatSet(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	ctx := context.Background()

	res, err := StockPreheatSet(ctx, rdb, 1, 100)
	if err != nil || res != PreheatFilled {
		t.Fatalf("fill: res=%v err=%v", res, err)
	}
	if mustGet(t, mr, "catalog:sku:stock:1") != "100" {
		t.Fatalf("want 100 got %s", mustGet(t, mr, "catalog:sku:stock:1"))
	}

	_ = mr.Set("catalog:sku:stock:1", "90")
	res, err = StockPreheatSet(ctx, rdb, 1, 100)
	if err != nil || res != PreheatKept {
		t.Fatalf("keep: res=%v err=%v", res, err)
	}
	if mustGet(t, mr, "catalog:sku:stock:1") != "90" {
		t.Fatalf("want keep 90 got %s", mustGet(t, mr, "catalog:sku:stock:1"))
	}

	_ = mr.Set("catalog:sku:stock:1", "120")
	res, err = StockPreheatSet(ctx, rdb, 1, 100)
	if err != nil || res != PreheatPulledDown {
		t.Fatalf("pull: res=%v err=%v", res, err)
	}
	if mustGet(t, mr, "catalog:sku:stock:1") != "100" {
		t.Fatalf("want 100 got %s", mustGet(t, mr, "catalog:sku:stock:1"))
	}
}
