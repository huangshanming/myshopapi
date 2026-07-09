#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

echo "==> 创建 namespace"
kubectl apply -f "$ROOT/deploy/k8s/namespace.yaml"

echo "==> 创建 JWT Secret"
kubectl apply -f "$ROOT/deploy/k8s/secrets/mymall-jwt-auth.yaml"

echo "==> 配置 APISIX JWT Consumer"
kubectl apply -f "$ROOT/deploy/apisix/jwt-consumer.yaml"

echo "==> 配置 APISIX 路由（需先 bash deploy/k8s/apply.sh 部署后端）"
kubectl apply -f "$ROOT/deploy/apisix/public-routes.yaml"
kubectl apply -f "$ROOT/deploy/apisix/protected-routes.yaml"

echo "==> 完成。验证 Consumer："
kubectl get apisixconsumer -n mymall
