# mymall 管理端 / 商家端

Vue3 + Vite + Element Plus。开发端口 `5174`。

## 启动

```bash
cd apps/admin-web
cp .env.example .env   # 首次可复制；仓库已带默认 .env
npm install
npm run dev
```

API 地址在 `.env` 里配置 `VITE_API_BASE`（改完需重启 `npm run dev`）：

| 场景 | `VITE_API_BASE` |
|------|-----------------|
| 默认：Vite 代理直连 8881–8884 | 留空 |
| 经 APISIX | `http://localhost:9080` |

本地覆盖可建 `.env.local`（已 gitignore），不会提交。

RBAC 库表（首次）：

```bash
mysql -u homestead -p mymall < scripts/init-rbac-tables.sql
mysql -u homestead -p mymall < scripts/seed-rbac.sql
```

## 种子账号

| 角色 | 手机号 | 密码 |
|------|--------|------|
| 平台管理员 | 13900000001 | 123456 |
| 商家 | 13900000002 | 123456 |
