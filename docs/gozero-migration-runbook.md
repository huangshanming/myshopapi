# goctl HTTP 迁移 Runbook

## 工具

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5
export PATH="$(go env GOPATH)/bin:$PATH"
```

## 路由 SSOT（现行规范）

业务 HTTP 路由以各服务 `api/*.api` 为唯一来源，运行时由 `internal/handler/routes.go` 的 `RegisterHandlers` 注册。

| 文件 | 职责 |
|------|------|
| `api/<name>.api` | 声明 path / method / `@server` 中间件分组 / `@handler` |
| `internal/handler/routes.go` | `AddRoutes` 唯一落点，绑定现有按端 handler |
| `cmd/main.go` | 启动、依赖注入、组装 middleware Bundle、调用 `RegisterHandlers`；**禁止**再写 `Path: "/api/..."` |

改路由流程：

1. 编辑 `api/*.api`
2. 同步修改 `internal/handler/routes.go`（中间件组与路径一致）
3. 校验：

```bash
./scripts/check-api-routes.sh merchant-service   # 或 order-service / user-service / catalog-service
./scripts/check-api-routes.sh all
```

参考实现：[`services/merchant-service`](../services/merchant-service)（`api/merchant.api` + `internal/handler/routes.go` + `internal/middleware`）。

## 可选：goctl 对照生成

在服务目录准备好 `.api` 后，可生成到临时目录对照，**不要**直接覆盖现有 `handler/{admin,user,...}`：

```bash
OUT=./_goctl_out
rm -rf "$OUT"
goctl api go -api api/<name>.api -dir "$OUT" -style go_zero
```

注意：`goctl` 的 `-dir` 请用**仓库内相对/绝对路径**，避免 `/tmp` 被错误解析。

将生成物中有价值的 types/注释合并进服务时，**保留**现有 `pkg/config` + `CONFIG_PATH` 启动方式，不要直接替换为 goctl 默认 `etc/*.yaml` RestConf（除非单独评估）。

## Handler 模板（后续可选闭包形态）

```go
func XxxHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var req types.XxxReq
    if err := httpx.Parse(r, &req); err != nil {
      httpx.ErrorCtx(r.Context(), w, xerr.BadRequest(err.Error()))
      return
    }
    l := logic.NewXxxLogic(r.Context(), svcCtx)
    resp, err := l.Xxx(&req)
    if err != nil {
      httpx.ErrorCtx(r.Context(), w, err)
      return
    }
    httpx.OkJsonCtx(r.Context(), w, resp)
  }
}
```

## main 必备

```go
xerr.RegisterErrorHandler()
// ... NewServiceContext ...
handler.RegisterHandlers(server, svcCtx, healthReg, svcMW.NewBundle())
```

## 契约

见 [gozero-http-contract.md](gozero-http-contract.md)。

## 参考实现

- [`services/merchant-service`](../services/merchant-service)：路由 SSOT + `RegisterHandlers`
- [`services/inventory-sync-service`](../services/inventory-sync-service)：小型 goctl 形态 + `httpx.OkJsonCtx`
