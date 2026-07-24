import http from 'node:http'
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

function proxyByPath(match, target) {
  return {
    name: `proxy-${match.source}`,
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        const url = req.url || ''
        if (!match.test(url.split('?')[0])) {
          next()
          return
        }
        const t = new URL(target)
        const p = http.request(
          {
            protocol: t.protocol,
            hostname: t.hostname,
            port: t.port,
            path: url,
            method: req.method,
            headers: { ...req.headers, host: t.host },
          },
          (pr) => {
            res.writeHead(pr.statusCode || 502, pr.headers)
            pr.pipe(res)
          },
        )
        p.on('error', (err) => {
          res.statusCode = 502
          res.end(String(err.message || err))
        })
        req.pipe(p)
      })
    },
  }
}

export default defineConfig({
  plugins: [
    vue(),
    proxyByPath(/^\/api\/v1\/products\/\d+\/reviews$/, svc.order),
    proxyByPath(/^\/api\/v1\/admin\/users\/\d+\/favorites$/, svc.catalog),
  ],
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
      '/uploads/points-mall': to(svc.user),
      '/uploads': to(svc.catalog),
      '/api/v1/admin/products': to(svc.catalog),
      '/api/v1/admin/categories': to(svc.catalog),
      '/api/v1/admin/articles': to(svc.catalog),
      '/api/v1/admin/article-categories': to(svc.catalog),
      '/api/v1/admin/article-comments': to(svc.catalog),
      '/api/v1/admin/comment-emojis': to(svc.catalog),
      '/api/v1/admin/article-uploads': to(svc.catalog),
      '/api/v1/admin/banners': to(svc.catalog),
      '/api/v1/admin/shop-uploads': to(svc.catalog),
      '/api/v1/products': to(svc.catalog),
      '/api/v1/product_category': to(svc.catalog),
      '/api/v1/banners': to(svc.catalog),
      '/api/v1/articles': to(svc.catalog),
      // order
      '/api/v1/merchant/orders': to(svc.order),
      '/api/v1/merchant/after-sales': to(svc.order),
      '/api/v1/merchant/reviews': to(svc.order),
      '/api/v1/admin/orders': to(svc.order),
      '/api/v1/admin/after-sales': to(svc.order),
      '/api/v1/admin/reviews': to(svc.order),
      '/api/v1/admin/logistics': to(svc.order),
      '/api/v1/logistics': to(svc.order),
      '/api/v1/orders': to(svc.order),
      // merchant（店铺 / 入驻；勿盖住上面的 catalog merchant 路径）
      '/api/v1/admin/shops': to(svc.merchant),
      '/api/v1/admin/applications': to(svc.merchant),
      '/api/v1/admin/seckill': to(svc.merchant),
      '/api/v1/admin/homepage-packages': to(svc.merchant),
      '/api/v1/admin/homepage-settings': to(svc.merchant),
      '/api/v1/admin/homepage-orders': to(svc.merchant),
      '/api/v1/admin/theme-slots': to(svc.merchant),
      '/api/v1/admin/theme-packages': to(svc.merchant),
      '/api/v1/admin/theme-orders': to(svc.merchant),
      '/api/v1/admin/coupons': to(svc.merchant),
      '/api/v1/admin/notifications': to(svc.user),
      '/api/v1/admin/tasks': to(svc.user),
      '/api/v1/admin/points-products': to(svc.user),
      '/api/v1/admin/points-orders': to(svc.user),
      '/api/v1/user/points-mall': to(svc.user),
      '/uploads/points-mall': to(svc.user),
      '/api/v1/shops': to(svc.merchant),
      '/api/v1/home': to(svc.merchant),
      '/api/v1/map': to(svc.merchant),
      '/api/v1/seckill': to(svc.merchant),
      '/api/v1/coupons': to(svc.merchant),
      '/api/v1/user/coupons': to(svc.merchant),
      '/api/v1/merchant/wallet': to(svc.merchant),
      '/api/v1/merchant/seckill': to(svc.merchant),
      '/api/v1/merchant/homepage-packages': to(svc.merchant),
      '/api/v1/merchant/homepage-orders': to(svc.merchant),
      '/api/v1/merchant/theme-slots': to(svc.merchant),
      '/api/v1/merchant/theme-packages': to(svc.merchant),
      '/api/v1/merchant/theme-orders': to(svc.merchant),
      '/api/v1/merchant/coupons': to(svc.merchant),
      '/api/v1/merchant': to(svc.merchant),
      '/api/v1/articles': to(svc.catalog),
      // user + 其余 /api
      '/api': to(svc.user),
    },
  },
})
