# 主题好物集市（固定坑位付费）

## 产品

首页 2×2 固定 4 坑；商家购买某坑展示权，自定义标题/副文案/封面，跳转自选：店铺 / 分类 / 商品。空置显示平台默认图文（通常绑分类）。

## 表

脚本：`scripts/alter-homepage-theme-slots.sql`、`scripts/seed-homepage-theme-menus.sql`

| 表 | 说明 |
|----|------|
| `homepage_theme_slots` | 坑位定义（`slot_key`/`position`/默认创意与跳转） |
| `homepage_theme_packages` | 套餐；`theme_slot_id=0` 通用 |
| `homepage_theme_orders` | 购买/代开通；同坑 `active` 顺延排队 |

钱包流水 `change_type=theme_slot`。

## API（merchant-service）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/home/theme-tiles` | C 端拼装 4 格 |
| GET/PUT | `/api/v1/admin/theme-slots[/:id]` | 坑位 |
| CRUD | `/api/v1/admin/theme-packages` | 套餐 |
| GET | `/api/v1/admin/theme-orders` | 订单 |
| POST | `/api/v1/admin/theme-orders/grant` | 代开通 |
| GET | `/api/v1/merchant/theme-slots` | 可买坑位（含占用结束时间） |
| GET | `/api/v1/merchant/theme-packages` | 上架套餐 |
| POST/GET | `/api/v1/merchant/theme-orders` | 钱包购买 / 本店订单 |

## 前端

- 超管：`/admin/themes`（`ThemeSlots.vue`），perms `marketing:theme:*`
- 商家：`/merchant/themes`（`ThemePromote.vue`），`theme:list` / `theme:buy`
- mall-uni：首页 `listThemeTiles()`，付费角标「推广」

## 与首页展位关系

独立表与常量，支付/顺延规则对齐 [`homepage-slots.md`](homepage-slots.md)，不混用 `slot_type`。
