#!/usr/bin/env bash
# Docker 打包三个微服务镜像（不启动）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
COMPOSE_FILE="$ROOT/deploy/local/docker-compose.yaml"

export GOLANG_IMAGE="${GOLANG_IMAGE:-docker.m.daocloud.io/library/golang:1.24-alpine}"
export ALPINE_IMAGE="${ALPINE_IMAGE:-docker.m.daocloud.io/library/alpine:3.19}"

if ! docker info >/dev/null 2>&1; then
  echo "Docker 未运行，请先启动 Docker Desktop"
  exit 1
fi

echo "==> Docker 构建镜像"
docker compose -f "$COMPOSE_FILE" build

echo ""
echo "完成。镜像："
docker images | grep mymall- || true
echo ""
echo "启动: bash scripts/start-all.sh --no-build"
echo "K8s:  bash deploy/k8s/apply.sh"
