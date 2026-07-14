#!/usr/bin/env bash
set -euo pipefail

echo "==> 添加 APISIX Helm 仓库"
helm repo add apisix https://charts.apiseven.com 2>/dev/null || true
helm repo update

echo "==> 安装 APISIX（namespace: mymall）"
helm upgrade --install apisix apisix/apisix \
  --namespace mymall \
  --create-namespace \
  --set ingress-controller.enabled=true \
  --set ingress-controller.config.apisix.serviceNamespace=mymall \
  --set ingress-controller.gatewayProxy.createDefault=true \
  --wait --timeout 10m

echo "==> APISIX 已安装。网关 Service："
kubectl get svc -n mymall | grep apisix || kubectl get svc -n ingress-apisix 2>/dev/null || true
echo "本地访问: kubectl port-forward svc/apisix-gateway -n mymall 9080:80"
