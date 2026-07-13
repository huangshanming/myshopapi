# K8s 微服务部署说明

## 架构

```text
APISIX Gateway
    ├── user-service     (HTTP :8888, gRPC :9090) → mymall
    ├── catalog-service  (HTTP :8888, gRPC :9090) → mymall + Redis + MQ
    └── order-service    (HTTP :8888)              → mymall + MQ + gRPC client

本机基础设施（Pod 经 host.docker.internal 访问）:
  MySQL :3306 | Redis :6379 | RabbitMQ :5672
```

## 目录结构

```text
deploy/k8s/
├── namespace.yaml
├── secrets/           # mysql / jwt / redis / rabbitmq
├── services/
│   ├── user-service/
│   ├── catalog-service/
│   └── order-service/
└── observability/     # Jaeger 骨架
```

## 首次准备

```bash
# 1. 建库
mysql -u homestead -p < scripts/migrate-db.sql
mysql -u homestead -p mymall < scripts/init-order-tables.sql

# 2. 启动 Redis + RabbitMQ（可选）
docker compose -f deploy/local/docker-compose.infra.yaml up -d

# 3. 生成 gRPC 代码（修改 proto 后）
bash scripts/generate-proto.sh
```

## 部署

```bash
bash deploy/k8s/apply.sh
bash deploy/apisix/install.sh   # 首次安装 APISIX
bash deploy/apisix/apply.sh
kubectl apply -f deploy/k8s/observability/jaeger.yaml   # 可选

kubectl port-forward svc/apisix-gateway -n mymall 9080:80
```

## 本机开发（不用 K8s）

```bash
bash scripts/dev-run.sh
# user-service    → http://localhost:8881  (gRPC :9090)
# catalog-service → http://localhost:8882  (gRPC :9091)
# order-service   → http://localhost:8883
```

下单 API 需携带 `X-User-Id` Header（经 APISIX 时由 JWT 自动注入）。

## 全链路验证

```text
注册 → 登录 → 浏览商品 → 下单 → 查订单 → 取消 → 库存回滚
```

经 APISIX `:9080` 或本机直连各服务端口。

## Secret 说明

| Secret | 内容 |
|---|---|
| `mymall-mysql` | host、username、password |
| `mymall-jwt-auth` | JWT secret、consumer key |
| `mymall-redis` | host、password |
| `mymall-rabbitmq` | host、username、password |

Deployment 通过 `MYMALL_*` 环境变量注入，覆盖 ConfigMap 占位符。
