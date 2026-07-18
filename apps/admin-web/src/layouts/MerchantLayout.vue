<template>
  <el-container class="layout">
    <el-aside width="220px" class="aside">
      <div class="brand">商家后台</div>
      <el-menu :default-active="$route.path" :default-openeds="opened" router>
        <template v-for="item in menuNodes" :key="item.id">
          <el-sub-menu v-if="item.children?.length" :index="'dir-' + item.id">
            <template #title>{{ item.name }}</template>
            <el-menu-item
              v-for="child in item.children"
              :key="child.id"
              :index="child.path || ('/merchant/' + child.id)"
            >{{ child.name }}</el-menu-item>
          </el-sub-menu>
          <el-menu-item v-else-if="item.path" :index="item.path">{{ item.name }}</el-menu-item>
        </template>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <span class="shop-info">店铺 #{{ auth.shopId || '-' }} · {{ auth.user?.nickname || auth.user?.mobile }}</span>
        <div class="header-actions">
          <el-badge :value="unread" :hidden="!unread" :max="99" class="bell-badge">
            <el-button
              class="bell-btn"
              :class="{ 'bell-shake': unread > 0 }"
              circle
              @click="goNotifications"
            >
              <span class="bell-icon">🔔</span>
            </el-button>
          </el-badge>
          <el-button link type="danger" @click="onLogout">退出</el-button>
        </div>
      </el-header>
      <el-main><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { unreadNotificationCount } from '../api/merchant-notification'

const auth = useAuthStore()
const router = useRouter()
const unread = ref(0)
let pollTimer = null

const fallbackMenus = [
  {
    id: 100,
    name: '商品中心',
    type: 'dir',
    children: [
      { id: 1, name: '商品列表', path: '/merchant/products', type: 'menu' },
      { id: 2, name: '发布商品', path: '/merchant/products/edit', type: 'menu' },
      { id: 3, name: '回收站', path: '/merchant/products/recycle', type: 'menu' },
      { id: 7, name: '操作日志', path: '/merchant/products/op-logs', type: 'menu' },
    ],
  },
  {
    id: 101,
    name: '库存管理',
    type: 'dir',
    children: [
      { id: 4, name: '库存预警', path: '/merchant/stocks/warnings', type: 'menu' },
    ],
  },
  {
    id: 102,
    name: '订单中心',
    type: 'dir',
    children: [
      { id: 5, name: '店铺订单', path: '/merchant/orders', type: 'menu' },
      { id: 19, name: '售后管理', path: '/merchant/after-sales', type: 'menu' },
    ],
  },
  {
    id: 104,
    name: '内容中心',
    type: 'dir',
    children: [
      { id: 8, name: '我的文章', path: '/merchant/articles', type: 'menu' },
      { id: 9, name: '发布文章', path: '/merchant/articles/edit', type: 'menu' },
    ],
  },
  {
    id: 103,
    name: '店铺设置',
    type: 'dir',
    children: [
      { id: 6, name: '员工权限', path: '/merchant/staff', type: 'menu' },
      { id: 10, name: '消息通知', path: '/merchant/notifications', type: 'menu' },
    ],
  },
]

function normalizeTree(tree) {
  return (tree || [])
    .filter((n) => n.type !== 'button' && n.visible !== 0)
    .map((n) => ({
      ...n,
      children: (n.children || []).filter((c) => c.type === 'menu' && c.visible !== 0 && c.path),
    }))
    .filter((n) => n.children?.length || n.path)
}

const menuNodes = computed(() => {
  const tree = normalizeTree(auth.menuTree)
  return tree.length ? tree : fallbackMenus
})

const opened = computed(() => menuNodes.value.filter((n) => n.children?.length).map((n) => 'dir-' + n.id))

async function refreshUnread() {
  try {
    const res = await unreadNotificationCount()
    unread.value = Number(res.data?.count || 0)
  } catch (_) {
    /* ignore */
  }
}

function goNotifications() {
  router.push('/merchant/notifications')
}

onMounted(async () => {
  if (auth.isMerchant) {
    try {
      await auth.loadMerchantMe()
    } catch (_) {
      /* ignore */
    }
    await refreshUnread()
    pollTimer = setInterval(refreshUnread, 30000)
  }
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})

function onLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout { min-height: 100vh; }
.aside { background: #0f3d3e; }
.brand { padding: 20px 16px; font-weight: 700; color: #e8f5f5; }
.aside :deep(.el-menu) { border-right: none; background: transparent; }
.aside :deep(.el-menu-item),
.aside :deep(.el-sub-menu__title) { color: #a8c5c6; }
.aside :deep(.el-menu-item.is-active) { background: #155e63; color: #fff; }
.aside :deep(.el-sub-menu .el-menu-item) { min-width: auto; }
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #eee;
}
.header-actions { display: flex; align-items: center; gap: 12px; }
.bell-badge { line-height: 1; }
.bell-btn { width: 36px; height: 36px; padding: 0; }
.bell-icon { font-size: 16px; line-height: 1; }
.bell-shake {
  animation: bell-shake 2.4s ease-in-out infinite;
  transform-origin: top center;
}
@keyframes bell-shake {
  0% { transform: rotate(0deg); }
  4% { transform: rotate(16deg); }
  8% { transform: rotate(-14deg); }
  12% { transform: rotate(12deg); }
  16% { transform: rotate(-10deg); }
  20% { transform: rotate(8deg); }
  24% { transform: rotate(-4deg); }
  28% { transform: rotate(0deg); }
  100% { transform: rotate(0deg); }
}
</style>
