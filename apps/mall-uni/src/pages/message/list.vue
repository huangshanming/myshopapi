<template>
  <view class="page">
    <view class="toolbar">
      <text class="title">消息中心</text>
      <text v-if="list.length" class="read-all" @tap="readAll">全部已读</text>
    </view>
    <view
      v-for="m in list"
      :key="m.id"
      class="item"
      :class="{ unread: !m.is_read }"
      @tap="openMsg(m)"
    >
      <view class="row">
        <text class="name">{{ m.title }}</text>
        <view v-if="!m.is_read" class="dot" />
      </view>
      <text class="content">{{ m.content }}</text>
      <text class="time">{{ m.created_at }}</text>
    </view>
    <view v-if="!list.length && !loading" class="empty">暂无消息</view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import {
  listNotifications, markAllNotificationsRead, markNotificationRead,
} from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const list = ref([])
const loading = ref(false)

async function load() {
  if (!isLoggedIn()) {
    uni.redirectTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/message/list') })
    return
  }
  loading.value = true
  try {
    const res = await listNotifications({ page: 1, page_size: 50 })
    list.value = res?.list || []
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

async function openMsg(m) {
  if (!m.is_read) {
    try { await markNotificationRead(m.id) } catch { /* ignore */ }
    m.is_read = 1
  }
  if (m.link_type === 'order' && m.link_id) {
    uni.navigateTo({ url: `/pages/order/detail?id=${m.link_id}` })
    return
  }
  if (m.link_type === 'article' && m.link_id) {
    let commentId = 0
    try {
      const extra = typeof m.extra === 'string' ? JSON.parse(m.extra || '{}') : (m.extra || {})
      commentId = Number(extra.comment_id) || 0
    } catch { /* ignore */ }
    const q = commentId ? `&comment_id=${commentId}` : ''
    uni.navigateTo({ url: `/pages/community/detail?id=${m.link_id}${q}` })
  }
}

async function readAll() {
  try {
    await markAllNotificationsRead()
    list.value.forEach((x) => { x.is_read = 1 })
    uni.showToast({ title: '已全部已读', icon: 'none' })
  } catch { /* toast */ }
}

onShow(load)
</script>

<style scoped>
.page { min-height: 100vh; background: #f7f3ec; padding-bottom: 40rpx; }
.toolbar {
  display: flex; justify-content: space-between; align-items: center;
  padding: 24rpx 28rpx; background: #fff;
}
.title { font-size: 32rpx; font-weight: 700; }
.read-all { font-size: 26rpx; color: #c4894a; }
.item {
  margin: 16rpx 24rpx 0; background: #fff; border-radius: 16rpx; padding: 24rpx;
}
.item.unread { box-shadow: 0 0 0 2rpx rgba(196,137,74,.25); }
.row { display: flex; align-items: center; justify-content: space-between; }
.name { font-size: 30rpx; font-weight: 600; color: #18181b; }
.item.unread .name { font-weight: 700; }
.dot { width: 14rpx; height: 14rpx; border-radius: 50%; background: #ef4444; }
.content {
  display: block; margin-top: 10rpx; font-size: 26rpx; color: #52525b; line-height: 1.5;
}
.time { display: block; margin-top: 12rpx; font-size: 22rpx; color: #a1a1aa; }
.empty { text-align: center; color: #94a3b8; padding: 100rpx 0; }
</style>
