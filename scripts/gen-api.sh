#!/usr/bin/env bash
# Generate go-zero HTTP stubs from api/*.api (routes/handler/logic/middleware).
# Usage: ./scripts/gen-api.sh <service-name>
#   e.g. ./scripts/gen-api.sh merchant-service
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SERVICE="${1:-}"
HOME_TPL="${GOCTL_HOME:-$ROOT/deploy/goctl-template}"

if [[ -z "$SERVICE" ]]; then
  echo "usage: $0 <service-name>" >&2
  exit 2
fi

DIR="$ROOT/services/$SERVICE"
if [[ ! -d "$DIR" ]]; then
  echo "missing service dir: $DIR" >&2
  exit 1
fi

API="$(find "$DIR/api" -name '*.api' | head -1 || true)"
if [[ -z "$API" ]]; then
  echo "missing api/*.api under $DIR" >&2
  exit 1
fi

if ! command -v goctl >/dev/null 2>&1; then
  echo "goctl not found; install: go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5" >&2
  exit 1
fi

# Preserve hand-maintained ServiceContext across goctl (goctl would overwrite).
SVC="$DIR/internal/svc/service_context.go"
SVC_BAK=""
if [[ -f "$SVC" ]]; then
  SVC_BAK="$(mktemp)"
  cp "$SVC" "$SVC_BAK"
fi

echo "==> goctl api go ($SERVICE)"
(
  cd "$DIR"
  # Generate into service root; style go_zero
  goctl api go -api "$API" -dir . -style go_zero --home "$HOME_TPL"
)

# Drop goctl default main/config/etc — project uses cmd/main.go + pkg/config
rm -f "$DIR"/*api.go 2>/dev/null || true
rm -rf "$DIR/internal/config" 2>/dev/null || true
# keep etc/*.yaml service configs; remove only goctl-named etc if present
find "$DIR/etc" -name '*api.yaml' -delete 2>/dev/null || true

if [[ -n "$SVC_BAK" && -f "$SVC_BAK" ]]; then
  cp "$SVC_BAK" "$SVC"
  rm -f "$SVC_BAK"
  echo "==> restored internal/svc/service_context.go"
fi

echo "==> done. Edit logic/middleware as needed; do NOT hand-edit handler/routes.go"
echo "    verify: ./scripts/check-api-routes.sh $SERVICE"
