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

    <button
      v-if="canCancel"
      class="btn-cancel"
      :loading="cancelling"
      @tap="onCancel"
    >取消订单</button>
  </view>
  <view v-else class="empty">{{ loading ? '加载中...' : '订单不存在' }}</view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { cancelOrder, getOrder, ORDER_STATUS } from '../../api/index'

const order = ref(null)
const loading = ref(false)
const cancelling = ref(false)
let orderId = 0

const canCancel = computed(() => {
  const s = order.value?.status
  return s === 'pending' || s === 'confirmed'
})

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
    order.value = res.data || null
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
  cancelling.value = true
  try {
    await cancelOrder(orderId)
    uni.showToast({ title: '已取消', icon: 'success' })
    load()
  } catch {
    /* handled */
  } finally {
    cancelling.value = false
  }
}
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; }
.card {
  background: #fff; border-radius: 24rpx; padding: 28rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.row { display: flex; justify-content: space-between; padding: 12rpx 0; font-size: 26rpx; }
.label { color: #71717a; }
.gold { color: #c8a876; }
.price { color: #d83636; font-weight: 700; }
.sec { font-weight: 600; font-size: 28rpx; display: block; margin-bottom: 16rpx; }
.item { padding: 12rpx 0; border-top: 1rpx solid #f5f5f5; }
.name { display: block; font-size: 28rpx; }
.sub { color: #71717a; font-size: 22rpx; }
.btn-cancel {
  margin-top: 24rpx; background: #fff; color: #d83636; border: 2rpx solid #d83636;
  border-radius: 999rpx; height: 80rpx; line-height: 80rpx; font-size: 28rpx;
}
.empty { text-align: center; padding: 120rpx; color: #71717a; }
</style>
