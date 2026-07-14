# mymall 管理端 / 商家端

Vue3 + Vite + Element Plus。开发端口 `5174`。

## 启动

```bash
cd apps/admin-web
npm install
npm run dev
```

默认通过 Vite 代理直连本地微服务（8881–8884）。若走 APISIX：

```bash
VITE_API_BASE=http://localhost:9080 npm run dev
```

## 种子账号

| 角色 | 手机号 | 密码 |
|------|--------|------|
| 平台管理员 | 13900000001 | 123456 |
| 商家 | 13900000002 | 123456 |
