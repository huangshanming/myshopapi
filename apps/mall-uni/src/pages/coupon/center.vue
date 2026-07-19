<template>
  <view class="page">
    <view v-for="c in list" :key="c.id" class="card" :class="{ dim: c.claimed_by_me || c.display_status === 'sold_out' }">
      <view class="left">
        <text class="amt">{{ face(c) }}</text>
        <text class="rule">{{ rule(c) }}</text>
      </view>
      <view class="right">
        <text class="name">{{ c.name }}</text>
        <text class="sub">{{ c.display_status === 'sold_out' ? '已领完' : (c.claimed_by_me ? '已领取' : '可领取') }}</text>
        <button
          class="btn"
          :disabled="c.claimed_by_me || c.display_status === 'sold_out'"
          @tap="claim(c)"
        >{{ c.claimed_by_me ? '已领取' : '立即领取' }}</button>
      </view>
    </view>
    <view v-if="!list.length" class="empty">暂无可领优惠券</view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { claimCoupon, listCouponCenter } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const list = ref([])

function face(c) {
  if (c.coupon_type === 'discount') return `${(Number(c.discount_rate) * 10).toFixed(1)}折`
  return `¥${c.discount_amount}`
}
function rule(c) {
  if (c.coupon_type === 'no_threshold') return '无门槛'
  if (c.coupon_type === 'discount') return c.threshold_amount > 0 ? `满${c.threshold_amount}可用` : '无门槛折扣'
  return `满${c.threshold_amount}减${c.discount_amount}`
}

async function load() {
  try {
    const res = await listCouponCenter()
    list.value = res?.list || []
  } catch {
    list.value = []
  }
}

async function claim(c) {
  if (!isLoggedIn()) {
    uni.navigateTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/coupon/center') })
    return
  }
  try {
    await claimCoupon(c.id, 'direct')
    uni.showToast({ title: '领取成功', icon: 'success' })
    load()
  } catch { /* toast */ }
}

onShow(load)
</script>

<style scoped>
.page { padding: 24rpx; min-height: 100vh; background: #f7f3ec; }
.card {
  display: flex; background: #fff; border-radius: 16rpx; margin-bottom: 20rpx; overflow: hidden;
  box-shadow: 0 4rpx 16rpx rgba(180,140,80,.08);
}
.card.dim { opacity: .55; }
.left {
  width: 200rpx; background: linear-gradient(160deg, #d4a574, #c4894a);
  color: #fff; display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 24rpx 12rpx;
}
.amt { font-size: 40rpx; font-weight: 700; }
.rule { font-size: 22rpx; margin-top: 8rpx; opacity: .9; }
.right { flex: 1; padding: 24rpx; }
.name { font-size: 30rpx; font-weight: 600; display: block; }
.sub { font-size: 22rpx; color: #94a3b8; margin: 8rpx 0 16rpx; display: block; }
.btn {
  display: inline-block; margin: 0; padding: 0 28rpx; height: 56rpx; line-height: 56rpx;
  font-size: 24rpx; color: #fff; background: #c4894a; border-radius: 28rpx; border: none;
}
.btn[disabled] { background: #cbd5e1; }
.empty { text-align: center; color: #94a3b8; padding: 80rpx; }
</style>
