<template>
  <view class="page">
    <view class="head">
      <text class="title">限时秒杀</text>
      <view class="countdown">
        <text class="cd-label">距结束</text>
        <text class="count-num">{{ cd.h }}</text>:
        <text class="count-num">{{ cd.m }}</text>:
        <text class="count-num">{{ cd.s }}</text>
      </view>
    </view>

    <view class="list">
      <view
        v-for="p in items"
        :key="p.id"
        class="item"
        @tap="goDetail(p)"
      >
        <image class="img" :src="p.img || placeholder" mode="aspectFill" />
        <view class="body">
          <text class="name line-2">{{ p.name }}</text>
          <view class="price-row">
            <text class="price">¥{{ p.price }}</text>
            <text class="old">¥{{ p.old }}</text>
          </view>
          <text class="left">仅剩{{ p.left }}件</text>
        </view>
      </view>
    </view>

    <view class="foot">
      <text v-if="loading">加载中...</text>
      <text v-else-if="finished">没有更多了</text>
      <text v-else @tap="loadMore">上拉加载更多</text>
    </view>
  </view>
</template>

<script setup>
import { onMounted, onUnmounted, reactive, ref } from 'vue'
import { onReachBottom } from '@dcloudio/uni-app'
import { listSeckill } from '../../api/index'

const placeholder = 'https://picsum.photos/id/96/400/400'
const items = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)
const endAt = ref(Date.now() + 3600e3)
const cd = reactive({ h: '00', m: '00', s: '00' })
let timer

function pad(n) {
  return String(n).padStart(2, '0')
}

function tick() {
  let left = Math.max(0, endAt.value - Date.now())
  const h = Math.floor(left / 3600e3)
  left -= h * 3600e3
  const m = Math.floor(left / 60e3)
  left -= m * 60e3
  const s = Math.floor(left / 1e3)
  cd.h = pad(h)
  cd.m = pad(m)
  cd.s = pad(s)
}

function parseEndAt(v) {
  if (!v) return Date.now() + 3600e3
  const t = new Date(String(v).replace(/-/g, '/')).getTime()
  return Number.isFinite(t) ? t : Date.now() + 3600e3
}

async function load(reset = false) {
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      items.value = []
    }
    const res = await listSeckill({ page: page.value, page_size: pageSize })
    const data = res.data || {}
    if (data.end_at) {
      endAt.value = parseEndAt(data.end_at)
      tick()
    }
    const list = data.list || []
    items.value = reset ? list : items.value.concat(list)
    const total = data.total || 0
    if (items.value.length >= total || list.length < pageSize) {
      finished.value = true
    } else {
      page.value += 1
    }
  } catch {
    if (reset) items.value = []
  } finally {
    loading.value = false
  }
}

function loadMore() {
  load(false)
}

function goDetail(p) {
  uni.navigateTo({
    url: `/pages/product/detail?id=${p.product_id}&seckill_entry_id=${p.id}`,
  })
}

onReachBottom(() => loadMore())

onMounted(() => {
  tick()
  timer = setInterval(tick, 1000)
  load(true)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; }
.head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24rpx; }
.title { font-size: 36rpx; font-weight: 700; }
.countdown { display: flex; align-items: center; gap: 6rpx; font-size: 22rpx; color: #71717a; }
.cd-label { margin-right: 4rpx; }
.count-num {
  min-width: 40rpx; text-align: center; background: #18181b; color: #fff;
  border-radius: 8rpx; padding: 4rpx 6rpx; font-weight: 600;
}
.item {
  display: flex; gap: 20rpx; background: #fff; border-radius: 24rpx; padding: 20rpx;
  margin-bottom: 20rpx; box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.img { width: 200rpx; height: 200rpx; border-radius: 16rpx; background: #f3f3f3; flex-shrink: 0; }
.body { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.name { font-size: 28rpx; font-weight: 600; }
.price-row { margin-top: 16rpx; display: flex; align-items: baseline; gap: 12rpx; }
.price { color: #d83636; font-size: 36rpx; font-weight: 700; }
.old { color: #a1a1aa; font-size: 22rpx; text-decoration: line-through; }
.left { margin-top: 8rpx; font-size: 22rpx; color: #71717a; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2; overflow: hidden;
}
.foot { text-align: center; padding: 24rpx; color: #a1a1aa; font-size: 24rpx; }
</style>
