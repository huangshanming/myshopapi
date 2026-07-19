<template>
  <view class="page">
    <view class="tabs">
      <text v-for="t in tabs" :key="t.v" class="tab" :class="{ on: status === t.v }" @tap="switchTab(t.v)">{{ t.l }}</text>
    </view>
    <view v-for="uc in list" :key="uc.id" class="card" :class="{ dim: status !== 'unused' }">
      <view class="left">
        <text class="amt">{{ face(uc.coupon) }}</text>
        <text class="rule">{{ rule(uc.coupon) }}</text>
      </view>
      <view class="right">
        <text class="name">{{ uc.coupon_name || uc.coupon?.name }}</text>
        <text class="sub">有效期至 {{ uc.valid_end }}</text>
        <text v-if="status === 'expired'" class="tag">已失效</text>
        <button v-if="status === 'unused'" class="btn" @tap="goUse">去使用</button>
      </view>
    </view>
    <view v-if="!list.length" class="empty">暂无优惠券</view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { listMyCoupons } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const tabs = [
  { v: 'unused', l: '未使用' },
  { v: 'used', l: '已使用' },
  { v: 'expired', l: '已过期' },
]
const status = ref('unused')
const list = ref([])

function face(c) {
  if (!c) return '-'
  if (c.coupon_type === 'discount') return `${(Number(c.discount_rate) * 10).toFixed(1)}折`
  return `¥${c.discount_amount}`
}
function rule(c) {
  if (!c) return ''
  if (c.coupon_type === 'no_threshold') return '无门槛'
  return `满${c.threshold_amount}可用`
}

async function load() {
  if (!isLoggedIn()) {
    uni.redirectTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/coupon/mine') })
    return
  }
  try {
    const res = await listMyCoupons({ page: 1, page_size: 50, status: status.value })
    list.value = res?.list || []
  } catch {
    list.value = []
  }
}

function switchTab(v) {
  status.value = v
  load()
}

function goUse() {
  uni.switchTab({ url: '/pages/index/index' })
}

onShow(load)
</script>

<style scoped>
.page { min-height: 100vh; background: #f7f3ec; padding-bottom: 40rpx; }
.tabs { display: flex; background: #fff; padding: 0 12rpx; }
.tab { flex: 1; text-align: center; padding: 24rpx 0; color: #64748b; font-size: 28rpx; }
.tab.on { color: #c4894a; font-weight: 600; border-bottom: 4rpx solid #c4894a; }
.card {
  display: flex; margin: 20rpx 24rpx 0; background: #fff; border-radius: 16rpx; overflow: hidden;
}
.card.dim { opacity: .55; }
.left {
  width: 200rpx; background: linear-gradient(160deg, #d4a574, #c4894a);
  color: #fff; display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 24rpx;
}
.amt { font-size: 40rpx; font-weight: 700; }
.rule { font-size: 22rpx; margin-top: 8rpx; }
.right { flex: 1; padding: 24rpx; }
.name { font-weight: 600; font-size: 28rpx; display: block; }
.sub { font-size: 22rpx; color: #94a3b8; margin: 8rpx 0; display: block; }
.tag { font-size: 22rpx; color: #ef4444; }
.btn {
  margin: 8rpx 0 0; display: inline-block; height: 52rpx; line-height: 52rpx; padding: 0 24rpx;
  font-size: 24rpx; color: #fff; background: #c4894a; border-radius: 26rpx; border: none;
}
.empty { text-align: center; color: #94a3b8; padding: 80rpx; }
</style>
