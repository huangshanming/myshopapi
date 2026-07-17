#!/usr/bin/env bash
# 本地热更新：air 监听 .go/.yaml，保存后自动重新编译启动（无需 Docker/K8s）
#
# 用法:
#   bash scripts/dev.sh              # 四个服务一起（air）
#   bash scripts/dev.sh user         # 只跑 user-service
#   bash scripts/dev.sh catalog
#   bash scripts/dev.sh order
#   bash scripts/dev.sh merchant
#
# 前提: 本机 MySQL 已启动；脚本会自动拉起 Redis + RabbitMQ 容器
# 首次会 go install github.com/air-verse/air@latest
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
INFRA_FILE="$ROOT/deploy/local/docker-compose.infra.yaml"

TARGET="${1:-all}"

yellow() { printf '\033[33m%s\033[0m\n' "$*"; }
green()  { printf '\033[32m%s\033[0m\n' "$*"; }
red()    { printf '\033[31m%s\033[0m\n' "$*"; }

ensure_air() {
  local gobin
  gobin="$(go env GOPATH)/bin"
  export PATH="$gobin:$PATH"
  if command -v air >/dev/null 2>&1; then
    return 0
  fi
  yellow "==> 安装 air（https://github.com/air-verse/air）"
  go install github.com/air-verse/air@latest
  if ! command -v air >/dev/null 2>&1; then
    red "air 安装后仍找不到，请把 $gobin 加入 PATH 后重试"
    exit 1
  fi
  green "==> air 已就绪: $(command -v air)"
}

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
  yellow "==> $name  →  air 热更新（改 .go / .yaml 自动重启）"
  green "    http://localhost:$port/healthz"
  green "    配置: services/${dir}/.air.toml"
  cd "$ROOT/services/$dir"
  exec air -c .air.toml
}

ensure_air
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
  merchant)
    run_service "merchant-service" "merchant-service" "8884"
    ;;
  all)
    yellow "==> 启动四个服务（air 热更新）"
    green "    user-service     → http://localhost:8881"
    green "    catalog-service  → http://localhost:8882"
    green "    order-service    → http://localhost:8883"
    green "    merchant-service → http://localhost:8884"
    echo "    单独调试: bash scripts/dev.sh <user|catalog|order|merchant>"
    echo "    保存 .go 文件后会自动重编；Ctrl+C 停止全部"
    echo ""
    PIDS=()
    cleanup() { for pid in "${PIDS[@]}"; do kill "$pid" 2>/dev/null || true; done; }
    trap cleanup EXIT INT TERM
    (cd "$ROOT/services/user-service"     && air -c .air.toml) & PIDS+=($!)
    (cd "$ROOT/services/catalog-service"  && air -c .air.toml) & PIDS+=($!)
    (cd "$ROOT/services/order-service"    && air -c .air.toml) & PIDS+=($!)
    (cd "$ROOT/services/merchant-service" && air -c .air.toml) & PIDS+=($!)
    wait
    ;;
  *)
    red "未知参数: $TARGET"
    echo "用法: bash scripts/dev.sh [user|catalog|order|merchant|all]"
    exit 1
    ;;
esac
