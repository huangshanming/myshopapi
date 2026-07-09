# K8s 微服务部署说明

## 架构

```text
user-service Pod     → 本机 MySQL user_db     (host.docker.internal:3306)
catalog-service Pod  → 本机 MySQL catalog_db (host.docker.internal:3306)
```

MySQL **不在 K8s 内运行**，使用本机数据库；上线阿里云时只需改 ConfigMap 中 `mysql.host` 为 RDS 地址。

## 目录结构

```text
deploy/k8s/
├── namespace.yaml
├── secrets/mymall-jwt-auth.yaml
├── services/user-service/      # Deployment + Service + ConfigMap
└── services/catalog-service/
```

## 配置说明

| 类型 | 资源 | 内容 |
|---|---|---|
| **Secret** | `mymall-mysql` | host、username、password（敏感） |
| **Secret** | `mymall-jwt-auth` | JWT secret、consumer key（敏感） |
| **ConfigMap** | `user-service-config` | 端口、dbname、连接池等非敏感配置 |
| **ConfigMap** | `catalog-service-config` | 同上，dbname 为 `catalog_db` |

Deployment 通过环境变量 `MYMALL_MYSQL_*`、`MYMALL_JWT_*` 注入 Secret，覆盖 ConfigMap 中的空值（[`pkg/config/config.go`](../../pkg/config/config.go) 已支持 viper env 覆盖）。

## 首次准备

1. 在本机 MySQL 执行数据库拆分：

```bash
mysql -u homestead -p < scripts/migrate-db.sql
```

2. 确保 MySQL 允许 Docker 访问（`bind-address = 0.0.0.0`，用户有 `@'%'` 权限）

## 部署

```bash
kubectl get nodes
bash deploy/k8s/apply.sh
bash deploy/apisix/apply.sh   # APISIX 已安装时
```

## 本机开发（不用 K8s）

```bash
bash scripts/dev-run.sh
# user-service    → http://localhost:8881
# catalog-service → http://localhost:8882
```

## 验证服务隔离

```bash
kubectl port-forward svc/user-service 8881:8888 -n mymall
curl http://localhost:8881/healthz
curl http://localhost:8881/api/v1/products/list    # 404

kubectl port-forward svc/catalog-service 8882:8888 -n mymall
curl http://localhost:8882/api/v1/products/list    # 200
curl http://localhost:8882/api/v1/user/login       # 404
```

## 上线阿里云 RDS

1. 更新 Secret（不改 ConfigMap）：

```bash
kubectl create secret generic mymall-mysql -n mymall \
  --from-literal=host=rm-xxxxx.mysql.rds.aliyuncs.com \
  --from-literal=username=mymall \
  --from-literal=password='your-rds-password' \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl rollout restart deployment/user-service deployment/catalog-service -n mymall
```

2. 镜像推送到 ACR 后更新 Deployment 中 `image` 字段。
