import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Login from '../views/Login.vue'
import Forbidden from '../views/Forbidden.vue'
import AdminLayout from '../layouts/AdminLayout.vue'
import MerchantLayout from '../layouts/MerchantLayout.vue'
import AdminShops from '../views/admin/Shops.vue'
import AdminApplications from '../views/admin/Applications.vue'
import AdminOrders from '../views/admin/Orders.vue'
import SystemMenus from '../views/admin/system/Menus.vue'
import SystemRoles from '../views/admin/system/Roles.vue'
import SystemUsers from '../views/admin/system/Users.vue'
import SystemAdmins from '../views/admin/system/Admins.vue'
import SystemConfigs from '../views/admin/system/Configs.vue'
import MerchantProducts from '../views/merchant/Products.vue'
import MerchantOrders from '../views/merchant/Orders.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login },
    { path: '/403', component: Forbidden },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { role: 'admin' },
      children: [
        { path: '', redirect: '/admin/applications' },
        { path: 'applications', component: AdminApplications, meta: { perms: 'business:application:list' } },
        { path: 'shops', component: AdminShops, meta: { perms: 'business:shop:list' } },
        { path: 'orders', component: AdminOrders, meta: { perms: 'business:order:list' } },
        { path: 'system/menus', component: SystemMenus, meta: { perms: 'system:menu:list' } },
        { path: 'system/roles', component: SystemRoles, meta: { perms: 'system:role:list' } },
        { path: 'system/users', component: SystemUsers, meta: { perms: 'system:user:list' } },
        { path: 'system/admins', component: SystemAdmins, meta: { perms: 'system:admin:list' } },
        { path: 'system/configs', component: SystemConfigs, meta: { perms: 'system:config:list' } },
      ],
    },
    {
      path: '/merchant',
      component: MerchantLayout,
      meta: { role: 'merchant' },
      children: [
        { path: '', redirect: '/merchant/products' },
        { path: 'products', component: MerchantProducts },
        { path: 'orders', component: MerchantOrders },
      ],
    },
    { path: '/', redirect: '/login' },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.path === '/login' || to.path === '/403') return true
  if (!auth.token) return '/login'
  if (to.meta.role === 'admin' && !auth.isAdmin) return '/login'
  if (to.meta.role === 'merchant' && !auth.isMerchant) return '/login'
  if (to.meta.role === 'admin' && auth.isAdmin && !auth.menuTree?.length) {
    try {
      await auth.loadAuthMe()
    } catch (_) {
      /* ignore */
    }
  }
  if (to.meta.perms && !auth.hasPerm(to.meta.perms)) {
    return '/403'
  }
  return true
})

export default router
