#!/usr/bin/env bash
# Compare METHOD+path sets between api/*.api and internal/handler/routes.go.
# Also fail if cmd/main.go still hardcodes business API paths.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SERVICE="${1:-}"

usage() {
  echo "usage: $0 <service-name|all>" >&2
  echo "  e.g. $0 merchant-service" >&2
  exit 2
}

[[ -n "$SERVICE" ]] || usage

check_one() {
  local name="$1"
  local dir="$ROOT/services/$name"
  local api routes main
  api="$(find "$dir/api" -name '*.api' 2>/dev/null | head -1 || true)"
  routes="$dir/internal/handler/routes.go"
  # catalog may nest product/content; allow override
  if [[ ! -f "$routes" ]]; then
    routes="$(find "$dir/internal" -path '*/handler/routes.go' 2>/dev/null | head -1 || true)"
  fi
  main="$dir/cmd/main.go"

  if [[ -z "$api" || ! -f "$api" ]]; then
    echo "[$name] FAIL: missing api/*.api" >&2
    return 1
  fi
  if [[ -z "$routes" || ! -f "$routes" ]]; then
    echo "[$name] FAIL: missing handler/routes.go" >&2
    return 1
  fi
  if [[ ! -f "$main" ]]; then
    echo "[$name] FAIL: missing cmd/main.go" >&2
    return 1
  fi

  if grep -E 'Path:\s*"/api/' "$main" >/dev/null 2>&1; then
    echo "[$name] FAIL: cmd/main.go still contains Path: \"/api/..." >&2
    grep -nE 'Path:\s*"/api/' "$main" >&2 || true
    return 1
  fi

  local tmp
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' RETURN

  # .api: lines like "get /api/v1/foo" or "post /healthz"
  python3 - "$api" "$tmp/api.txt" <<'PY'
import re, sys
api, out = sys.argv[1], sys.argv[2]
text = open(api).read()
# strip block comments roughly
text = re.sub(r'//.*?$', '', text, flags=re.M)
pat = re.compile(r'(?im)^\s*(get|post|put|delete|patch)\s+(\S+)')
rows = set()
for m in pat.finditer(text):
    rows.add(f"{m.group(1).upper()} {m.group(2)}")
open(out, 'w').write('\n'.join(sorted(rows)) + ('\n' if rows else ''))
print(len(rows))
PY
  local api_n=$?
  api_count=$(wc -l < "$tmp/api.txt" | tr -d ' ')

  python3 - "$routes" "$tmp/routes.txt" <<'PY'
import re, sys
path = sys.argv[1]
out = sys.argv[2]
text = open(path).read()
# Method: http.MethodGet, Path: "/foo"
pat = re.compile(
    r'Method:\s*http\.Method(Get|Post|Put|Delete|Patch)\s*,\s*Path:\s*"([^"]+)"',
    re.I,
)
rows = set()
for m in pat.finditer(text):
    rows.add(f"{m.group(1).upper()} {m.group(2)}")
open(out, 'w').write('\n'.join(sorted(rows)) + ('\n' if rows else ''))
print(len(rows))
PY
  routes_count=$(wc -l < "$tmp/routes.txt" | tr -d ' ')

  if ! diff -u "$tmp/api.txt" "$tmp/routes.txt" > "$tmp/diff.txt"; then
    echo "[$name] FAIL: .api ($api_count) vs routes.go ($routes_count) mismatch:" >&2
    cat "$tmp/diff.txt" >&2
    return 1
  fi
  echo "[$name] OK: $api_count routes match ($api <-> $(basename "$routes"))"
}

if [[ "$SERVICE" == "all" ]]; then
  ec=0
  for s in merchant-service order-service user-service catalog-service inventory-sync-service; do
    if [[ -d "$ROOT/services/$s" ]]; then
      check_one "$s" || ec=1
    fi
  done
  exit $ec
fi

check_one "$SERVICE"
