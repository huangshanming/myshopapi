<template>
  <view class="page" v-if="order">
    <view class="status-card">
      <text class="st">{{ statusText(order.status) }}</text>
      <text class="no">单号 {{ order.order_no }}</text>
    </view>

    <view class="panel">
      <view class="row">
        <image class="img" :src="order.product_cover || placeholder" mode="aspectFill" />
        <view class="body">
          <text class="name">{{ order.product_name }}</text>
          <text class="qty">数量 ×{{ order.quantity || 1 }}</text>
          <text class="cost">{{ order.points_cost }} 积分</text>
        </view>
      </view>
    </view>

    <view class="panel">
      <text class="sec-title">收货信息</text>
      <text class="line">{{ order.receiver_name }} {{ order.receiver_phone }}</text>
      <text class="line muted">{{ order.receiver_address || '—' }}</text>
    </view>

    <view class="panel" v-if="order.ship_company || order.ship_no">
      <text class="sec-title">物流信息</text>
      <text class="line">承运商 {{ order.ship_company || '—' }}</text>
      <text class="line muted">运单号 {{ order.ship_no || '—' }}</text>
    </view>

    <view class="panel">
      <text class="sec-title">订单信息</text>
      <text class="line muted">下单时间 {{ order.created_at || '—' }}</text>
      <text v-if="order.shipped_at" class="line muted">发货时间 {{ order.shipped_at }}</text>
      <text v-if="order.completed_at" class="line muted">完成时间 {{ order.completed_at }}</text>
      <text v-if="order.admin_remark" class="line muted">备注 {{ order.admin_remark }}</text>
    </view>
  </view>
  <view v-else class="loading">{{ loadErr || '加载中...' }}</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getPointsMallOrder } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/200/200'
const order = ref(null)
const loadErr = ref('')

function statusText(s) {
  const map = {
    pending: '待发货',
    shipped: '已发货',
    completed: '已完成',
    cancelled: '已取消',
  }
  return map[s] || s || '-'
}

async function load(id) {
  if (!isLoggedIn()) {
    uni.redirectTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent(`/pages/points-mall/order-detail?id=${id}`),
    })
    return
  }
  try {
    const res = await getPointsMallOrder(id)
    order.value = res?.data || res || null
    if (!order.value) loadErr.value = '订单不存在'
  } catch (e) {
    loadErr.value = e.message || '加载失败'
    order.value = null
  }
}

onLoad((q) => {
  load(Number(q.id) || 0)
})
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; min-height: 100vh; background: #f7f3ec; }
.status-card {
  background: linear-gradient(135deg, #e8d5b5, #c8a876);
  border-radius: 24rpx; padding: 36rpx 28rpx; color: #fff; margin-bottom: 20rpx;
}
.st { display: block; font-size: 36rpx; font-weight: 700; }
.no { display: block; margin-top: 8rpx; font-size: 24rpx; opacity: .9; }
.panel {
  background: #fff; border-radius: 24rpx; padding: 28rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.row { display: flex; gap: 20rpx; }
.img { width: 160rpx; height: 160rpx; border-radius: 16rpx; background: #f3f3f3; flex-shrink: 0; }
.body { flex: 1; min-width: 0; }
.name { display: block; font-size: 28rpx; font-weight: 600; }
.qty { display: block; margin-top: 10rpx; font-size: 24rpx; color: #71717a; }
.cost { display: block; margin-top: 12rpx; font-size: 30rpx; font-weight: 700; color: #b8860b; }
.sec-title { display: block; font-size: 28rpx; font-weight: 600; margin-bottom: 12rpx; }
.line { display: block; font-size: 26rpx; color: #18181b; line-height: 1.6; }
.line.muted { color: #71717a; }
.loading { text-align: center; padding: 120rpx 32rpx; color: #a1a1aa; }
</style>
