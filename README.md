# mymall 微服务商城

Go 1.24 Monorepo，user / catalog / order / merchant 微服务，**HTTP 框架为 go-zero rest**，服务间 **zrpc**（基于现有 proto），经 APISIX 统一入口。

## 架构

```
Client → APISIX (JWT) → user-service     (go-zero rest :8888, zrpc :9090)
                       → catalog-service (go-zero rest :8888, zrpc :9090, Redis, MQ)
                       → order-service   (go-zero rest :8888, zrpc Client, MQ)
                       → merchant-service(go-zero rest :8888)
```

- **HTTP**：go-zero `rest`（已替换 Gin）
- **同步调用**：order-service 通过 **zrpc** 调用 user / catalog（沿用 `api/proto` 生成代码）
- **异步事件**：RabbitMQ `mymall.events` exchange（下单、库存 Saga）
- **数据库**：本机 MySQL，四服务共用 **`mymall`** 库（本地开发简化；上线可再拆库）
- **缓存**：本机 Redis（catalog 商品列表）
- **K8s Pod 访问本机**：`host.docker.internal`
- **接口文档**：`bash scripts/serve-docs.sh`（Scalar）；服务内不再挂 gin-swagger

## 目录结构

```
mymall/
├── api/proto/           # gRPC 定义
├── api/gen/             # protoc 生成代码
├── pkg/                 # 共享库（config/jwt/database/httpserver/middleware…）
├── common/              # LocalTime 等
├── apps/admin-web/      # 管理端 / 商家端 Vue3
├── services/<svc>/      # goctl 风格目录（四服务同构）
│   ├── etc/<svc>.yaml   # 本地/镜像默认配置（CONFIG_PATH）
│   ├── cmd/main.go
│   └── internal/
│       ├── handler/     # HTTP 薄入口
│       ├── logic/       # 业务逻辑
│       ├── types/       # HTTP Req/Resp
│       ├── svc/         # ServiceContext（DB/Redis/MQ/RPC）
│       ├── model/       # 表模型；model/sql 指向 scripts/*.sql
│       ├── repository/  # 数据访问（迁移期内保留）
│       ├── server/      # zrpc 服务端（user/catalog）
│       └── client/      # zrpc 客户端（order: userrpc/catalogrpc）
├── deploy/
│   ├── k8s/             # Deployment / Service / Secret；ConfigMap 挂载 /app/etc/<svc>.yaml
│   ├── apisix/          # 路由 + JWT
│   └── local/           # docker-compose（Redis + RabbitMQ + 四服务）
└── scripts/
    ├── migrate-db.sql
    ├── init-schema.sql
    ├── init-order-tables.sql
    ├── init-merchant-tables.sql
    ├── seed-admin-merchant.sql
    ├── init-rbac-tables.sql
    ├── seed-rbac.sql
    └── dev.sh           # go run；默认读 ./etc/<svc>.yaml
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
# 后台 RBAC（菜单/角色/权限）
mysql -u homestead -p mymall < scripts/init-rbac-tables.sql
mysql -u homestead -p mymall < scripts/seed-rbac.sql

# 可选：商品中台表结构 + 演示商品（images/ 素材）
# mysql -u homestead -p mymall < scripts/alter-product-center.sql
# python3 scripts/seed-shop-demo-products.py

# Redis + RabbitMQ（可选 docker-compose）
docker compose -f deploy/local/docker-compose.infra.yaml up -d
```

若基础表早已建好，**补** `init-merchant-tables.sql` + `seed-admin-merchant.sql`；后台 RBAC 再补 `init-rbac-tables.sql` + `seed-rbac.sql`。改代码后需重建对应服务镜像。

本地 `server.mode: debug`（默认）启动时会按各服务 model 做 **GORM AutoMigrate**（加表/加列，类似 Beego sync；**不删列**）。K8s `release` 不会跑。可用 `MYMALL_MYSQL_AUTO_MIGRATE=true|false` 强制开关。种子数据仍靠 SQL，AutoMigrate 不会插账号。

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
| admin-web（Vue） | `cd apps/admin-web && npm i && npm run dev` → http://localhost:5174 |

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

# 或四个服务一起跑
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
