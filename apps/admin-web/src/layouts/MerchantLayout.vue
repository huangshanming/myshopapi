<template>
  <el-container class="layout">
    <el-aside width="200px" class="aside">
      <div class="brand">商家后台</div>
      <el-menu :default-active="$route.path" router>
        <el-menu-item index="/merchant/products">商品管理</el-menu-item>
        <el-menu-item index="/merchant/orders">店铺订单</el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <span>店铺 #{{ auth.shopId || '-' }} · {{ auth.user?.nickname || auth.user?.mobile }}</span>
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
.aside { background: #0f3d3e; }
.brand { padding: 20px 16px; font-weight: 700; color: #e8f5f5; }
.aside :deep(.el-menu) { border-right: none; background: transparent; }
.aside :deep(.el-menu-item) { color: #a8c5c6; }
.aside :deep(.el-menu-item.is-active) { background: #155e63; color: #fff; }
.header { display: flex; justify-content: flex-end; align-items: center; gap: 12px; border-bottom: 1px solid #eee; }
</style>
