# K8s 部署说明

## 目录结构

```text
deploy/k8s/
├── namespace.yaml
├── secrets/
│   └── mymall-jwt-auth.yaml
├── infra/mysql/          # MySQL 数据库
├── services/mymall/      # 过渡期单体 Deployment
├── services/user-service/    # Service（指向 mymall-api）
├── services/catalog-service/   # Service（指向 mymall-api）
└── services/order-service/     # Service（占位，指向 mymall-api）
```

## 当前状态（过渡期）

微服务尚未拆分，采用 **一个 Deployment + 多个 Service 名称** 的方式，让 APISIX 路由能正常工作。Phase 1 拆完后再替换为独立 Deployment。

## 部署步骤

```bash
# 1. 确保 K8s 集群 Ready
kubectl get nodes

# 2. 部署 MySQL + 应用
bash deploy/k8s/apply.sh

# 3. 部署 APISIX 路由（需先安装 APISIX Ingress Controller）
bash deploy/apisix/apply.sh
```

## 手动分步部署

```bash
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/secrets/
kubectl apply -f deploy/k8s/infra/mysql/
docker build -t mymall-api:local .
kubectl apply -f deploy/k8s/services/mymall/
kubectl apply -f deploy/k8s/services/user-service/
kubectl apply -f deploy/k8s/services/catalog-service/
kubectl apply -f deploy/k8s/services/order-service/
```

## 验证

```bash
kubectl get pods -n mymall
kubectl port-forward svc/user-service 8888:8888 -n mymall
curl http://localhost:8888/healthz
curl http://localhost:8888/api/v1/products/list
```

## 注意事项

- 镜像名 `mymall-api:local` 需在本地 Docker 中 build，Docker Desktop K8s 可直接使用
- MySQL 首次启动需等待约 30~60 秒
- 数据库表需自行导入或迁移（项目未含 AutoMigrate）
