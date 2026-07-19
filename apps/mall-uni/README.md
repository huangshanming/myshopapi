# mymall 用户端（UniApp）

Vue3 + Vite + uni-app，H5 优先，可编译微信小程序。

## 启动

```bash
# 需先启动 user(:8881) / catalog(:8882) / order(:8883) / merchant(:8884)
cd apps/mall-uni
yarn install
yarn dev:h5
```

浏览器打开控制台提示的本地地址（默认 `http://localhost:5175`）。

首页商户种子（可选）：

```bash
bash scripts/seed-home-shop-images.sh
mysql -u homestead -p mymall < scripts/seed-home-shops.sql
```

## 功能（MVP）

- 首页：对齐 `test/index.html` 结构；推荐商品接 `GET /api/v1/products/list`；商户接 `GET /api/v1/shops/list`
- 商品详情 + 立即购买下单
- 登录 / 注册、我的订单、取消订单
- 分类 / 购物车 Tab 占位

## 代理

开发态 Vite 将 `/api/v1/products`、`/api/v1/product_category` → catalog，`/api/v1/shops` → merchant，`/api/v1/orders` → order，`/api/v1/user` → user，`/uploads` → catalog。
