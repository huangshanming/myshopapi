#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "本机启动 user-service :8881 和 catalog-service :8882"
echo "请确保 MySQL 已执行 scripts/migrate-db.sql"

cd "$ROOT/services/user-service"
go run ./cmd &
USER_PID=$!

cd "$ROOT/services/catalog-service"
go run ./cmd &
CATALOG_PID=$!

trap "kill $USER_PID $CATALOG_PID 2>/dev/null" EXIT

wait
