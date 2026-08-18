# mymall 微服务商城

Go 1.24 Monorepo，user / catalog / order / merchant 微服务，**HTTP 为 go-zero rest + goctl**，服务间 **zrpc**，数据层 **sqlx**，经 APISIX 统一入口。

## 架构

```
Client → APISIX (JWT) → user-service     (rest :8881, zrpc :9090)
                       → catalog-service (rest :8882, zrpc :9091, Redis, MQ)
                       → order-service   (rest :8883, zrpc Client → user/catalog/merchant)
                       → merchant-service(rest :8884, zrpc :9092)
                       → inventory-sync  (health :8885, Canal→Redis 库存预热)
                       → agent-service   (FastAPI :8886, Python Agents / 导购)
                       → lottery-service (HTTP :8887)
                       → recommend-service (FastAPI :8888, 协同推荐)
                       → askdata-service (FastAPI :8889, 电商问数)
```

- **HTTP**：go-zero `rest`；契约见 [`docs/gozero-http-contract.md`](docs/gozero-http-contract.md)
- **配置**：每服务 `internal/config` 嵌入 `rest.RestConf`，`conf.MustLoad` + `MYMALL_*` 覆盖
- **同步调用**：order → user/catalog/merchant **zrpc**（goctl rpc：`api/gen` + `api/rpcclient`；可选 etcd）
- **边缘网关**：**APISIX**（JWT / 路由 / `X-User-*`；刻意不换成 go-zero gateway）
- **异步事件**：RabbitMQ `mymall.events`
- **数据层**：四业务服务 **sqlx**；goctl model → `internal/modelgen`；schema 靠 `scripts/*.sql` + `scripts/ddl/`
- **库存**：Redis 预扣 + MQ + MySQL；`inventory-sync-service` 经 Canal 同步 Redis
- **goctl**：见 [`docs/gozero-migration-runbook.md`](docs/gozero-migration-runbook.md)

## 目录结构

```
mymall/
├── api/proto/           # gRPC 定义（user / catalog / merchant）
├── api/gen/             # goctl/protoc 生成 pb
├── api/rpcclient/       # goctl 生成的 zrpc client
├── pkg/                 # 共享库（xerr/jwt/middleware/cache/mq/zrpcx…；database 仅 inventory-sync）
├── common/              # LocalTime 等
├── apps/admin-web/      # 管理端 / 商家端 Vue3
├── apps/mall-uni/       # 用户端 UniApp
├── services/<svc>/
│   ├── api/*.api        # goctl HTTP 源
│   ├── etc/<svc>.yaml   # RestConf + 业务字段（CONFIG_PATH）
│   ├── internal/modelgen/ # goctl model 实体
│   ├── internal/rpclogic/ # goctl rpc 业务（与 HTTP logic 分离）
│   ├── cmd/main.go      # conf.MustLoad + rest + 可选 zrpc
│   └── internal/
│       ├── config/      # RestConf 嵌入
│       ├── handler/     # goctl 生成
│       ├── logic/       # 业务逻辑
│       ├── types/
│       ├── svc/         # ServiceContext（sqlx.Conn / Redis / MQ / RPC）
│       ├── model/       # 表实体（db tag）
│       ├── repository/  # sqlx 数据访问
│       ├── server/      # zrpc 服务端（user/catalog/merchant）
│       └── client/      # zrpc 客户端
├── services/agent-service/  # Python FastAPI Agents（:8886，智能导购等）
├── services/recommend-service/  # Python FastAPI 协同推荐（:8888）
├── services/askdata-service/  # Python FastAPI 电商问数（:8889，uv）
├── services/lottery-service/ # 九宫格抽奖（:8887）
├── deploy/
│   ├── k8s/
│   ├── apisix/
│   └── local/docker-compose.yaml
└── scripts/
    ├── gen-api.sh / generate-proto.sh / gen-model.sh
    ├── init-*.sql / alter-*.sql / seed-*.sql
    └── dev.sh
```

## 本地开发

### 1. 基础设施

```bash
# MySQL 建库 + 建表（无 migrate 工具，按脚本顺序手动执行；均可重复跑）
mysql -u homestead -p < scripts/migrate-db.sql
mysql -u homestead -p mymall < scripts/init-schema.sql
mysql -u homestead -p mymall < scripts/init-order-tables.sql
# 后台多租户：补 role / shop_id，建 shops 等表
mysql -u homestead -p mymall < scripts/init-merchant-tables.sql
# 种子账号（超管 13900000001 / 商家 13900000002，密码均为 123456）
mysql -u homestead -p mymall < scripts/seed-admin-merchant.sql
# 首页商户（小米/生鲜/服饰 + 门头图；先落盘图片）
bash scripts/seed-home-shop-images.sh
mysql -u homestead -p mymall < scripts/seed-home-shops.sql
# 秒杀场次 + 商家钱包
mysql -u homestead -p --default-character-set=utf8mb4 mymall < scripts/alter-seckill-wallet.sql
# 订单项秒杀关联
mysql -u homestead -p --default-character-set=utf8mb4 mymall < scripts/alter-order-seckill.sql
# 后台 RBAC（菜单/角色/权限）
mysql -u homestead -p mymall < scripts/init-rbac-tables.sql
mysql -u homestead -p --default-character-set=utf8mb4 mymall < scripts/seed-rbac.sql

# 可选：商品中台表结构 + 演示商品（images/ 素材）
# mysql -u homestead -p mymall < scripts/alter-product-center.sql
# python3 scripts/seed-shop-demo-products.py

# Redis + RabbitMQ + Canal + Milvus（可选 docker-compose）
docker compose -f deploy/local/docker-compose.infra.yaml up -d
# Milvus gRPC :19530 / health :19091

# Canal 复制账号（一次即可；需本机 MySQL 已开 ROW binlog）
mysql -u root -p < scripts/init-canal-mysql.sql
```

若基础表早已建好，**补** `init-merchant-tables.sql` + `seed-admin-merchant.sql`；后台 RBAC 再补 `init-rbac-tables.sql` + `seed-rbac.sql`。改代码后需重建对应服务镜像。

### Canal / 库存预热

本机 MySQL 需开启 ROW binlog（Homestead / 自建均可），例如 `my.cnf`：

```ini
[mysqld]
server-id=1
log-bin=mysql-bin
binlog_format=ROW
binlog_row_image=FULL
```

检查：

```bash
mysql -u homestead -p -e "SHOW VARIABLES LIKE 'log_bin'; SHOW VARIABLES LIKE 'binlog_format';"
```

`inventory-sync-service` 启动时全量预热 `product_skus.stock` 与 `lottery_prizes.stock` → Redis，随后用 [canal-go](https://github.com/withlin/canal-go) 增量同步。下单路径：Redis 预扣 → MQ → MySQL 乐观锁；抽奖有限库存：Redis 预扣 → MySQL 事务扣库存（binlog 回写 Redis）。`inventory.failed` / 取消订单会补偿 Redis。

```bash
bash scripts/dev.sh inventory   # 单独跑 sync
# 或
bash scripts/dev.sh             # 含 inventory-sync
```

验证：`redis-cli KEYS 'catalog:sku:stock:*'`；管理端改 SKU 库存后 Redis 应跟随。

表结构变更请执行 `scripts/*.sql`（四业务服务已不用 GORM AutoMigrate）。

### 2. 一键启动（Docker 打包）

```bash
bash scripts/start-all.sh
```

| 服务 | 地址 |
|------|------|
| user-service | http://localhost:8881 |
| catalog-service | http://localhost:8882 |
| order-service | http://localhost:8883 |
| merchant-service | http://localhost:8884 |
| inventory-sync-service | http://localhost:8885 |
| admin-web（Vue） | `cd apps/admin-web && npm i && npm run dev` → http://localhost:5174 |
| mall-uni（用户端） | `cd apps/mall-uni && yarn && yarn dev:h5` → http://localhost:5175 |

停止：`bash scripts/stop-all.sh`

仅重新构建镜像：`bash scripts/build-all.sh`

**Docker 拉取基础镜像 429？** `start-all.sh` / `build-all.sh` 已默认走 DaoCloud 镜像源；若仍失败，可手动指定：

```bash
export GOLANG_IMAGE=docker.m.daocloud.io/library/golang:1.24-alpine
export ALPINE_IMAGE=docker.m.daocloud.io/library/alpine:3.19
bash scripts/start-all.sh
```

或在 Docker Desktop → Settings → Docker Engine 中移除限流的 registry mirror（如 `docker.xuanyuan.me`）。

### 3. 本地热更新（air，无需 K8s）

Docker 每次改代码都要重建镜像，**日常开发请用 `dev.sh`**（[air](https://github.com/air-verse/air)，保存 `.go` / `.yaml` 自动重编重启）：

```bash
# 只调试正在改的服务（推荐）
bash scripts/dev.sh order

# 或全部服务一起跑（含 inventory-sync）
bash scripts/dev.sh
```

首次会自动 `go install github.com/air-verse/air@latest`（需 `$(go env GOPATH)/bin` 在 PATH）。改共享包 `pkg/` / `common/` 也会触发对应服务重建。前端仍用 `cd apps/admin-web && npm run dev`（Vite HMR）。

### 3. K8s 部署

```bash
bash deploy/k8s/apply.sh
bash deploy/apisix/install.sh   # 首次安装 APISIX
bash deploy/apisix/apply.sh
kubectl port-forward svc/apisix-gateway -n mymall 9080:80
```

### 4. 全链路验证（经 APISIX :9080）

注册 → 登录 → 浏览商品 → 下单 → 查订单 → 取消

## 接口文档（OpenAPI + Apifox）

### 生成文档

```bash
bash scripts/gen-docs.sh
```

产出：
- `docs/openapi/mymall.yaml` — **整份复制给 AI 生成前端页面**
- `docs/openapi/mymall.swagger.json` — Swagger 2.0 备份

### 浏览文档

**先在一个终端启动服务（保持运行不要关）：**

```bash
bash scripts/serve-docs.sh
```

**再另开浏览器访问：** http://localhost:9099/scalar/index.html

若提示「连接被拒绝」，说明文档服务没启动或终端已关闭。

**不要**直接 `open docs/scalar/index.html`（`file://` 无法加载 yaml）。

各服务开发时也可访问 Swagger UI：`http://localhost:8881/swagger/index.html`

### Apifox 导入与调试

1. 打开 Apifox → **导入** → **OpenAPI** → 选择 `docs/openapi/mymall.yaml`
2. 新建环境 **本地开发**：
   - `baseUrl` = `http://localhost:9080`（经 APISIX）或 `http://localhost:8881`（直连 user）
   - 变量 `token` = 登录后获得的 JWT
3. 在「Auth」配置 **Bearer Token** → `{{token}}`
4. 调试流程：调用 `POST /api/v1/user/login` → 复制 token → 再调需鉴权接口
5. 每次改接口注释后执行 `bash scripts/gen-docs.sh`，Apifox **重新导入** 同步

### 给 AI 生成页面的模板

```text
以下是商城 OpenAPI 3.0 规范，请生成 Vue3 页面：
- 统一响应 { code, msg, data }，code=200 为成功
- JWT: Authorization: Bearer {token}

[粘贴 docs/openapi/mymall.yaml 全文]
```

可选自动生成 TS 类型：

```bash
npx openapi-typescript docs/openapi/mymall.yaml -o frontend/src/types/api.ts
```

## gRPC 代码生成

```bash
bash scripts/generate-proto.sh
```
