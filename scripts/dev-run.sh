#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "本机启动 user-service :8881、catalog-service :8882、order-service :8883"
echo "请确保 MySQL 已执行 scripts/migrate-db.sql 与 scripts/init-order-tables.sql"
echo "Redis/RabbitMQ 可通过 deploy/local/docker-compose.infra.yaml 启动"

PIDS=()

cleanup() {
  for pid in "${PIDS[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
}
trap cleanup EXIT

(cd "$ROOT/services/user-service" && CONFIG_PATH=./etc/user-service.yaml go run ./cmd) &
PIDS+=($!)

(cd "$ROOT/services/catalog-service" && CONFIG_PATH=./etc/catalog-service.yaml go run ./cmd) &
PIDS+=($!)

(cd "$ROOT/services/order-service" && CONFIG_PATH=./etc/order-service.yaml go run ./cmd) &
PIDS+=($!)

wait
