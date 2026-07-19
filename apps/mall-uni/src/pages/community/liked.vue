<template>
  <view class="page">
    <view
      v-for="a in list"
      :key="a.id"
      class="card"
      :class="{ invalid: a.invalid }"
      @tap="goDetail(a)"
    >
      <image class="cover" :src="a.cover_url || placeholder" mode="aspectFill" />
      <view class="body">
        <text class="title line-2">{{ a.title || '文章已失效' }}</text>
        <view class="stats">
          <text v-if="a.invalid" class="tag">已失效</text>
          <template v-else>
            <text>阅 {{ a.read_count || 0 }}</text>
            <text>❤ {{ a.like_count || 0 }}</text>
          </template>
        </view>
        <text class="cancel" @tap.stop="cancelLike(a)">取消点赞</text>
      </view>
    </view>

    <view class="foot">
      <text v-if="loading">加载中...</text>
      <text v-else-if="!list.length">暂无点赞</text>
      <text v-else-if="finished">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onReachBottom, onShow } from '@dcloudio/uni-app'
import { listMyArticleLikes, unlikeArticle } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/102/400/300'
const list = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

function goDetail(a) {
  if (a.invalid) {
    uni.showToast({ title: '文章已失效', icon: 'none' })
    return
  }
  uni.navigateTo({ url: `/pages/community/detail?id=${a.id}` })
}

async function cancelLike(a) {
  try {
    await unlikeArticle(a.id)
    list.value = list.value.filter((x) => x.id !== a.id)
    uni.showToast({ title: '已取消', icon: 'none' })
  } catch { /* handled */ }
}

async function load(reset = false) {
  if (!isLoggedIn()) {
    uni.redirectTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/community/liked') })
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
    const res = await listMyArticleLikes({ page: page.value, page_size: pageSize })
    const rows = res?.list || []
    list.value = reset ? rows : list.value.concat(rows)
    const total = res?.total || 0
    if (list.value.length >= total || rows.length < pageSize) finished.value = true
    else page.value += 1
  } catch {
    if (reset) list.value = []
  } finally {
    loading.value = false
  }
}

onShow(() => load(true))
onReachBottom(() => load(false))
</script>

<style scoped>
.page { padding: 16rpx 24rpx 40rpx; }
.card {
  display: flex; gap: 20rpx; background: #fff; border-radius: 16rpx;
  padding: 20rpx; margin-bottom: 16rpx;
}
.card.invalid { opacity: 0.55; }
.cover { width: 200rpx; height: 150rpx; border-radius: 12rpx; flex-shrink: 0; background: #f4f4f5; }
.body { flex: 1; min-width: 0; display: flex; flex-direction: column; justify-content: space-between; }
.title { font-size: 28rpx; color: #18181b; font-weight: 600; }
.stats { display: flex; gap: 20rpx; font-size: 22rpx; color: #a1a1aa; margin-top: 12rpx; }
.tag { color: #a1a1aa; font-size: 22rpx; }
.cancel { display: block; margin-top: 12rpx; color: #71717a; font-size: 22rpx; }
.foot { text-align: center; color: #a1a1aa; font-size: 24rpx; padding: 24rpx 0; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2;
  overflow: hidden;
}
</style>
