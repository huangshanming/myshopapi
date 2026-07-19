<template>
  <view class="page" v-if="order">
    <view class="card">
      <view class="row">
        <text class="label">订单号</text>
        <text>{{ order.order_no }}</text>
      </view>
      <view class="row">
        <text class="label">状态</text>
        <text class="gold">{{ statusText(order.status) }}</text>
      </view>
      <view class="row">
        <text class="label">金额</text>
        <text class="price">¥{{ order.total_amount }}</text>
      </view>
      <view class="row">
        <text class="label">下单时间</text>
        <text>{{ order.created_at }}</text>
      </view>
      <view v-if="order.receiver_name || order.receiver_address" class="row col">
        <text class="label">收货信息</text>
        <text class="addr">{{ order.receiver_name }} {{ order.receiver_phone }}</text>
        <text class="addr">{{ order.receiver_address }}</text>
      </view>
      <view v-if="order.ship_company" class="row">
        <text class="label">物流</text>
        <text>{{ order.ship_company }} {{ order.ship_no }}</text>
      </view>
    </view>

    <view class="card">
      <text class="sec">商品</text>
      <view v-for="it in order.items || []" :key="it.id" class="item">
        <text class="name">{{ it.product_name }}</text>
        <text class="sub">×{{ it.quantity }} · ¥{{ it.price }}</text>
      </view>
    </view>

    <view class="actions">
      <button v-if="canCancel" class="btn outline" :loading="busy" @tap="onCancel">取消订单</button>
      <button v-if="canConfirm" class="btn primary" :loading="busy" @tap="onConfirm">确认收货</button>
      <button v-if="canReview" class="btn primary" @tap="goReview">去评价</button>
      <button v-if="canViewReview" class="btn outline" @tap="goViewReview">查看评价</button>
    </view>
  </view>
  <view v-else class="empty">{{ loading ? '加载中...' : '订单不存在' }}</view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { cancelOrder, confirmReceive, getOrder, ORDER_STATUS } from '../../api/index'

const order = ref(null)
const loading = ref(false)
const busy = ref(false)
let orderId = 0

const canCancel = computed(() => {
  const s = order.value?.status
  return s === 'pending' || s === 'confirmed'
})
const canConfirm = computed(() => order.value?.status === 'shipped')
const canReview = computed(() => order.value?.status === 'completed')
const canViewReview = computed(() => order.value?.status === 'reviewed')

function statusText(s) {
  return ORDER_STATUS[s] || s
}

onLoad((q) => {
  orderId = Number(q.id || 0)
  load()
})

async function load() {
  if (!orderId) return
  loading.value = true
  try {
    const res = await getOrder(orderId)
    order.value = res || null
  } catch {
    order.value = null
  } finally {
    loading.value = false
  }
}

async function onCancel() {
  const ok = await new Promise((resolve) => {
    uni.showModal({
      title: '取消订单',
      content: '确认取消该订单？',
      success: (r) => resolve(r.confirm),
    })
  })
  if (!ok) return
  busy.value = true
  try {
    await cancelOrder(orderId)
    uni.showToast({ title: '已取消', icon: 'success' })
    load()
  } catch {
    /* handled */
  } finally {
    busy.value = false
  }
}

async function onConfirm() {
  const ok = await new Promise((resolve) => {
    uni.showModal({
      title: '确认收货',
      content: '确认已收到商品？',
      success: (r) => resolve(r.confirm),
    })
  })
  if (!ok) return
  busy.value = true
  try {
    await confirmReceive(orderId)
    uni.showToast({ title: '已确认收货', icon: 'success' })
    load()
  } catch {
    /* handled */
  } finally {
    busy.value = false
  }
}

function goReview() {
  uni.navigateTo({ url: `/pages/order/review?id=${orderId}` })
}

function goViewReview() {
  uni.navigateTo({ url: `/pages/order/review-view?id=${orderId}` })
}
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; }
.card {
  background: #fff; border-radius: 24rpx; padding: 28rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.row { display: flex; justify-content: space-between; padding: 12rpx 0; font-size: 26rpx; }
.row.col { flex-direction: column; align-items: flex-start; gap: 8rpx; }
.addr { color: #3f3f46; line-height: 1.4; }
.label { color: #71717a; }
.gold { color: #c8a876; }
.price { color: #d83636; font-weight: 700; }
.sec { font-weight: 600; font-size: 28rpx; display: block; margin-bottom: 16rpx; }
.item { padding: 12rpx 0; border-top: 1rpx solid #f5f5f5; }
.name { display: block; font-size: 28rpx; }
.sub { color: #71717a; font-size: 22rpx; }
.actions { display: flex; flex-direction: column; gap: 16rpx; margin-top: 8rpx; }
.btn {
  border-radius: 999rpx; height: 80rpx; line-height: 80rpx; font-size: 28rpx; margin: 0;
}
.btn.outline { background: #fff; color: #d83636; border: 2rpx solid #d83636; }
.btn.primary { background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; border: none; }
.empty { text-align: center; padding: 120rpx; color: #71717a; }
</style>
