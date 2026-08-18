import http from 'node:http'
import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'

const svc = {
  user: 'http://127.0.0.1:8881',
  catalog: 'http://127.0.0.1:8882',
  order: 'http://127.0.0.1:8883',
  merchant: 'http://127.0.0.1:8884',
  agent: 'http://127.0.0.1:8886',
  lottery: 'http://127.0.0.1:8887',
  recommend: 'http://127.0.0.1:8888',
}

function to(target) {
  return { target, changeOrigin: true }
}

/** Vite 的 proxy.router 对带 ID 的子路径不可靠，评价接口需打到 order-service */
function proxyProductReviewsPlugin() {
  return {
    name: 'proxy-product-reviews-to-order',
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        const url = req.url || ''
        if (!/^\/api\/v1\/products\/\d+\/reviews(?:\?|$)/.test(url)) {
          next()
          return
        }
        const target = new URL(svc.order)
        const headers = { ...req.headers, host: target.host }
        const opts = {
          protocol: target.protocol,
          hostname: target.hostname,
          port: target.port,
          path: url,
          method: req.method,
          headers,
        }
        const p = http.request(opts, (pr) => {
          res.writeHead(pr.statusCode || 502, pr.headers)
          pr.pipe(res)
        })
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
  plugins: [uni(), proxyProductReviewsPlugin()],
  server: {
    port: 5175,
    proxy: {
      '/api/v1/user/favorites': to(svc.catalog),
      '/api/v1/user/article-favorites': to(svc.catalog),
      '/api/v1/user/article-likes': to(svc.catalog),
      '/api/v1/user/articles': to(svc.catalog),
      '/api/v1/user/article-uploads': to(svc.catalog),
      '/api/v1/user/review-uploads': to(svc.order),
      '/api/v1/products': to(svc.catalog),
      '/api/v1/product_category': to(svc.catalog),
      '/api/v1/banners': to(svc.catalog),
      '/api/v1/articles': to(svc.catalog),
      '/api/v1/comment-emojis': to(svc.catalog),
      '/api/v1/shops': to(svc.merchant),
      '/api/v1/home': to(svc.merchant),
      '/api/v1/map': to(svc.merchant),
      '/api/v1/seckill': to(svc.merchant),
      '/api/v1/coupons': to(svc.merchant),
      '/api/v1/user/coupons': to(svc.merchant),
      '/api/v1/user/points-mall': to(svc.user),
      '/api/v1/user/cps': to(svc.user),
      '/api/v1/lottery': to(svc.lottery),
      '/api/v1/recommend': to(svc.recommend),
      '/api/v1/orders': to(svc.order),
      '/api/v1/agents': to(svc.agent),
      '/api/v1/user': to(svc.user),
      '/api/v1/regions': to(svc.user),
      '/uploads/points-mall': to(svc.user),
      '/uploads': to(svc.catalog),
    },
  },
})
