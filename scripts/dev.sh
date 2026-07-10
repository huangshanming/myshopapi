#!/usr/bin/env bash
# 本地热调试：改代码后 Ctrl+C 再执行即可立即生效（go run，无需 Docker 重建）
#
# 用法:
#   bash scripts/dev.sh              # 启动三个服务
#   bash scripts/dev.sh user         # 只调试 user-service
#   bash scripts/dev.sh catalog      # 只调试 catalog-service
#   bash scripts/dev.sh order        # 只调试 order-service
#
# 前提: 本机 MySQL 已启动；脚本会自动拉起 Redis + RabbitMQ 容器
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
INFRA_FILE="$ROOT/deploy/local/docker-compose.infra.yaml"

TARGET="${1:-all}"

yellow() { printf '\033[33m%s\033[0m\n' "$*"; }
green()  { printf '\033[32m%s\033[0m\n' "$*"; }
red()    { printf '\033[31m%s\033[0m\n' "$*"; }

start_infra() {
  if ! docker info >/dev/null 2>&1; then
    red "Docker 未运行。Redis/RabbitMQ 需要 Docker，请先启动 Docker Desktop"
    exit 1
  fi
  yellow "==> 启动 Redis + RabbitMQ（基础设施）"
  docker compose -f "$INFRA_FILE" up -d
}

run_service() {
  local name="$1"
  local dir="$2"
  local port="$3"
  yellow "==> $name  →  go run（改代码后 Ctrl+C 再执行本脚本即可）"
  green "    http://localhost:$port/healthz"
  cd "$ROOT/services/$dir"
  exec go run ./cmd
}

start_infra

case "$TARGET" in
  user)
    run_service "user-service" "user-service" "8881"
    ;;
  catalog)
    run_service "catalog-service" "catalog-service" "8882"
    ;;
  order)
    run_service "order-service" "order-service" "8883"
    ;;
  all)
    yellow "==> 启动三个服务（go run）"
    green "    user-service    → http://localhost:8881"
    green "    catalog-service → http://localhost:8882"
    green "    order-service   → http://localhost:8883"
    echo "    改某个服务代码后，单独调试: bash scripts/dev.sh <user|catalog|order>"
    echo ""
    PIDS=()
    cleanup() { for pid in "${PIDS[@]}"; do kill "$pid" 2>/dev/null || true; done; }
    trap cleanup EXIT
    cd "$ROOT/services/user-service"    && go run ./cmd & PIDS+=($!)
    cd "$ROOT/services/catalog-service" && go run ./cmd & PIDS+=($!)
    cd "$ROOT/services/order-service"   && go run ./cmd & PIDS+=($!)
    wait
    ;;
  *)
    red "未知参数: $TARGET"
    echo "用法: bash scripts/dev.sh [user|catalog|order|all]"
    exit 1
    ;;
esac
