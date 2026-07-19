import { defineStore } from 'pinia'
import http from '../api/http'
import { fetchAuthMe } from '../api/system'
import { fetchMerchantMe } from '../api/merchant-product'

function parseJwt(token) {
  try {
    let payload = token.split('.')[1] || ''
    payload = payload.replace(/-/g, '+').replace(/_/g, '/')
    const pad = payload.length % 4
    if (pad) payload += '='.repeat(4 - pad)
    return JSON.parse(atob(payload))
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
    roles: JSON.parse(localStorage.getItem('mymall_roles') || '[]'),
    perms: JSON.parse(localStorage.getItem('mymall_perms') || '[]'),
    menuTree: JSON.parse(localStorage.getItem('mymall_menus') || '[]'),
    isShopOwner: localStorage.getItem('mymall_shop_owner') === '1',
  }),
  getters: {
    isAdmin: (s) => s.role === 'platform_admin' || s.role === 'admin',
    isMerchant: (s) => s.role === 'merchant_owner' || s.role === 'merchant_staff',
    isSuperAdmin: (s) => Array.isArray(s.roles) && s.roles.includes('super_admin'),
  },
  actions: {
    hasPerm(code) {
      if (!code) return true
      if (this.isSuperAdmin) return true
      if (this.isMerchant && this.isShopOwner) return true
      // 商家权限未加载时放行，避免首屏误拦；加载后按 perms 校验
      if (this.isMerchant && (!this.perms || !this.perms.length)) return true
      return this.perms.includes(code)
    },
    async loadAuthMe() {
      if (!this.isAdmin) return
      const body = await fetchAuthMe()
      this.roles = body?.roles || []
      this.perms = body?.perms || []
      this.menuTree = body?.menu_tree || []
      localStorage.setItem('mymall_roles', JSON.stringify(this.roles))
      localStorage.setItem('mymall_perms', JSON.stringify(this.perms))
      localStorage.setItem('mymall_menus', JSON.stringify(this.menuTree))
    },
    async loadMerchantMe() {
      if (!this.isMerchant) return
      const body = await fetchMerchantMe()
      this.perms = body?.perms || []
      this.menuTree = body?.menu_tree || body?.menus || []
      this.isShopOwner = !!body?.is_owner
      localStorage.setItem('mymall_perms', JSON.stringify(this.perms))
      localStorage.setItem('mymall_menus', JSON.stringify(this.menuTree))
      localStorage.setItem('mymall_shop_owner', this.isShopOwner ? '1' : '0')
    },
    async login(mobile, password, shopId = 0) {
      const body = await http.post('/api/v1/user/login', {
        mobile,
        password,
        shop_id: shopId || undefined,
      })
      const token = body?.token
      const user = body?.user
      if (!token || !user) {
        throw new Error('登录响应异常，请重启 user-service 后重试')
      }
      const claims = parseJwt(token)
      this.token = token
      this.user = user
      this.role = claims.role || user?.role || ''
      this.shopId = Number(claims.shop_id || shopId || 0)
      this.userId = Number(claims.user_id || user?.id || 0)
      localStorage.setItem('mymall_token', token)
      localStorage.setItem('mymall_user', JSON.stringify(user))
      localStorage.setItem('mymall_role', this.role)
      localStorage.setItem('mymall_shop_id', String(this.shopId))
      localStorage.setItem('mymall_user_id', String(this.userId))
      if (this.isAdmin) {
        try {
          await this.loadAuthMe()
        } catch (e) {
          console.warn('loadAuthMe failed', e)
        }
      }
      if (this.isMerchant) {
        try {
          await this.loadMerchantMe()
        } catch (e) {
          console.warn('loadMerchantMe failed', e)
        }
      }
      return this.role
    },
    logout() {
      this.token = ''
      this.user = null
      this.role = ''
      this.shopId = 0
      this.userId = 0
      this.roles = []
      this.perms = []
      this.menuTree = []
      this.isShopOwner = false
      localStorage.clear()
    },
  },
})
