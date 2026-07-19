# 优惠券系统

平台券 + 商家券完整 MVP。定义/发放在 **merchant-service**，下单锁券与核销编排在 **order-service**。

## 表

脚本：`scripts/alter-coupons.sql`、`scripts/alter-orders-coupon.sql`、`scripts/seed-coupon-menus.sql`

| 表 | 说明 |
|----|------|
| `coupons` | 券模板（类型/门槛/库存/有效期/渠道/身份） |
| `coupon_scopes` | 品类/商品适用范围 |
| `user_coupons` | 用户券实例：`unused`/`locked`/`used`/`expired` |
| `coupon_grants` | 定向发放批次 |
| `coupon_redeem_logs` | 核销/解锁/返还流水 |
| `orders` | `goods_amount` / `discount_amount` / `pay_amount` / `user_coupon_id` |

类型：`full_reduce` / `no_threshold` / `category` / `product` / `discount`。

## 核心规则

- 一单仅一张券；领取扣库存；限领按人计数
- 下单锁券 → 取消/`failed` 解锁 → 确认收货（或商家完成）核销
- 全额退款返还券（未过期可再用，过期标失效）；部分退不返还
- 不可叠加：`stackable=0` 或订单含秒杀行则不可用券
- 实付保底 ¥0.01；钱包冻结金额为 `pay_amount`

## API 摘要

| 域 | 路径 |
|----|------|
| 超管 | `/api/v1/admin/coupons` CRUD、`/off`、`/copy`、`/grant`、claims/redeems/stats |
| 商家 | `/api/v1/merchant/coupons` 同构（本店） |
| C 端 | `/api/v1/coupons/center`、`/popup`、`/:id/claim`、`/api/v1/user/coupons` |
| 结算 | `POST /api/v1/orders/coupon-preview`；`POST /api/v1/orders` 带 `user_coupon_id` |
| 内部 | `/api/v1/internal/coupons/{match,lock,unlock,redeem,return,order-gift}` |

## 前端

- 超管：`/admin/coupons`；商家：`/merchant/coupons`
- mall-uni：领券中心、我的优惠券、确认订单选券、首页通栏

## 权限

- 超管：`marketing:coupon:list|edit|grant`（菜单 id 121–123）
- 商家：`coupon:list|edit|grant`（shop_menu 34–36）
