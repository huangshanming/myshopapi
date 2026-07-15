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

# 同 tag :local 时：1) Docker 易全 CACHED；2) kind/Docker Desktop 节点 containerd 不会自动看到宿主机 docker build 的新镜像
# 注意：set -u 下空数组 "${arr[@]}" 在部分 bash 会报 unbound variable
NO_CACHE=""
if [[ "${FORCE_REBUILD:-}" == "1" ]] || [[ "${1:-}" == "--no-cache" ]]; then
  NO_CACHE=1
  echo "==> 强制无缓存构建 (--no-cache)"
fi
CACHEBUST="${CACHEBUST:-$(date +%s)}"
IMAGE_TAG="local-${CACHEBUST}"

# Docker Desktop / kind 节点名即 docker 容器名（如 desktop-control-plane）
cluster_node() {
  kubectl get nodes -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true
}

# 把刚 build 的镜像灌进集群 containerd，否则 Pod 会一直用旧 digest（表现为 CrashLoop / 旧 Gin 日志）
load_image_to_cluster() {
  local image="$1"
  local node
  node="$(cluster_node)"
  if [[ -z "$node" ]] || ! docker inspect "$node" >/dev/null 2>&1; then
    echo "WARN: 无法向节点导入镜像（找不到节点容器 ${node:-none}），跳过 $image"
    return 0
  fi
  echo "==> 导入 $image → 节点 $node"
  docker save "$image" | docker exec -i "$node" ctr -n k8s.io images import -
}

build_svc() {
  local name="$1"
  local image="mymall-${name}:${IMAGE_TAG}"
  echo "==> 构建 ${image} (+ :local) CACHEBUST=$CACHEBUST"
  local -a build_cmd=(docker build)
  if [[ -n "$NO_CACHE" ]]; then
    build_cmd+=(--no-cache)
  fi
  build_cmd+=(
    --build-arg "CACHEBUST=$CACHEBUST"
    -t "$image"
    -t "mymall-${name}:local"
    -f "$ROOT/services/${name}/Dockerfile"
    "$ROOT"
  )
  "${build_cmd[@]}"
  load_image_to_cluster "$image"
}

echo "==> 构建微服务镜像并导入集群"
build_svc user-service
build_svc catalog-service
build_svc order-service
build_svc merchant-service

echo "==> 部署 Manifest"
kubectl apply -f "$ROOT/deploy/k8s/services/user-service/"
kubectl apply -f "$ROOT/deploy/k8s/services/catalog-service/"
kubectl apply -f "$ROOT/deploy/k8s/services/order-service/"
kubectl apply -f "$ROOT/deploy/k8s/services/merchant-service/"

echo "==> 切换到本次构建的镜像标签 ${IMAGE_TAG}"
kubectl set image deployment/user-service user-service="mymall-user-service:${IMAGE_TAG}" -n mymall
kubectl set image deployment/catalog-service catalog-service="mymall-catalog-service:${IMAGE_TAG}" -n mymall
kubectl set image deployment/order-service order-service="mymall-order-service:${IMAGE_TAG}" -n mymall
kubectl set image deployment/merchant-service merchant-service="mymall-merchant-service:${IMAGE_TAG}" -n mymall

# 清除此前 ProgressDeadlineExceeded，避免仍卡在旧状态
for d in user-service catalog-service order-service merchant-service; do
  kubectl patch deployment "$d" -n mymall --type merge -p '{"spec":{"progressDeadlineSeconds":600}}' >/dev/null
done

echo "==> 等待 Pod 就绪"
kubectl rollout status deployment/user-service -n mymall --timeout=300s
kubectl rollout status deployment/catalog-service -n mymall --timeout=300s
kubectl rollout status deployment/order-service -n mymall --timeout=300s
kubectl rollout status deployment/merchant-service -n mymall --timeout=300s

echo "==> 完成。当前资源："
kubectl get pods,svc -n mymall
echo "镜像标签: ${IMAGE_TAG}（下次 apply 会换新 tag）"
