# goctl / go-zero 工作流

## 工具

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5
export PATH="$(go env GOPATH)/bin:$PATH"
```

自定义模板：[`deploy/goctl-template`](../deploy/goctl-template)

- **handler**：唯一接触 `w`/`r`；`httpx.Parse` → 调 logic → `OkJsonCtx` / `ErrorCtx`
- **logic**：不存 `ctx`、不接收 `w`/`r`；签名 `func (l *XxxLogic) Xxx(ctx context.Context, req *types.XxxReq) (*types.XxxResp, error)`
- **数据层**：业务四服务用 go-zero **sqlx**（`sqlx.NewMysql` + 手写 repository）；entity 由 goctl model 生成到 `internal/modelgen`
- **禁止** `internal/httpapi`；业务在 `logic` / `biz` / domain

## 边界（刻意保留）

| 层 | 方案 |
|---|---|
| **边缘网关** | **APISIX**（JWT、CORS、路径分流、`X-User-*` 注入）— **不用** go-zero gateway |
| **HTTP** | goctl api + rest |
| **RPC** | goctl rpc（pb + server 桩 + `api/rpcclient`）+ [`pkg/zrpcx`](../pkg/zrpcx) |
| **Model** | goctl model → `modelgen`；repository **不**用生成 CRUD |
| **inventory-sync** | Viper + GORM（Canal worker） |

## HTTP：唯一正确流程

1. 编辑服务 `api/*.api`（路由 / `@server` 中间件 / `group` / `@handler` / **type + returns**）
   - **全服务统一**：`group` 必须写成 `<端>/<模块>`
   - 每个带参接口必须定义 `type XxxReq`；列表/详情等需 `returns (XxxResp)`
   - 分页：`page`、`page_size`（见 `pkg/middleware/pagination.go`）
   - `@handler` 名全服务唯一且语义化
2. 生成：

```bash
./scripts/gen-api.sh merchant-service   # 或 order / user / catalog
```

   - 同模块 handler 合并为单文件（`merge-handlers.py`）
   - **保留**已有 `internal/config/config.go` 与 `cmd/main.go`；只删 goctl 默认 `*api.go` / `etc/*api.yaml`
   - goctl **不会覆盖**已存在的 logic/handler；需按新模板重生成时用 `FORCE_REGEN=1` 或先删对应文件

3. **只改**业务实现：`logic` / `middleware` / `svc/service_context.go`（gen 会还原 svc）
4. **禁止手改** `internal/handler/routes.go`
5. 校验：`./scripts/check-api-routes.sh <svc|all>`

## RPC / Proto

```bash
./scripts/gen-rpc.sh [user|catalog|merchant|all]
# generate-proto.sh 为别名，转调 gen-rpc.sh
```

- pb → [`api/gen`](../api/gen)；共享 client → [`api/rpcclient`](../api/rpcclient)
- server 桩 → `services/<svc>/internal/server/*_service_server.go`
- RPC 业务 → `services/<svc>/internal/rpclogic`（`FORCE_RPC_LOGIC=1` 可覆盖空桩）
- 同一进程仍由 `cmd/main.go` 起 rest + zrpc（`pkg/zrpcx`，etcd 可选）
- **发现**：`MYMALL_ETCD_HOSTS` 非空时注册/发现（key：`user.rpc` / `catalog.rpc` / `merchant.rpc`）；否则 Endpoints 直连

## Model

```bash
./scripts/gen-model.sh [user|catalog|order|merchant|all]
```

- DDL 输入：[`scripts/ddl/<svc>-tables.sql`](../scripts/ddl/)（从 `scripts/*.sql` 抽出；含 goctl 所需 PRIMARY KEY 形态）
- 输出：`services/<svc>/internal/modelgen`
- 后处理：`time.Time`→`common.LocalTime`、字段 initialism（`UserId`→`UserID`）、联合主键表手写补齐
- **repository 仍手写**；user-service 的 `internal/model` 以 type alias 引用 modelgen；其它服务可在 ddl 列齐全后同样切换

## 启动（业务服务）

```go
var c config.Config
conf.MustLoad(configPath, &c)
c.OverlayFromEnv() // MYMALL_*
sqlConn := sqlx.NewMysql(c.MySQL.DSN())
svcCtx := svc.NewServiceContext(&c, sqlConn /* ... */)
server := rest.MustNewServer(c.RestConf, rest.WithCors())
handler.RegisterHandlers(server, svcCtx)
```

`CONFIG_PATH` 默认 `./etc/<svc>-service.yaml`（go-zero RestConf 风格：`Name`/`Host`/`Port`/`Mode`/`Log`…）。

错误契约：[`pkg/xerr`](../pkg/xerr) + [gozero-http-contract.md](gozero-http-contract.md)。

## 分层门禁

```bash
rg 'http\.ResponseWriter|\*http\.Request' services/<svc>/internal/logic
rg '^\s+ctx\s+context\.Context' services/<svc>/internal/logic --glob '*_logic.go'
rg 'internal/httpapi|gorm\.io' services/<svc>
```

## 参考服务

- user / catalog / order / merchant：RestConf + sqlx + goctl HTTP/RPC/model
- inventory-sync：Canal worker，仍用 [`pkg/config`](../pkg/config) + [`pkg/database`](../pkg/database)（GORM）
