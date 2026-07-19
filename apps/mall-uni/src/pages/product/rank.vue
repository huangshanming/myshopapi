<template>
  <view class="page">
    <view class="tip">优先按今日销量排序，今日无销量时按总销量</view>
    <view
      v-for="(r, i) in list"
      :key="r.id"
      class="rank-item"
      @tap="goDetail(r.id)"
    >
      <view class="rank-no" :class="rankClass(i)">{{ i + 1 }}</view>
      <image class="rank-img" :src="r.main_image || placeholder" mode="aspectFill" />
      <view class="rank-info">
        <text class="line-2 rank-name">{{ displayName(r) }}</text>
        <view class="rank-foot">
          <text class="price">¥{{ r.sale_price }}</text>
          <text class="sub">{{ soldText(r) }}</text>
        </view>
      </view>
    </view>
    <view class="foot">
      <text v-if="loading">加载中...</text>
      <text v-else-if="!list.length">暂无上榜商品</text>
      <text v-else-if="finished">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onReachBottom, onShow } from '@dcloudio/uni-app'
import { listSalesRank } from '../../api/index'

const placeholder = 'https://picsum.photos/id/96/200/200'
const list = ref([])
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const finished = ref(false)

function displayName(r) {
  if (r.shop_name) return `${r.name} | ${r.shop_name}`
  return r.name
}

function soldText(r) {
  if (r.today_sold > 0) return `今日售出${r.today_sold}`
  return `总销量${r.sold_count || 0}`
}

function rankClass(i) {
  const n = i + 1
  if (n <= 3) return 'r' + n
  return ''
}

function goDetail(id) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
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
    const res = await listSalesRank({ page: page.value, page_size: pageSize })
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
.page { padding: 24rpx 32rpx 40rpx; }
.tip { font-size: 22rpx; color: #a1a1aa; margin-bottom: 20rpx; }
.rank-item {
  display: flex; align-items: center; gap: 20rpx;
  background: #fff; border-radius: 24rpx; padding: 20rpx; margin-bottom: 16rpx;
}
.rank-no {
  width: 40rpx; height: 40rpx; border-radius: 50%; color: #fff; font-size: 22rpx; font-weight: 700;
  display: flex; align-items: center; justify-content: center; background: #c8a876; flex-shrink: 0;
}
.rank-no.r1 { background: linear-gradient(135deg,#d83636,#f25757); }
.rank-no.r2 { background: linear-gradient(135deg,#f59e0b,#fbbf24); }
.rank-no.r3 { background: linear-gradient(135deg,#94a3b8,#cbd5e1); }
.rank-img { width: 128rpx; height: 128rpx; border-radius: 16rpx; background: #f4f4f5; flex-shrink: 0; }
.rank-info { flex: 1; min-width: 0; }
.rank-name { font-size: 26rpx; color: #18181b; font-weight: 500; }
.rank-foot { display: flex; justify-content: space-between; align-items: baseline; margin-top: 12rpx; }
.price { color: #d83636; font-weight: 700; font-size: 30rpx; }
.sub { color: #71717a; font-size: 22rpx; }
.foot { text-align: center; color: #a1a1aa; font-size: 24rpx; padding: 24rpx 0; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2; overflow: hidden;
}
</style>
