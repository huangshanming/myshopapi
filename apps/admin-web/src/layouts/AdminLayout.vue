<template>
  <el-container class="layout">
    <el-aside width="220px" class="aside">
      <div class="brand">mymall 平台</div>
      <el-menu :default-active="$route.path" router>
        <template v-for="item in menuNodes" :key="item.id">
          <el-sub-menu v-if="item.children?.length" :index="'dir-' + item.id">
            <template #title>{{ item.name }}</template>
            <el-menu-item
              v-for="child in item.children"
              :key="child.id"
              :index="child.path || ('/admin/' + child.id)"
            >{{ child.name }}</el-menu-item>
          </el-sub-menu>
          <el-menu-item v-else-if="item.path" :index="item.path">{{ item.name }}</el-menu-item>
        </template>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <span>{{ auth.user?.nickname || auth.user?.mobile }}</span>
        <el-button link type="danger" @click="onLogout">退出</el-button>
      </el-header>
      <el-main><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const menuNodes = computed(() => {
  const tree = auth.menuTree || []
  if (tree.length) {
    return tree
      .filter((n) => n.type !== 'button' && n.visible !== 0)
      .map((n) => ({
        ...n,
        children: (n.children || []).filter((c) => c.type === 'menu' && c.visible !== 0),
      }))
  }
  // fallback 静态
  return [
    {
      id: 1,
      name: '业务管理',
      children: [
        { id: 2, name: '入驻审核', path: '/admin/applications' },
        { id: 3, name: '店铺管理', path: '/admin/shops' },
        { id: 4, name: '全站订单', path: '/admin/orders' },
      ],
    },
  ]
})

onMounted(async () => {
  if (auth.isAdmin && (!auth.menuTree || !auth.menuTree.length)) {
    try {
      await auth.loadAuthMe()
    } catch (_) {
      /* ignore */
    }
  }
})

function onLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout { min-height: 100vh; }
.aside { background: #1f2d3d; color: #fff; }
.brand { padding: 20px 16px; font-weight: 700; color: #fff; }
.aside :deep(.el-menu) { border-right: none; background: transparent; }
.aside :deep(.el-menu-item),
.aside :deep(.el-sub-menu__title) { color: #c0c4cc; }
.aside :deep(.el-menu-item.is-active) { background: #304156; color: #fff; }
.header { display: flex; justify-content: flex-end; align-items: center; gap: 12px; border-bottom: 1px solid #eee; }
</style>
