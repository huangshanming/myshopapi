<template>
  <div class="login-page">
    <el-card class="card">
      <h1>mymall 后台登录</h1>
      <p class="tip">平台管理员或商家账号</p>
      <el-form @submit.prevent="onSubmit">
        <el-form-item label="手机号">
          <el-input v-model="mobile" placeholder="13900000001 / 13900000002" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" show-password />
        </el-form-item>
        <el-button type="primary" :loading="loading" style="width:100%" @click="onSubmit">登录</el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'

const mobile = ref('13900000001')
const password = ref('123456')
const loading = ref(false)
const router = useRouter()
const auth = useAuthStore()

async function onSubmit() {
  loading.value = true
  try {
    const role = await auth.login(mobile.value, password.value)
    ElMessage.success('登录成功')
    if (role === 'platform_admin' || role === 'admin') {
      router.push('/admin')
    } else {
      router.push('/merchant')
    }
  } catch (e) {
    ElMessage.error(e.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: linear-gradient(160deg, #e8f0f2 0%, #f7f3ea 100%);
}
.card { width: 380px; }
h1 { margin: 0 0 8px; font-size: 22px; }
.tip { color: #888; margin: 0 0 20px; font-size: 13px; }
</style>
