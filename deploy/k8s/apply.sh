#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

echo "==> 创建 namespace"
kubectl apply -f "$ROOT/deploy/k8s/namespace.yaml"

echo "==> 创建 Secret"
kubectl apply -f "$ROOT/deploy/k8s/secrets/mymall-jwt-auth.yaml"
kubectl apply -f "$ROOT/deploy/k8s/infra/mysql/secret.yaml"

echo "==> 部署 MySQL"
kubectl apply -f "$ROOT/deploy/k8s/infra/mysql/pvc.yaml"
kubectl apply -f "$ROOT/deploy/k8s/infra/mysql/deployment.yaml"
kubectl apply -f "$ROOT/deploy/k8s/infra/mysql/service.yaml"

echo "==> 构建本地镜像 mymall-api:local"
docker build -t mymall-api:local "$ROOT"

echo "==> 部署 mymall-api（过渡期单体）"
kubectl apply -f "$ROOT/deploy/k8s/services/mymall/configmap.yaml"
kubectl apply -f "$ROOT/deploy/k8s/services/mymall/deployment.yaml"

echo "==> 创建微服务 Service（过渡期均指向 mymall-api）"
kubectl apply -f "$ROOT/deploy/k8s/services/user-service/service.yaml"
kubectl apply -f "$ROOT/deploy/k8s/services/catalog-service/service.yaml"
kubectl apply -f "$ROOT/deploy/k8s/services/order-service/service.yaml"

echo "==> 等待 Pod 就绪"
kubectl rollout status deployment/mysql -n mymall --timeout=180s
kubectl rollout status deployment/mymall-api -n mymall --timeout=120s

echo "==> 完成。当前资源："
kubectl get pods,svc -n mymall
