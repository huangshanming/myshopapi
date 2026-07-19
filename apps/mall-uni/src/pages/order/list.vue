<template>
  <view class="page">
    <view v-if="!list.length && !loading" class="empty">暂无订单</view>
    <view v-for="o in list" :key="o.id" class="card" @tap="goDetail(o.id)">
      <view class="row">
        <text class="no">{{ o.order_no }}</text>
        <text class="status">{{ statusText(o.status) }}</text>
      </view>
      <view v-for="it in o.items || []" :key="it.id" class="item">
        <text class="line-1">{{ it.product_name }}</text>
        <text class="sub">×{{ it.quantity }} · ¥{{ it.price }}</text>
      </view>
      <view class="foot">
        <text class="sub">{{ o.created_at }}</text>
        <text class="price">¥{{ o.total_amount }}</text>
      </view>
      <view v-if="o.status === 'shipped' || o.status === 'completed' || o.status === 'reviewed'" class="actions" @tap.stop>
        <text v-if="o.status === 'shipped'" class="act" @tap="goDetail(o.id)">确认收货</text>
        <text v-if="o.status === 'completed'" class="act" @tap="goReview(o.id)">去评价</text>
        <text v-if="o.status === 'reviewed'" class="act muted" @tap="goViewReview(o.id)">查看评价</text>
      </view>
    </view>
    <view class="more">
      <text v-if="loading">加载中...</text>
      <text v-else-if="finished && list.length">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onReachBottom, onShow } from '@dcloudio/uni-app'
import { listOrders, ORDER_STATUS } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const list = ref([])
const page = ref(1)
const loading = ref(false)
const finished = ref(false)

function statusText(s) {
  return ORDER_STATUS[s] || s
}

function goDetail(id) {
  uni.navigateTo({ url: `/pages/order/detail?id=${id}` })
}

function goReview(id) {
  uni.navigateTo({ url: `/pages/order/review?id=${id}` })
}

function goViewReview(id) {
  uni.navigateTo({ url: `/pages/order/review-view?id=${id}` })
}

async function load(reset = false) {
  if (!isLoggedIn()) {
    uni.redirectTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/order/list') })
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
    const res = await listOrders({ page: page.value, page_size: 10 })
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

onShow(() => load(true))
onReachBottom(() => load(false))
</script>

<style scoped>
.page { padding: 24rpx 32rpx; }
.card {
  background: #fff; border-radius: 24rpx; padding: 24rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.row { display: flex; justify-content: space-between; margin-bottom: 16rpx; }
.no { font-size: 24rpx; color: #71717a; }
.status { color: #c8a876; font-size: 24rpx; }
.item { margin-bottom: 8rpx; }
.item .line-1 { font-size: 28rpx; display: block; }
.sub { color: #71717a; font-size: 22rpx; }
.foot { display: flex; justify-content: space-between; align-items: center; margin-top: 16rpx; }
.price { color: #d83636; font-weight: 700; font-size: 30rpx; }
.actions { display: flex; justify-content: flex-end; gap: 24rpx; margin-top: 16rpx; padding-top: 16rpx; border-top: 1rpx solid #f5f5f5; }
.act { color: #c8a876; font-size: 24rpx; }
.act.muted { color: #71717a; }
.empty, .more { text-align: center; color: #71717a; padding: 40rpx; font-size: 24rpx; }
</style>
