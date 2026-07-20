<template>
  <view class="page">
    <view class="top-bar">
      <text class="link" @tap="goMine">我的笔记</text>
      <text class="pub" @tap="goPublish">发笔记</text>
    </view>
    <view
      v-for="a in list"
      :key="a.id"
      class="card"
      @tap="goDetail(a.id)"
    >
      <image class="cover" :src="a.cover_url || placeholder" mode="aspectFill" />
      <view class="body">
        <view class="row">
          <text class="title line-2">{{ a.title }}</text>
          <text v-if="a.paid" class="badge">推广</text>
        </view>
        <view class="stats">
          <text>阅 {{ a.read_count || 0 }}</text>
          <text>❤ {{ a.like_count || 0 }}</text>
          <text>☆ {{ a.collect_count || 0 }}</text>
        </view>
      </view>
    </view>

    <view class="foot">
      <text v-if="loading">加载中...</text>
      <text v-else-if="!list.length">暂无种草内容</text>
      <text v-else-if="finished">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onReachBottom } from '@dcloudio/uni-app'
import { listArticles } from '../../api/index'

const placeholder = 'https://picsum.photos/id/102/400/300'
const list = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

function goDetail(id) {
  uni.navigateTo({ url: `/pages/community/detail?id=${id}` })
}
function goPublish() {
  uni.navigateTo({ url: '/pages/community/publish' })
}
function goMine() {
  uni.navigateTo({ url: '/pages/community/mine' })
}

async function load(reset = false) {
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      list.value = []
    }
    const res = await listArticles({ page: page.value, page_size: pageSize })
    const rows = res?.list || []
    list.value = reset ? rows : list.value.concat(rows)
    const total = res?.total || 0
    if (list.value.length >= total || rows.length < pageSize) {
      finished.value = true
    } else {
      page.value += 1
    }
  } catch {
    if (reset) list.value = []
  } finally {
    loading.value = false
  }
}

onLoad(() => load(true))
onReachBottom(() => load(false))
</script>

<style scoped>
.page { padding: 16rpx 24rpx 40rpx; }
.top-bar {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 16rpx; padding: 0 4rpx;
}
.link { font-size: 26rpx; color: #71717a; }
.pub {
  font-size: 26rpx; color: #fff; background: #c8a876;
  padding: 10rpx 28rpx; border-radius: 28rpx;
}
.card {
  display: flex; gap: 20rpx; background: #fff; border-radius: 16rpx;
  padding: 20rpx; margin-bottom: 16rpx;
}
.cover { width: 200rpx; height: 150rpx; border-radius: 12rpx; flex-shrink: 0; background: #f4f4f5; }
.body { flex: 1; min-width: 0; display: flex; flex-direction: column; justify-content: space-between; }
.row { display: flex; gap: 12rpx; align-items: flex-start; }
.title { flex: 1; font-size: 28rpx; color: #18181b; font-weight: 600; }
.badge {
  font-size: 20rpx; color: #c8a876; background: rgba(200,168,118,.15);
  padding: 2rpx 10rpx; border-radius: 6rpx; flex-shrink: 0;
}
.stats {
  display: flex; gap: 20rpx; font-size: 22rpx; color: #a1a1aa; margin-top: 12rpx;
}
.foot { text-align: center; color: #a1a1aa; font-size: 24rpx; padding: 24rpx 0; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2;
  overflow: hidden;
}
</style>
