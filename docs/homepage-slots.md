# 首页付费展位 + 种草互动

## 表结构

脚本：`scripts/alter-homepage-slots.sql`、`scripts/alter-community-article-engagement.sql`、`scripts/seed-homepage-menus.sql`

| 表 | 说明 |
|----|------|
| `homepage_slot_packages` | 套餐：`slot_type` + 价格 + `duration_days` |
| `homepage_slot_settings` | 各类型首页条带 `home_limit` |
| `homepage_slot_orders` | 购买/代开通记录，`active` 至 `end_at` |
| `community_article` | 增加 `audience_count` / `read_count` / `collect_count` |
| `article_likes` / `article_favorites` / `article_audiences` | 互动与去重观众 |

`slot_type`：`brand_shop` | `quality_shop` | `article`

## 规则

- 支付：店铺钱包自购（`pay_source=wallet`）或超管代开通（`admin`，不扣款）
- 续费顺延：同店同类型（文章同 `target_id`）若已有 `active`，新单 `start_at = max(now, 旧 end_at)`
- 首页条带数量 = `homepage_slot_settings.home_limit`，不是套餐字段
- 列表排序：生效付费优先，其后有机内容；读时可将过期单标 `expired`
- 种草：`audience_count` 同一用户进详情计 1；`read_count` 每次打开 +1；另有点赞/收藏

## merchant-service API

| 方法 | 路径 | 说明 |
|------|------|------|
| CRUD | `/api/v1/admin/homepage-packages` | 超管套餐 |
| GET/PUT | `/api/v1/admin/homepage-settings` | 条带数量 |
| GET | `/api/v1/admin/homepage-orders` | 订单列表 |
| POST | `/api/v1/admin/homepage-orders/grant` | 代开通 |
| GET | `/api/v1/merchant/homepage-packages` | 上架套餐 |
| POST | `/api/v1/merchant/homepage-orders` | 钱包购买 |
| GET | `/api/v1/merchant/homepage-orders` | 本店记录 |
| GET | `/api/v1/shops/home-slots?slot_type=` | 首页条带（店铺两类） |
| GET | `/api/v1/shops/list?slot_type=` | 全部列表付费优先分页 |

钱包流水 `change_type=homepage_slot`。

相关：**主题好物集市**（固定 4 坑付费）见 [`homepage-theme-slots.md`](homepage-theme-slots.md)。

## catalog-service API（种草 C 端）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/articles/list` | 已发布；`home=1` 用 `home_limit` |
| GET | `/api/v1/articles/:id` | 详情；可选 `X-User-Id` 记观众/阅读 |
| POST/DELETE | `/api/v1/articles/:id/like` | 点赞（JWT） |
| POST/DELETE | `/api/v1/articles/:id/favorite` | 收藏（JWT） |
| GET | `/api/v1/articles/:id/engagement` | 状态（JWT） |

文章展位排序：同库 `LEFT JOIN homepage_slot_orders`（`slot_type=article` 且 `active`）。

## 响应契约

成功：HTTP 200 + 业务 DTO；失败：非 2xx + `{code,msg}`。详见 [`docs/gozero-http-contract.md`](gozero-http-contract.md)。

## 权限与菜单

- 超管：`business:homepage:list|package|grant`（`seed-homepage-menus.sql` id 112–114）
- 商家：`homepage:list|buy`（catalog `EnsureShopMenus` id 30–31，「首页推广」）

## 前端

- admin-web：`/admin/homepage`、`/merchant/homepage`；Vite 将 homepage-* 代理到 merchant `:8884`
- mall-uni：首页三块接 `home-slots` / `articles?home=1`；`pages/shop/*`、`pages/community/*`；Vite `/api/v1/articles` → catalog

## 代理注意

- 本地 Vite 更具体前缀写在通用 `/api/v1/merchant` 之前
- APISIX：公开文章 GET → catalog；点赞收藏 → protected + JWT；homepage 管理 → merchant-service
