# C 端站内信（消息中心）

仅站内信，无微信/短信；C 端轮询未读。与商家 `shop_notifications` 分表。实现在 **user-service**；订单状态变更由 **order-service** 调内部接口写信。

## 表

脚本：`scripts/alter-user-notifications.sql`、`scripts/seed-user-message-menus.sql`

| 表 | 说明 |
|----|------|
| `user_notifications` | 用户消息（每人一条；全局公告展开写入） |
| `user_notification_batches` | 超管发送批次记录 |

`msg_type`：`system` / `order` / `announce`；`link_type`：`none` / `order`（`link_id` 为订单 ID）。

## 规则

- 全局公告按用户展开；单次最多约 5000 人
- 订单消息：`msg_type=order`，`link_type=order`，`link_id=order.id`
- 触发点：支付成功、取消、发货、完成/确认收货、售后退款完成

## API 摘要

| 域 | 路径 |
|----|------|
| C 端 | `GET /api/v1/user/notifications`、`/unread-count`；`POST .../:id/read`、`/read-all` |
| 超管 | `POST /api/v1/admin/notifications/send`；`GET .../sends`；`GET .../sends/:id/recipients` |
| 内部 | `POST /api/v1/internal/notifications`（order-service 调用） |

## 前端

- mall-uni：首页铃铛（未读角标 + 抖动）→ `/pages/message/list`；订单消息可跳转订单详情
- 超管：`/admin/messages`（权限 `business:message:send`，菜单 id 124）

## 网关 / 本地代理

- APISIX：C 端 notifications → user-service（protected）；超管 send/sends → user-service（user-admin-rbac）
- Vite：mall-uni `/api/v1/user` → user；admin-web `/api/v1/admin/notifications` → user
