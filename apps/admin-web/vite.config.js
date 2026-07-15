import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const proxyTarget = env.VITE_API_PROXY || 'http://localhost:8881'

  return {
    plugins: [vue()],
    server: {
      port: 5174,
      proxy: {
        '/api': {
          target: proxyTarget,
          changeOrigin: true,
          configure: (proxy) => {
            // 开发时按路径分流到不同微服务（无需 APISIX）
            proxy.on('proxyReq', (proxyReq, req) => {
              const url = req.url || ''
              let target = 'http://localhost:8881'
              if (url.includes('/merchant/products') || url.includes('/admin/products') || url.includes('/admin/categories')) {
                target = 'http://localhost:8882'
              } else if (url.includes('/merchant/orders') || url.includes('/admin/orders') || url.match(/\/api\/v1\/orders/)) {
                target = 'http://localhost:8883'
              } else if (
                url.includes('/merchant/') ||
                url.includes('/admin/applications') ||
                (url.includes('/admin/shops') && !url.includes('/admin/system'))
              ) {
                target = 'http://localhost:8884'
              }
              // /admin/auth|/menus|/roles|/users|/admins|/configs → user-service :8881
              proxyReq.setHeader('host', new URL(target).host)
            })
          },
          router: (req) => {
            const url = req.url || ''
            if (url.includes('/merchant/products') || url.includes('/admin/products') || url.includes('/admin/categories')) {
              return 'http://localhost:8882'
            }
            if (url.includes('/merchant/orders') || url.includes('/admin/orders') || /^\/api\/v1\/orders/.test(url)) {
              return 'http://localhost:8883'
            }
            if (
              url.includes('/merchant/') ||
              url.includes('/admin/applications') ||
              (url.includes('/admin/shops') && !url.includes('/admin/system'))
            ) {
              return 'http://localhost:8884'
            }
            return 'http://localhost:8881'
          },
        },
      },
    },
  }
})
