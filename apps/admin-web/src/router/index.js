import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Login from '../views/Login.vue'
import AdminLayout from '../layouts/AdminLayout.vue'
import MerchantLayout from '../layouts/MerchantLayout.vue'
import AdminShops from '../views/admin/Shops.vue'
import AdminApplications from '../views/admin/Applications.vue'
import AdminOrders from '../views/admin/Orders.vue'
import MerchantProducts from '../views/merchant/Products.vue'
import MerchantOrders from '../views/merchant/Orders.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { role: 'admin' },
      children: [
        { path: '', redirect: '/admin/applications' },
        { path: 'applications', component: AdminApplications },
        { path: 'shops', component: AdminShops },
        { path: 'orders', component: AdminOrders },
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

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.path === '/login') return true
  if (!auth.token) return '/login'
  if (to.meta.role === 'admin' && !auth.isAdmin) return '/login'
  if (to.meta.role === 'merchant' && !auth.isMerchant) return '/login'
  return true
})

export default router
