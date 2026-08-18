#!/usr/bin/env bash
# Smoke mall-uni / admin critical HTTP paths against local go-zero services.
# Expect: no 5xx. Auth-gated routes may return 401/403.
set -euo pipefail

USER_BASE="${USER_BASE:-http://127.0.0.1:8881}"
CATALOG_BASE="${CATALOG_BASE:-http://127.0.0.1:8882}"
ORDER_BASE="${ORDER_BASE:-http://127.0.0.1:8883}"
MERCHANT_BASE="${MERCHANT_BASE:-http://127.0.0.1:8884}"
LOTTERY_BASE="${LOTTERY_BASE:-http://127.0.0.1:8887}"

fail=0
ok=0

check() {
  local method="$1" url="$2" body="${3:-}"
  local tmp code
  tmp="$(mktemp)"
  if [[ "$method" == "POST" ]]; then
    code=$(curl -sS -o "$tmp" -w '%{http_code}' -X POST -H 'Content-Type: application/json' \
      -H 'X-User-Id: 1' -H 'X-User-Role: user' \
      -d "${body}" "$url" || echo "000")
  else
    code=$(curl -sS -o "$tmp" -w '%{http_code}' \
      -H 'X-User-Id: 1' -H 'X-User-Role: user' \
      "$url" || echo "000")
  fi
  if [[ "$code" =~ ^5 ]] || [[ "$code" == "000" ]]; then
    echo "FAIL $code $method $url"
    head -c 240 "$tmp"; echo
    fail=$((fail + 1))
  else
    # pagination: if body has "list", it should be an array when present as top-level list
    if python3 - "$tmp" <<'PY' 2>/dev/null
import json,sys
p=sys.argv[1]
try:
  d=json.load(open(p))
except Exception:
  sys.exit(0)
if isinstance(d, dict) and "list" in d and d["list"] is not None and not isinstance(d["list"], list):
  # allow opaque maps for seckill-style only if nested list exists
  inner=d["list"]
  if isinstance(inner, dict) and isinstance(inner.get("list"), list):
    print("WARN nested list shape:", p)
    sys.exit(2)
  print("BAD list not array")
  sys.exit(1)
sys.exit(0)
PY
    then
      :
    else
      st=$?
      if [[ $st -eq 1 ]]; then
        echo "FAIL shape $method $url (list not array)"
        head -c 200 "$tmp"; echo
        fail=$((fail + 1))
        rm -f "$tmp"
        return
      fi
    fi
    echo "OK   $code $method $url"
    ok=$((ok + 1))
  fi
  rm -f "$tmp"
}

echo "== catalog =="
check GET "$CATALOG_BASE/api/v1/banners"
check GET "$CATALOG_BASE/api/v1/products/list?page=1&page_size=5"
check GET "$CATALOG_BASE/api/v1/products/detail?id=12"
check GET "$CATALOG_BASE/api/v1/products/sales-rank?page=1&page_size=3"
check GET "$CATALOG_BASE/api/v1/product_category/list?page=1&page_size=10"
check GET "$CATALOG_BASE/api/v1/articles/list?page=1&page_size=5"
check GET "$CATALOG_BASE/uploads/products/1/lipstick1.png"

echo "== merchant =="
check GET "$MERCHANT_BASE/api/v1/shops/list?page=1&page_size=5&home=1"
check GET "$MERCHANT_BASE/api/v1/shops/home-slots?slot_type=brand_shop"
check GET "$MERCHANT_BASE/api/v1/home/theme-tiles"
check GET "$MERCHANT_BASE/api/v1/seckill/current"
check GET "$MERCHANT_BASE/api/v1/seckill/list?page=1&page_size=5"
check GET "$MERCHANT_BASE/api/v1/coupons/center"

echo "== user =="
check GET "$USER_BASE/api/v1/user/notifications/unread-count"
check GET "$USER_BASE/api/v1/regions"
check POST "$USER_BASE/api/v1/user/tasks/events" '{"task_code":"browse_products","ref_type":"product","ref_id":"12"}'
check GET "$USER_BASE/api/v1/user/tasks"
check GET "$USER_BASE/uploads/points-mall/placeholder.png"

# Field-name contract: user-service modelgen must emit snake_case (not PascalCase).
assert_snake() {
  local url="$1" path="$2" expect_keys="$3"
  local tmp code
  tmp="$(mktemp)"
  code=$(curl -sS -o "$tmp" -w '%{http_code}' \
    -H 'X-User-Id: 1' -H 'X-User-Role: admin' "$url" || echo "000")
  if [[ "$code" =~ ^5 ]] || [[ "$code" == "000" ]]; then
    echo "FAIL $code GET $url (snake assert)"
    fail=$((fail + 1))
    rm -f "$tmp"
    return
  fi
  if python3 - "$tmp" "$path" "$expect_keys" <<'PY'
import json, sys
d = json.load(open(sys.argv[1]))
path = sys.argv[2].split(".") if sys.argv[2] else []
cur = d
for p in path:
  if p == "0" and isinstance(cur, list):
    cur = cur[0] if cur else {}
  elif isinstance(cur, dict):
    cur = cur.get(p)
  else:
    cur = None
    break
if not isinstance(cur, dict):
  print("no object at", sys.argv[2])
  sys.exit(1)
for k in sys.argv[3].split(","):
  if k not in cur:
    print("missing", k, "got", sorted(cur.keys()))
    sys.exit(1)
# PascalCase regressions
for bad in ("ID", "Mobile", "Nickname", "CreatedAt", "FrozenBalance", "ConfigKey"):
  if bad in cur:
    print("pascal field", bad)
    sys.exit(1)
sys.exit(0)
PY
  then
    echo "OK   $code snake $url"
    ok=$((ok + 1))
  else
    echo "FAIL snake $url"
    head -c 200 "$tmp"; echo
    fail=$((fail + 1))
  fi
  rm -f "$tmp"
}
assert_snake "$USER_BASE/api/v1/admin/users?page=1&page_size=1" "list.0" "id,mobile,nickname,status,role"
assert_snake "$USER_BASE/api/v1/user/wallet" "" "user_id,balance,frozen_balance"
assert_snake "$USER_BASE/api/v1/admin/configs" "list.0" "config_key,config_value"

echo "== order =="
check GET "$ORDER_BASE/api/v1/orders?page=1&page_size=5"
check GET "$ORDER_BASE/api/v1/logistics"

echo "== lottery =="
check GET "$LOTTERY_BASE/healthz"
check GET "$LOTTERY_BASE/api/v1/lottery/activity"

echo
echo "passed=$ok failed=$fail"
[[ "$fail" -eq 0 ]]
