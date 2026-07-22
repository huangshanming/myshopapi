#!/usr/bin/env bash
# Generate go-zero sqlx models from DDL (optional; user-service currently uses hand models + db tags).
# Usage: ./scripts/gen-model.sh <ddl.sql> <out-dir>
#   e.g. ./scripts/gen-model.sh scripts/ddl/user-tables.sql services/user-service/internal/modelgen
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DDL="${1:-}"
OUT="${2:-}"
HOME_TPL="${GOCTL_HOME:-$ROOT/deploy/goctl-template}"

if [[ -z "$DDL" || -z "$OUT" ]]; then
  echo "usage: $0 <ddl.sql> <out-dir>" >&2
  exit 2
fi

if ! command -v goctl >/dev/null 2>&1; then
  echo "goctl not found; install: go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5" >&2
  exit 1
fi

mkdir -p "$OUT"
goctl model mysql ddl -src "$DDL" -dir "$OUT" -style go_zero --home "$HOME_TPL"
echo "models generated -> $OUT"
