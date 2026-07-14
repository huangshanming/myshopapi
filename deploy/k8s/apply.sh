#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

echo "==> 创建 namespace"
kubectl apply -f "$ROOT/deploy/k8s/namespace.yaml"

echo "==> 创建 Secret"
kubectl apply -f "$ROOT/deploy/k8s/secrets/mymall-jwt-auth.yaml"
kubectl apply -f "$ROOT/deploy/k8s/secrets/mymall-mysql.yaml"
kubectl apply -f "$ROOT/deploy/k8s/secrets/mymall-redis.yaml"
kubectl apply -f "$ROOT/deploy/k8s/secrets/mymall-rabbitmq.yaml"

echo "==> 清理旧资源（单体 / K8s MySQL）"
kubectl delete deployment mymall-api -n mymall --ignore-not-found
kubectl delete deployment mysql -n mymall --ignore-not-found
kubectl delete svc mysql -n mymall --ignore-not-found
kubectl delete configmap mymall-config -n mymall --ignore-not-found

echo "==> 构建微服务镜像"
docker build -t mymall-user-service:local -f "$ROOT/services/user-service/Dockerfile" "$ROOT"
docker build -t mymall-catalog-service:local -f "$ROOT/services/catalog-service/Dockerfile" "$ROOT"
docker build -t mymall-order-service:local -f "$ROOT/services/order-service/Dockerfile" "$ROOT"
docker build -t mymall-merchant-service:local -f "$ROOT/services/merchant-service/Dockerfile" "$ROOT"

echo "==> 部署 user-service"
kubectl apply -f "$ROOT/deploy/k8s/services/user-service/"

echo "==> 部署 catalog-service"
kubectl apply -f "$ROOT/deploy/k8s/services/catalog-service/"

echo "==> 部署 order-service"
kubectl apply -f "$ROOT/deploy/k8s/services/order-service/"

echo "==> 部署 merchant-service"
kubectl apply -f "$ROOT/deploy/k8s/services/merchant-service/"

echo "==> 等待 Pod 就绪"
kubectl rollout status deployment/user-service -n mymall --timeout=180s
kubectl rollout status deployment/catalog-service -n mymall --timeout=180s
kubectl rollout status deployment/order-service -n mymall --timeout=180s
kubectl rollout status deployment/merchant-service -n mymall --timeout=180s

echo "==> 完成。当前资源："
kubectl get pods,svc -n mymall
