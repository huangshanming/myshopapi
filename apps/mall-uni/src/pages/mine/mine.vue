<template>
  <view class="page">
    <view class="profile card">
      <view class="avatar">{{ avatarText }}</view>
      <view v-if="user" class="info">
        <text class="name">{{ user.nickname || user.mobile }}</text>
        <text class="sub">{{ user.mobile }}</text>
      </view>
      <view v-else class="info" @tap="goLogin">
        <text class="name">未登录</text>
        <text class="sub">点击登录 / 注册</text>
      </view>
    </view>

    <view class="menu card">
      <view class="menu-item" @tap="goOrders">
        <text>我的订单</text>
        <text class="arrow">›</text>
      </view>
      <view v-if="user" class="menu-item" @tap="logout">
        <text>退出登录</text>
        <text class="arrow">›</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { clearAuth, getUser, isLoggedIn } from '../../stores/user'

const user = ref(null)
const avatarText = computed(() => {
  const n = user.value?.nickname || user.value?.mobile || '?'
  return String(n).slice(0, 1)
})

onShow(() => {
  user.value = getUser()
})

function goLogin() {
  uni.navigateTo({ url: '/pages/login/login' })
}

function goOrders() {
  if (!isLoggedIn()) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/order/list' })
}

function logout() {
  clearAuth()
  user.value = null
  uni.showToast({ title: '已退出', icon: 'none' })
}
</script>

<style scoped>
.page { padding: 32rpx; }
.card {
  background: #fff; border-radius: 32rpx; padding: 32rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08); margin-bottom: 24rpx;
}
.profile { display: flex; align-items: center; gap: 24rpx; }
.avatar {
  width: 112rpx; height: 112rpx; border-radius: 50%;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff;
  display: flex; align-items: center; justify-content: center; font-size: 44rpx; font-weight: 700;
}
.name { display: block; font-size: 34rpx; font-weight: 600; }
.sub { display: block; color: #71717a; font-size: 24rpx; margin-top: 8rpx; }
.menu-item {
  display: flex; justify-content: space-between; align-items: center;
  padding: 28rpx 0; border-bottom: 1rpx solid #f0f0f0; font-size: 28rpx;
}
.menu-item:last-child { border-bottom: none; }
.arrow { color: #c8a876; font-size: 36rpx; }
</style>
