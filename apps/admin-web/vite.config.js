import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 本地微服务端口（VITE_API_BASE 为空时生效；更具体的前缀写在前面）
const svc = {
  user: 'http://127.0.0.1:8881',
  catalog: 'http://127.0.0.1:8882',
  order: 'http://127.0.0.1:8883',
  merchant: 'http://127.0.0.1:8884',
}

function to(target) {
  return { target, changeOrigin: true }
}

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5174,
    proxy: {
      // catalog 商家商品中台（需在通用 /api/v1/merchant 之前）
      '/api/v1/merchant/products': to(svc.catalog),
      '/api/v1/merchant/skus': to(svc.catalog),
      '/api/v1/merchant/stocks': to(svc.catalog),
      '/api/v1/merchant/uploads': to(svc.catalog),
      '/api/v1/merchant/tags': to(svc.catalog),
      '/api/v1/merchant/attr-templates': to(svc.catalog),
      '/api/v1/merchant/schedules': to(svc.catalog),
      '/api/v1/merchant/auth': to(svc.catalog),
      '/api/v1/merchant/shop': to(svc.catalog),
      '/api/v1/merchant/articles': to(svc.catalog),
      '/api/v1/merchant/article-categories': to(svc.catalog),
      '/api/v1/merchant/article-comments': to(svc.catalog),
      '/api/v1/merchant/article-uploads': to(svc.catalog),
      '/api/v1/merchant/notifications': to(svc.catalog),
      '/uploads': to(svc.catalog),
      '/api/v1/admin/products': to(svc.catalog),
      '/api/v1/admin/categories': to(svc.catalog),
      '/api/v1/admin/articles': to(svc.catalog),
      '/api/v1/admin/article-categories': to(svc.catalog),
      '/api/v1/admin/article-comments': to(svc.catalog),
      '/api/v1/admin/article-uploads': to(svc.catalog),
      '/api/v1/products': to(svc.catalog),
      '/api/v1/product_category': to(svc.catalog),
      // order
      '/api/v1/merchant/orders': to(svc.order),
      '/api/v1/admin/orders': to(svc.order),
      '/api/v1/orders': to(svc.order),
      // merchant（店铺 / 入驻；勿盖住上面的 catalog merchant 路径）
      '/api/v1/admin/shops': to(svc.merchant),
      '/api/v1/admin/applications': to(svc.merchant),
      '/api/v1/merchant': to(svc.merchant),
      // user + 其余 /api
      '/api': to(svc.user),
    },
  },
})
