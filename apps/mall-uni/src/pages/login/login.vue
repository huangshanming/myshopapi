<template>
  <view class="page">
    <view class="card">
      <text class="title">登录 mymall</text>
      <view class="field">
        <text class="label">手机号</text>
        <input v-model="mobile" type="number" maxlength="11" placeholder="请输入手机号" />
      </view>
      <view class="field">
        <text class="label">密码</text>
        <input v-model="password" password placeholder="请输入密码" />
      </view>
      <button class="btn" :loading="loading" @tap="onLogin">登录</button>
      <button class="btn ghost" :loading="loading" @tap="onRegister">注册并登录</button>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { login, register } from '../../api/index'
import { setAuth } from '../../stores/user'

const mobile = ref('')
const password = ref('')
const loading = ref(false)
let redirect = ''

onLoad((q) => {
  redirect = q.redirect ? decodeURIComponent(q.redirect) : ''
})

async function onLogin() {
  if (!mobile.value || !password.value) {
    uni.showToast({ title: '请填写手机号和密码', icon: 'none' })
    return
  }
  loading.value = true
  try {
    const res = await login({ mobile: mobile.value, password: password.value })
    setAuth(res?.token, res?.user)
    uni.showToast({ title: '登录成功', icon: 'success' })
    setTimeout(() => {
      if (redirect) uni.redirectTo({ url: redirect })
      else uni.switchTab({ url: '/pages/mine/mine' })
    }, 400)
  } catch {
    /* handled */
  } finally {
    loading.value = false
  }
}

async function onRegister() {
  if (!mobile.value || !password.value) {
    uni.showToast({ title: '请填写手机号和密码', icon: 'none' })
    return
  }
  loading.value = true
  try {
    await register({ mobile: mobile.value, password: password.value })
    const res = await login({ mobile: mobile.value, password: password.value })
    setAuth(res?.token, res?.user)
    uni.showToast({ title: '注册成功', icon: 'success' })
    setTimeout(() => {
      if (redirect) uni.redirectTo({ url: redirect })
      else uni.switchTab({ url: '/pages/mine/mine' })
    }, 400)
  } catch {
    /* handled */
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.page { padding: 48rpx 32rpx; }
.card {
  background: #fff; border-radius: 32rpx; padding: 40rpx 32rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.title { font-size: 40rpx; font-weight: 700; color: #1d1d1f; }
.field { margin-top: 32rpx; }
.label { display: block; font-size: 24rpx; color: #71717a; margin-bottom: 12rpx; }
input {
  background: #fafafa; border-radius: 16rpx; padding: 20rpx 24rpx; font-size: 28rpx;
}
.btn {
  margin-top: 40rpx; height: 88rpx; line-height: 88rpx; border-radius: 999rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; font-size: 30rpx;
}
.btn.ghost {
  margin-top: 20rpx; background: #fff; color: #c8a876; border: 2rpx solid #c8a876;
}
</style>
