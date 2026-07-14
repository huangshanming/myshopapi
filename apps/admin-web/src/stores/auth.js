import { defineStore } from 'pinia'
import http from '../api/http'

function parseJwt(token) {
  try {
    const payload = token.split('.')[1]
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
  } catch {
    return {}
  }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('mymall_token') || '',
    user: JSON.parse(localStorage.getItem('mymall_user') || 'null'),
    role: localStorage.getItem('mymall_role') || '',
    shopId: Number(localStorage.getItem('mymall_shop_id') || 0),
    userId: Number(localStorage.getItem('mymall_user_id') || 0),
  }),
  getters: {
    isAdmin: (s) => s.role === 'platform_admin' || s.role === 'admin',
    isMerchant: (s) => s.role === 'merchant_owner' || s.role === 'merchant_staff',
  },
  actions: {
    async login(mobile, password, shopId = 0) {
      const body = await http.post('/api/v1/user/login', {
        mobile,
        password,
        shop_id: shopId || undefined,
      })
      const token = body.data.token
      const user = body.data.user
      const claims = parseJwt(token)
      this.token = token
      this.user = user
      this.role = claims.role || user.role || ''
      this.shopId = Number(claims.shop_id || shopId || 0)
      this.userId = Number(claims.user_id || user.id || 0)
      localStorage.setItem('mymall_token', token)
      localStorage.setItem('mymall_user', JSON.stringify(user))
      localStorage.setItem('mymall_role', this.role)
      localStorage.setItem('mymall_shop_id', String(this.shopId))
      localStorage.setItem('mymall_user_id', String(this.userId))
      return this.role
    },
    logout() {
      this.token = ''
      this.user = null
      this.role = ''
      this.shopId = 0
      this.userId = 0
      localStorage.clear()
    },
  },
})
