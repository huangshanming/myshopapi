# 商品收藏 + 订单评价

## 表结构

- `product_favorites`：用户收藏（`UNIQUE(user_id, product_id)`）
- `products.collect_count` / `avg_rating` / `review_count` / `good_rate`
- `product_reviews`：一单一评（`UNIQUE(order_id)`），`status=visible|deleted`
- `product_review_images`
- `orders.status` 增加 `reviewed`，可选 `reviewed_at`

脚本：`scripts/alter-product-favorites.sql`、`scripts/alter-product-reviews.sql`、`scripts/seed-review-menus.sql`

## 状态机（评价相关）

```
shipped --[C端确认收货]--> completed --[提交评价]--> reviewed
```

仅 `completed` 可「去评价」；`reviewed` 仅「查看评价」。用户不可改评/删评。

## 收藏 API（catalog-service）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/user/favorites` | `{product_id}` 幂等收藏 |
| DELETE | `/api/v1/user/favorites/:product_id` | 幂等取消 |
| POST | `/api/v1/user/favorites/batch-remove` | `{product_ids}` |
| GET | `/api/v1/user/favorites` | 分页列表，`invalid` 非在售 |
| GET | `/api/v1/products/:id/favorite` | `{favorited}` |
| GET | `/api/v1/products/:id/favorite-count` | `{count}` |
| GET | `/api/v1/admin/users/:id/favorites` | 总后台查看 |

## 确认收货 + 评价 API（order-service）

| 方法 | 路径 | 说明 |
|------|------|------|
| PUT | `/api/v1/orders/:id/confirm-receive` | shipped→completed |
| GET | `/api/v1/orders/:id/review-eligible` | 是否可评/已评 |
| POST | `/api/v1/orders/:id/reviews` | 提交评价 |
| GET | `/api/v1/orders/:id/review` | 本单评价 |
| GET | `/api/v1/products/:id/reviews` | 公开评价列表 |
| POST | `/api/v1/user/review-uploads` | 评价图 |
| GET/PUT/DELETE | `/api/v1/merchant/reviews...` | 本店列表/回复/软删 |
| GET/DELETE | `/api/v1/admin/reviews...` | 全站列表/软删 |

提交 body：`rating`(1-5)、`content`、`is_anonymous`、`order_item_id`、`images`(≤9)。

好评/中评/差评：`rating≥4` / `=3` / `≤2`。

## 防重

- 收藏：唯一键 + `OnConflict DoNothing`；取消按行删除后 `collect_count` 防负
- 评价：`uk_order` + 提交时校验订单 `completed` 且无评价记录；成功后写 `reviewed`

## 页面

- mall-uni：详情收藏与评价区、我的收藏、订单确认收货/去评价/查看评价
- 商家后台：评价管理（回复/删除）；商品列表均分/评价数/好评率/收藏
- 总后台：评价管理、商品收藏人数排序、用户收藏弹窗

## RBAC

- 商家：`product:review:list|reply|delete`（EnsureShopMenus id 27–29）
- 平台：`business:review:list|delete`（seed-review-menus）
