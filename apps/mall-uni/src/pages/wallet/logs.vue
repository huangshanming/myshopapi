<template>
  <view class="page">
    <view v-if="!list.length && !loading" class="empty">暂无流水</view>
    <view v-for="row in list" :key="row.id" class="card">
      <view class="row">
        <text class="type">{{ typeText(row.change_type) }}</text>
        <text class="amount" :class="{ plus: row.amount > 0 }">{{ formatAmount(row.amount) }}</text>
      </view>
      <view class="meta">
        <text>余额 ¥{{ row.balance_after }} · 冻结 ¥{{ row.frozen_after }}</text>
      </view>
      <view class="meta">
        <text>{{ row.remark || '—' }}</text>
        <text>{{ row.created_at }}</text>
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
import { listUserWalletLogs } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const TYPE_MAP = {
  admin_adjust: '后台调账',
  order_freeze: '下单冻结',
  order_unfreeze: '取消解冻',
  order_settle: '订单实扣',
}

const list = ref([])
const page = ref(1)
const loading = ref(false)
const finished = ref(false)

function typeText(t) {
  return TYPE_MAP[t] || t
}

function formatAmount(n) {
  const v = Number(n) || 0
  return (v > 0 ? '+' : '') + '¥' + v
}

async function load(reset = false) {
  if (!isLoggedIn()) {
    uni.redirectTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/wallet/logs') })
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
    const res = await listUserWalletLogs({ page: page.value, page_size: 15 })
    const rows = res.data?.list || []
    list.value = reset ? rows : list.value.concat(rows)
    const total = res.data?.total || 0
    if (list.value.length >= total || rows.length < 15) finished.value = true
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
.empty { text-align: center; color: #a1a1aa; padding: 80rpx 0; }
.card {
  background: #fff; border-radius: 24rpx; padding: 24rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.row { display: flex; justify-content: space-between; margin-bottom: 12rpx; }
.type { font-size: 28rpx; font-weight: 600; }
.amount { font-size: 30rpx; font-weight: 600; color: #b45309; }
.amount.plus { color: #15803d; }
.meta {
  display: flex; justify-content: space-between; gap: 16rpx;
  color: #71717a; font-size: 22rpx; margin-top: 8rpx;
}
.more { text-align: center; color: #a1a1aa; padding: 24rpx; font-size: 24rpx; }
</style>
