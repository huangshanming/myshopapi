# goctl HTTP 全套工作流

## 工具

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5
export PATH="$(go env GOPATH)/bin:$PATH"
```

本仓库自定义模板：[`deploy/goctl-template`](../deploy/goctl-template)（handler 调用 `logic.Xxx(w,r)`，便于对接现有 JSON 契约，无需一次性补全所有 `.api` types）。

## 唯一正确流程

1. 编辑服务 `api/*.api`（路由 / `@server` 中间件 / `group` / `@handler`）
   - `group` 写成 `<端>/<模块>`（可多层），goctl 生成到 `internal/handler/<端>/<模块>/`、`internal/logic/<端>/<模块>/`（叶子目录为 Go package 名）
   - 同一中间件链下按模块拆多个 `@server`（只改 `group`）
   - `@handler` 名全服务唯一且语义化（如 `MerchantCreateProduct`、`AdminGetArticle`），**禁止** `Create2` / `Detail3` 这种数字后缀
   - catalog 范本模块：`public/{health,product,category,banner,article}`、`user/{article,favorite}`、`merchant/{product,shopops,article,notification}`、`admin/{article,category,product,banner,comment,user_favorite,shop}`
2. 生成：

```bash
./scripts/gen-api.sh catalog-service
```

   `gen-api.sh` 会在 goctl 之后把同一模块下「一接口一文件」合并为单个文件（如 `admin/article/` → `admin_article_handler.go`）。**不要**手拆回多文件；下次 gen 会再合并。

3. **只改**业务实现：
   - `internal/logic/<端>/<模块>/*_logic.go`（可委托 `internal/httpapi/<端>/` 或域内 httpapi）
   - `internal/middleware/*_middleware.go`（接 `pkg/middleware`）
   - `internal/svc/service_context.go`（`gen-api.sh` 会保留该文件不被覆盖）
4. **禁止手改** `internal/handler/routes.go`（goctl 生成，带 `DO NOT EDIT`）
5. 校验：

```bash
./scripts/check-api-routes.sh catalog-service
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

- 嵌套 group 样板：[`services/catalog-service`](../services/catalog-service)（`group: admin/article` 等）
- 其它服务：[`services/merchant-service`](../services/merchant-service)
- 小型参考：[`services/inventory-sync-service`](../services/inventory-sync-service)
