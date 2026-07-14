<template>
  <el-container class="layout">
    <el-aside width="200px" class="aside">
      <div class="brand">mymall 平台</div>
      <el-menu :default-active="$route.path" router>
        <el-menu-item index="/admin/applications">入驻审核</el-menu-item>
        <el-menu-item index="/admin/shops">店铺管理</el-menu-item>
        <el-menu-item index="/admin/orders">全站订单</el-menu-item>
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
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
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
.aside :deep(.el-menu-item) { color: #c0c4cc; }
.aside :deep(.el-menu-item.is-active) { background: #304156; color: #fff; }
.header { display: flex; justify-content: flex-end; align-items: center; gap: 12px; border-bottom: 1px solid #eee; }
</style>
