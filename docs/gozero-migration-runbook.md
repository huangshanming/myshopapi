# goctl HTTP 全套工作流

## 工具

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5
export PATH="$(go env GOPATH)/bin:$PATH"
```

本仓库自定义模板：[`deploy/goctl-template`](../deploy/goctl-template)（handler 调用 `logic.Xxx(w,r)`，便于对接现有 JSON 契约，无需一次性补全所有 `.api` types）。

## 唯一正确流程

1. 编辑服务 `api/*.api`（路由 / `@server` 中间件 / `@handler`）
2. 生成：

```bash
./scripts/gen-api.sh merchant-service
```

3. **只改**业务实现：
   - `internal/logic/*_logic.go`
   - `internal/middleware/*_middleware.go`（接 `pkg/middleware`）
   - `internal/svc/service_context.go`（`gen-api.sh` 会保留该文件不被覆盖）
4. **禁止手改** `internal/handler/routes.go`（goctl 生成，带 `DO NOT EDIT`）
5. 校验：

```bash
./scripts/check-api-routes.sh merchant-service
./scripts/check-api-routes.sh all
```

## 启动方式（不变）

保留 `CONFIG_PATH` + [`pkg/httpserver.NewRest`](../pkg/httpserver/server.go) + [`pkg/xerr`](../pkg/xerr)。  
`goctl` 生成的根 `*api.go` / `internal/config` / `etc/*api.yaml` 会被 `gen-api.sh` 删除，不以它们替换 `cmd/main.go`。

## main

```go
xerr.RegisterErrorHandler()
// ... NewServiceContext ...
handler.RegisterHandlers(server, svcCtx)
```

## 契约

见 [gozero-http-contract.md](gozero-http-contract.md)。

## 参考

- 样板：[`services/merchant-service`](../services/merchant-service)
- 小型参考：[`services/inventory-sync-service`](../services/inventory-sync-service)
