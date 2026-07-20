<template>
  <view class="page">
    <view class="top">
      <text class="home" @tap="goHome">← 首页</text>
      <text class="tip">审核通过后才会在社区展示</text>
      <text class="btn" @tap="goPublish()">发笔记</text>
    </view>
    <view v-for="a in list" :key="a.id" class="card">
      <image class="cover" :src="a.cover_url || placeholder" mode="aspectFill" @tap="goDetail(a)" />
      <view class="body">
        <text class="title line-2">{{ a.title }}</text>
        <view class="meta">
          <text class="tag" :class="a.audit_status">{{ statusText(a) }}</text>
          <text class="time">{{ a.created_at }}</text>
        </view>
        <view class="ops">
          <text v-if="a.audit_status !== 'approved'" class="link" @tap="goPublish(a.id)">编辑</text>
          <text v-if="a.audit_status !== 'approved'" class="link danger" @tap="remove(a)">删除</text>
          <text v-if="a.audit_status === 'approved' && a.status === 'published'" class="link" @tap="goDetail(a)">查看</text>
        </view>
        <text v-if="a.audit_status === 'rejected' && a.reject_reason" class="reason">驳回：{{ a.reject_reason }}</text>
      </view>
    </view>
    <view class="foot">
      <text v-if="loading">加载中...</text>
      <text v-else-if="!list.length">还没有笔记，去发一篇吧</text>
      <text v-else-if="finished">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onReachBottom, onShow } from '@dcloudio/uni-app'
import { deleteMyArticle, listMyArticles } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/102/400/300'
const list = ref([])
const page = ref(1)
const loading = ref(false)
const finished = ref(false)

function statusText(a) {
  if (a.audit_status === 'pending') return '待审核'
  if (a.audit_status === 'rejected') return '已驳回'
  if (a.status === 'published') return '已发布'
  return a.audit_status || a.status
}

function goHome() {
  uni.switchTab({ url: '/pages/index/index' })
}

function goPublish(id) {
  uni.navigateTo({ url: id ? `/pages/community/publish?id=${id}` : '/pages/community/publish' })
}

function goDetail(a) {
  if (a.audit_status === 'approved' && a.status === 'published') {
    uni.navigateTo({ url: `/pages/community/detail?id=${a.id}` })
  }
}

async function remove(a) {
  uni.showModal({
    title: '删除笔记',
    content: `确认删除「${a.title}」？`,
    success: async (r) => {
      if (!r.confirm) return
      try {
        await deleteMyArticle(a.id)
        uni.showToast({ title: '已删除', icon: 'none' })
        load(true)
      } catch (e) {
        uni.showToast({ title: e.message || '删除失败', icon: 'none' })
      }
    },
  })
}

async function load(reset = false) {
  if (!isLoggedIn()) {
    uni.navigateTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/community/mine') })
    return
  }
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      list.value = []
    }
    const res = await listMyArticles({ page: page.value, page_size: 10 })
    const rows = res?.list || []
    list.value = reset ? rows : list.value.concat(rows)
    const total = res?.total || 0
    if (list.value.length >= total || rows.length < 10) finished.value = true
    else page.value += 1
  } catch {
    if (reset) list.value = []
  } finally {
    loading.value = false
  }
}

onLoad(() => load(true))
onShow(() => { if (list.value.length) load(true) })
onReachBottom(() => load(false))
</script>

<style scoped>
.page { padding: 16rpx 24rpx 40rpx; }
.top {
  display: flex; justify-content: space-between; align-items: center;
  gap: 16rpx; margin-bottom: 16rpx;
}
.home { font-size: 26rpx; color: #c4894a; flex-shrink: 0; }
.tip { flex: 1; font-size: 22rpx; color: #a1a1aa; text-align: center; }
.btn { font-size: 26rpx; color: #c4894a; flex-shrink: 0; }
.card {
  display: flex; gap: 20rpx; background: #fff; border-radius: 16rpx;
  padding: 20rpx; margin-bottom: 16rpx;
}
.cover { width: 160rpx; height: 160rpx; border-radius: 12rpx; background: #f4f4f5; flex-shrink: 0; }
.body { flex: 1; min-width: 0; }
.title { font-size: 28rpx; color: #18181b; font-weight: 600; }
.meta { display: flex; gap: 12rpx; align-items: center; margin-top: 12rpx; }
.tag { font-size: 20rpx; padding: 2rpx 10rpx; border-radius: 6rpx; background: #f4f4f5; color: #71717a; }
.tag.pending { background: #fef3c7; color: #d97706; }
.tag.rejected { background: #fee2e2; color: #dc2626; }
.tag.approved { background: #dcfce7; color: #16a34a; }
.time { font-size: 20rpx; color: #a1a1aa; }
.ops { display: flex; gap: 24rpx; margin-top: 16rpx; }
.link { font-size: 24rpx; color: #c4894a; }
.link.danger { color: #dc2626; }
.reason { display: block; margin-top: 8rpx; font-size: 22rpx; color: #dc2626; }
.foot { text-align: center; color: #a1a1aa; padding: 24rpx; font-size: 24rpx; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2; overflow: hidden;
}
</style>
