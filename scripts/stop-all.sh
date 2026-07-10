#!/usr/bin/env bash
# 停止本机 Docker 栈（微服务 + Redis + RabbitMQ）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
COMPOSE_FILE="$ROOT/deploy/local/docker-compose.yaml"

if ! docker info >/dev/null 2>&1; then
  echo "Docker 未运行，无需停止"
  exit 0
fi

echo "停止 Docker 服务栈..."
docker compose -f "$COMPOSE_FILE" down
echo "已停止"
