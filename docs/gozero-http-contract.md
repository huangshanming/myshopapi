# go-zero HTTP 契约（新）

成功与失败均对齐 go-zero `httpx.OkJsonCtx` / `httpx.ErrorCtx`。

## 成功

- **HTTP**：`200`
- **Body**：业务 DTO 本身（无 `{code,msg,data}` 外壳）

```json
{ "list": [], "total": 0 }
```

空成功可用 HTTP 200 无 body，或 `{}`。

## 失败

- **HTTP**：`4xx` / `5xx`（由业务码映射）
- **Body**：

```json
{ "code": 400, "msg": "参数错误" }
```

| code | HTTP | 含义 |
|-----:|-----:|------|
| 400 | 400 | 参数/业务校验 |
| 401 | 401 | 未登录 |
| 403 | 403 | 无权限 |
| 404 | 404 | 不存在 |
| 500 | 500 | 未分类系统错误（对前端文案泛化） |

实现：[`pkg/xerr`](../pkg/xerr)、各服务 `main` 调用 `xerr.RegisterErrorHandler()`。

## 与旧契约对照

| | 旧 | 新 |
|--|----|----|
| 成功 | HTTP 200 + `{code:200,msg,data}` | HTTP 200 + DTO |
| 业务失败 | HTTP 200 + `{code:4xx,msg}` | HTTP 4xx + `{code,msg}` |
| 鉴权失败 | HTTP 401/403 + `{code,msg,data}` | HTTP 401/403 + `{code,msg}` |

## 前端

迁移期拦截器双读：若 body 含数值 `code` 且含 `data` 字段，则按旧信封 unwrap；否则整包当 DTO。稳定后删除双读。

## goctl

```bash
# 安装（与 go.mod 对齐）
go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5

# 在服务目录生成（示例）
cd services/inventory-sync-service
goctl api go -api api/inventory.api -dir . -style go_zero
```

详见 [`docs/gozero-migration-runbook.md`](gozero-migration-runbook.md)。
