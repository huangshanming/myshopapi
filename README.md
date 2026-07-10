# mymall 微服务商城

Go 1.24 Monorepo，包含 user / catalog / order 三个独立微服务，通过 APISIX 网关统一入口。

## 架构

```
Client → APISIX (JWT) → user-service     (HTTP :8888, gRPC :9090)
                       → catalog-service (HTTP :8888, gRPC :9090, Redis, MQ)
                       → order-service   (HTTP :8888, gRPC client, MQ)
```

- **同步调用**：order-service 通过 gRPC 调用 user / catalog
- **异步事件**：RabbitMQ `mymall.events` exchange（下单、库存 Saga）
- **数据库**：本机 MySQL，`user_db` / `catalog_db` / `order_db` 各服务独立库
- **缓存**：本机 Redis（catalog 商品列表）
- **K8s Pod 访问本机**：`host.docker.internal`

## 目录结构

```
mymall/
├── api/proto/           # gRPC 定义
├── api/gen/             # protoc 生成代码
├── pkg/                 # 共享库（config/jwt/database/cache/mq/log/health）
├── common/              # LocalTime 等
├── services/
│   ├── user-service/
│   ├── catalog-service/
│   └── order-service/
├── deploy/
│   ├── k8s/             # Deployment / Service / Secret
│   ├── apisix/          # 路由 + JWT
│   └── local/           # docker-compose（Redis + RabbitMQ）
└── scripts/
    ├── migrate-db.sql
    ├── init-order-tables.sql
    └── dev-run.sh
```

## 本地开发

### 1. 基础设施

```bash
# MySQL 建库
mysql -u homestead -p < scripts/migrate-db.sql
mysql -u homestead -p order_db < scripts/init-order-tables.sql

# Redis + RabbitMQ（可选 docker-compose）
docker compose -f deploy/local/docker-compose.infra.yaml up -d
```

### 2. 一键启动（Docker 打包）

```bash
bash scripts/start-all.sh
```

| 服务 | 地址 |
|------|------|
| user-service | http://localhost:8881 |
| catalog-service | http://localhost:8882 |
| order-service | http://localhost:8883 |

停止：`bash scripts/stop-all.sh`

仅重新构建镜像：`bash scripts/build-all.sh`

### 3. 本地热调试（改代码立即生效）

Docker 每次改代码都要重建镜像，**日常开发请用 `dev.sh`**：

```bash
# 只调试正在改的服务（推荐）
bash scripts/dev.sh order

# 或三个服务一起跑
bash scripts/dev.sh
```

改完代码 → `Ctrl+C` → 再执行一次，立刻看到效果。

### 3. K8s 部署

```bash
bash deploy/k8s/apply.sh
bash deploy/apisix/install.sh   # 首次安装 APISIX
bash deploy/apisix/apply.sh
kubectl port-forward svc/apisix-gateway -n mymall 9080:80
```

### 4. 全链路验证（经 APISIX :9080）

注册 → 登录 → 浏览商品 → 下单 → 查订单 → 取消

## gRPC 代码生成

```bash
bash scripts/generate-proto.sh
```
