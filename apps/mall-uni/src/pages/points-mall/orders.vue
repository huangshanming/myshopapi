<template>
  <view class="page">
    <scroll-view class="tabs" scroll-x :show-scrollbar="false">
      <view
        v-for="t in tabs"
        :key="t.v"
        class="tab"
        :class="{ on: status === t.v }"
        @tap="onTab(t.v)"
      >{{ t.l }}</view>
    </scroll-view>

    <view
      v-for="o in displayList"
      :key="o.id"
      class="card"
      @tap="goDetail(o)"
    >
      <view class="top">
        <text class="no">单号 {{ o.order_no }}</text>
        <text class="st">{{ statusText(o.status) }}</text>
      </view>
      <view class="row">
        <image class="img" :src="o.product_cover || placeholder" mode="aspectFill" />
        <view class="body">
          <text class="name line-2">{{ o.product_name }}</text>
          <text class="cost">{{ o.points_cost }} 积分</text>
          <text class="time">{{ o.created_at || '' }}</text>
        </view>
      </view>
    </view>

    <view v-if="!loading && !displayList.length" class="empty">暂无兑换记录</view>
    <view class="foot">
      <text v-if="loading">加载中...</text>
      <text v-else-if="finished && list.length">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onReachBottom, onShow } from '@dcloudio/uni-app'
import { listPointsMallOrders } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/200/200'
const list = ref([])
const status = ref('')
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

const tabs = [
  { v: '', l: '全部' },
  { v: 'pending', l: '待发货' },
  { v: 'shipped', l: '已发货' },
  { v: 'completed', l: '已完成' },
  { v: 'cancelled', l: '已取消' },
]

const displayList = computed(() => {
  if (!status.value) return list.value
  return list.value.filter((o) => o.status === status.value)
})

function statusText(s) {
  const map = {
    pending: '待发货',
    shipped: '已发货',
    completed: '已完成',
    cancelled: '已取消',
  }
  return map[s] || s || '-'
}

async function load(reset = false) {
  if (!isLoggedIn()) {
    uni.redirectTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/points-mall/orders'),
    })
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
    const res = await listPointsMallOrders({ page: page.value, page_size: pageSize })
    const rows = res?.list || []
    list.value = reset ? rows : list.value.concat(rows)
    const total = Number(res?.total) || 0
    if (list.value.length >= total || rows.length < pageSize) {
      finished.value = true
    } else {
      page.value += 1
    }
  } catch {
    if (reset) list.value = []
  } finally {
    loading.value = false
  }
}

function onTab(v) {
  status.value = v
}

function goDetail(o) {
  uni.navigateTo({ url: `/pages/points-mall/order-detail?id=${o.id}` })
}

onReachBottom(() => load(false))
onShow(() => load(true))
</script>

<style scoped>
.page { padding: 0 32rpx 48rpx; min-height: 100vh; background: #f7f3ec; }
.tabs { white-space: nowrap; padding: 24rpx 0 8rpx; }
.tab {
  display: inline-block; margin-right: 16rpx; padding: 12rpx 28rpx;
  border-radius: 999rpx; font-size: 24rpx; color: #8b7355; background: #fff;
}
.tab.on { background: #c8a876; color: #fff; }
.card {
  background: #fff; border-radius: 24rpx; padding: 24rpx; margin-top: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.top { display: flex; justify-content: space-between; margin-bottom: 16rpx; }
.no { font-size: 22rpx; color: #a1a1aa; }
.st { font-size: 24rpx; color: #c8a876; font-weight: 600; }
.row { display: flex; gap: 20rpx; }
.img { width: 160rpx; height: 160rpx; border-radius: 16rpx; background: #f3f3f3; flex-shrink: 0; }
.body { flex: 1; min-width: 0; display: flex; flex-direction: column; }
.name { font-size: 28rpx; font-weight: 600; color: #18181b; }
.cost { margin-top: 12rpx; font-size: 28rpx; font-weight: 700; color: #b8860b; }
.time { margin-top: auto; font-size: 22rpx; color: #a1a1aa; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2; overflow: hidden;
}
.empty { text-align: center; color: #a1a1aa; padding: 100rpx 0; font-size: 26rpx; }
.foot { text-align: center; padding: 24rpx; color: #a1a1aa; font-size: 24rpx; }
</style>
