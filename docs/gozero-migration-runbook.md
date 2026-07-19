# goctl HTTP 迁移 Runbook

## 工具

```bash
go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5
export PATH="$(go env GOPATH)/bin:$PATH"
```

## 生成

在服务目录准备 `api/*.api` 后：

```bash
OUT=./_goctl_out
rm -rf "$OUT"
goctl api go -api api/<name>.api -dir "$OUT" -style go_zero
```

注意：`goctl` 的 `-dir` 请用**仓库内相对/绝对路径**，避免 `/tmp` 被错误解析。

将生成的 `internal/handler|logic|types|svc` 合并进服务，**保留**现有 `pkg/config` + `CONFIG_PATH` 启动方式，不要直接替换为 goctl 默认 `etc/*.yaml` RestConf（除非单独评估）。

## Handler 模板

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
handler.RegisterHandlers(server, svcCtx, middleware.RequestID())
```

## 契约

见 [gozero-http-contract.md](gozero-http-contract.md)。

## 参考实现

[`services/inventory-sync-service`](../services/inventory-sync-service)：healthz/readyz 已按 goctl 形态 + `httpx.OkJsonCtx`。
