#!/usr/bin/env bash
# 本机一键启动：Docker 打包并运行 Redis + RabbitMQ + 三个微服务
# 用法: bash scripts/start-all.sh              # 构建镜像 + 启动容器
#       bash scripts/start-all.sh --no-build   # 跳过构建，直接启动已有镜像
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
COMPOSE_FILE="$ROOT/deploy/local/docker-compose.yaml"

# Docker Hub 镜像加速源 429 时，默认走 DaoCloud 拉取基础镜像（可 export 覆盖）
export GOLANG_IMAGE="${GOLANG_IMAGE:-docker.m.daocloud.io/library/golang:1.24-alpine}"
export ALPINE_IMAGE="${ALPINE_IMAGE:-docker.m.daocloud.io/library/alpine:3.19}"

NO_BUILD=false
for arg in "$@"; do
  case "$arg" in
    --no-build|--skip-build) NO_BUILD=true ;;
  esac
done

red()  { printf '\033[31m%s\033[0m\n' "$*"; }
green(){ printf '\033[32m%s\033[0m\n' "$*"; }
yellow(){ printf '\033[33m%s\033[0m\n' "$*"; }

check_docker() {
  yellow "==> 检查 Docker"
  if ! docker info >/dev/null 2>&1; then
    red "Docker 未运行，请先启动 Docker Desktop"
    exit 1
  fi
  green "    Docker 已就绪"
}

check_mysql() {
  yellow "==> 检查 MySQL（容器通过 host.docker.internal 访问）"
  if command -v mysql >/dev/null 2>&1; then
    if mysql -h127.0.0.1 -uhomestead -psecret -e "SELECT 1" >/dev/null 2>&1; then
      green "    MySQL 已就绪 (homestead@127.0.0.1:3306)"
      return 0
    fi
  fi
  yellow "    无法自动连接 MySQL，请确认本机 MySQL/Homestead 已启动"
  yellow "    首次需执行: mysql -u homestead -p < scripts/migrate-db.sql"
  yellow "                mysql -u homestead -p mymall < scripts/init-schema.sql"
  yellow "                mysql -u homestead -p mymall < scripts/init-order-tables.sql"
  yellow "                mysql -u homestead -p mymall < scripts/init-merchant-tables.sql"
  yellow "                mysql -u homestead -p mymall < scripts/seed-admin-merchant.sql"
}

start_stack() {
  yellow "==> Docker 打包并启动全部服务"
  local args=( -f "$COMPOSE_FILE" up -d )
  if ! $NO_BUILD; then
    args+=( --build )
  fi
  docker compose "${args[@]}"

  green "    user-service            → http://localhost:8881  (gRPC :9090)"
  green "    catalog-service         → http://localhost:8882  (gRPC :9091)"
  green "    order-service           → http://localhost:8883"
  green "    merchant-service        → http://localhost:8884"
  green "    inventory-sync-service  → http://localhost:8885"
  green "    Redis                   → 127.0.0.1:6379"
  green "    Canal                   → 127.0.0.1:11111"
  green "    RabbitMQ 管理台         → http://localhost:15672  (mymall/mymall)"
}

wait_healthy() {
  yellow "==> 等待服务就绪"
  local urls=(
    "http://localhost:8881/healthz"
    "http://localhost:8882/healthz"
    "http://localhost:8883/healthz"
    "http://localhost:8884/healthz"
    "http://localhost:8885/healthz"
  )
  for url in "${urls[@]}"; do
    for i in $(seq 1 60); do
      if curl -sf "$url" >/dev/null 2>&1; then
        green "    OK $url"
        break
      fi
      if [[ $i -eq 60 ]]; then
        red "    超时 $url"
        yellow "    查看日志: docker compose -f $COMPOSE_FILE logs"
        if [[ "$url" == *":8883"* ]]; then
          yellow "    order-service 常见原因: mymall 库未建表、MySQL 未启动"
          yellow "    docker compose -f $COMPOSE_FILE logs order-service"
        fi
        exit 1
      fi
      sleep 2
    done
  done
}

main() {
  green "mymall 本机一键启动（Docker）"
  check_docker
  check_mysql
  start_stack
  wait_healthy
  echo ""
  green "全部启动完成！"
  echo "  查看容器: docker compose -f deploy/local/docker-compose.yaml ps"
  echo "  查看日志: docker compose -f deploy/local/docker-compose.yaml logs -f"
  echo "  停止服务: bash scripts/stop-all.sh"
  echo "  快速重启: bash scripts/start-all.sh --no-build"
  echo ""
  echo "  K8s 部署请用: bash deploy/k8s/apply.sh"
}

main "$@"
