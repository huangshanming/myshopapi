import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'

const svc = {
  user: 'http://127.0.0.1:8881',
  catalog: 'http://127.0.0.1:8882',
  order: 'http://127.0.0.1:8883',
}

function to(target) {
  return { target, changeOrigin: true }
}

export default defineConfig({
  plugins: [uni()],
  server: {
    port: 5175,
    proxy: {
      '/api/v1/products': to(svc.catalog),
      '/api/v1/product_category': to(svc.catalog),
      '/api/v1/orders': to(svc.order),
      '/api/v1/user': to(svc.user),
      '/uploads': to(svc.catalog),
    },
  },
})
