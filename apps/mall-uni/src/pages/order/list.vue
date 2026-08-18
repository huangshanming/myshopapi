<template>
  <view class="page">
    <scroll-view class="tabs" scroll-x :show-scrollbar="false">
      <view
        v-for="t in tabs"
        :key="t.v"
        class="tab"
        :class="{ on: status === t.v }"
        @tap="switchTab(t.v)"
      >{{ t.l }}</view>
    </scroll-view>

    <view v-if="status === 'after_sale'">
      <view v-if="!afterList.length && !loading" class="empty">暂无售后记录</view>
      <view v-for="as in afterList" :key="as.id" class="card" @tap="goDetail(as.order_id)">
        <view class="row">
          <text class="no">{{ as.order_no }}</text>
          <text class="status">{{ afterStatusText(as.status) }}</text>
        </view>
        <text class="line-1">{{ as.type === 'return_refund' ? '退货退款' : '仅退款' }} · ¥{{ as.amount }}</text>
        <text class="sub">{{ as.reason || '—' }}</text>
        <view class="foot">
          <text class="sub">{{ as.created_at }}</text>
          <text class="link">查看订单 ›</text>
        </view>
      </view>
    </view>

    <view v-else>
      <view v-if="!list.length && !loading" class="empty">暂无订单</view>
      <view v-for="o in list" :key="o.id" class="card" @tap="goDetail(o.id)">
        <view class="row">
          <text class="no">{{ o.order_no }}</text>
          <text class="status">{{ statusText(o.status) }}</text>
        </view>
        <view v-for="it in o.items || []" :key="it.id" class="item">
          <image class="cover" :src="itemImage(it) || placeholder" mode="aspectFill" />
          <view class="item-meta">
            <text class="line-1">{{ it.product_name }}</text>
            <text class="sub">×{{ it.quantity }} · ¥{{ it.price }}</text>
          </view>
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
    </view>

    <view class="more">
      <text v-if="loading">加载中...</text>
      <text v-else-if="finished && (status === 'after_sale' ? afterList.length : list.length)">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onReachBottom, onShow } from '@dcloudio/uni-app'
import {
  AFTER_SALE_STATUS, ORDER_STATUS, getProductDetail, listMyAfterSales, listOrders,
} from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const tabs = [
  { v: '', l: '全部' },
  { v: 'pending', l: '待付款' },
  { v: 'confirmed', l: '待发货' },
  { v: 'shipped', l: '待收货' },
  { v: 'completed', l: '待评价' },
  { v: 'reviewed', l: '已评价' },
  { v: 'after_sale', l: '退款/售后' },
]

const placeholder = 'https://picsum.photos/id/96/200/200'
const status = ref('')
const list = ref([])
const afterList = ref([])
const page = ref(1)
const loading = ref(false)
const finished = ref(false)
const coverCache = ref({})

function statusText(s) {
  return ORDER_STATUS[s] || s
}

function afterStatusText(s) {
  return AFTER_SALE_STATUS[s] || s
}

function parseSnap(snap) {
  if (!snap) return null
  if (typeof snap === 'object') return snap
  try {
    return JSON.parse(snap)
  } catch {
    return null
  }
}

function itemImage(it) {
  if (!it) return ''
  if (it._cover) return it._cover
  const snap = parseSnap(it.sku_snapshot)
  const fromSnap = snap?.image || snap?.cover || snap?.main_image || ''
  if (fromSnap) return fromSnap
  return coverCache.value[it.product_id] || ''
}

async function enrichCovers(orders) {
  const need = []
  const seen = new Set()
  for (const o of orders || []) {
    for (const it of o.items || []) {
      if (itemImage(it) || !it.product_id || seen.has(it.product_id)) continue
      if (coverCache.value[it.product_id]) continue
      seen.add(it.product_id)
      need.push(it.product_id)
    }
  }
  if (!need.length) return
  await Promise.all(
    need.slice(0, 20).map(async (id) => {
      try {
        const p = await getProductDetail(id)
        const url = p?.main_image || ''
        if (url) coverCache.value = { ...coverCache.value, [id]: url }
      } catch {
        /* ignore */
      }
    }),
  )
  // trigger view update after cache fill
  list.value = list.value.map((o) => ({
    ...o,
    items: (o.items || []).map((it) => ({
      ...it,
      _cover: itemImage(it),
    })),
  }))
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

function switchTab(v) {
  if (status.value === v) return
  status.value = v
  load(true)
}

async function load(reset = false) {
  if (!isLoggedIn()) {
    uni.redirectTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/order/list' + (status.value ? `?status=${status.value}` : '')),
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
      afterList.value = []
    }
    if (status.value === 'after_sale') {
      const res = await listMyAfterSales({ page: page.value, page_size: 10 })
      const rows = res?.list || []
      afterList.value = reset ? rows : afterList.value.concat(rows)
      const total = res?.total || 0
      if (afterList.value.length >= total || rows.length < 10) finished.value = true
      else page.value += 1
    } else {
      const params = { page: page.value, page_size: 10 }
      if (status.value) params.status = status.value
      const res = await listOrders(params)
      const rows = res?.list || []
      list.value = reset ? rows : list.value.concat(rows)
      const total = res?.total || 0
      if (list.value.length >= total || rows.length < 10) finished.value = true
      else page.value += 1
      await enrichCovers(rows)
    }
  } catch {
    if (reset) {
      list.value = []
      afterList.value = []
    }
  } finally {
    loading.value = false
  }
}

onLoad((q) => {
  status.value = q?.status || ''
})

onShow(() => load(true))
onReachBottom(() => load(false))
</script>

<style scoped>
.page { padding-bottom: 24rpx; min-height: 100vh; background: #f7f3ec; }
.tabs {
  white-space: nowrap; background: #fff; padding: 0 8rpx;
  border-bottom: 1rpx solid #f0ebe3; position: sticky; top: 0; z-index: 2;
}
.tab {
  display: inline-block; padding: 24rpx 28rpx; font-size: 28rpx; color: #78716c;
}
.tab.on { color: #c4894a; font-weight: 600; border-bottom: 4rpx solid #c4894a; }
.card {
  background: #fff; border-radius: 24rpx; padding: 24rpx; margin: 20rpx 24rpx 0;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.row { display: flex; justify-content: space-between; margin-bottom: 16rpx; }
.no { font-size: 24rpx; color: #71717a; }
.status { color: #c8a876; font-size: 24rpx; }
.item {
  display: flex; align-items: center; gap: 16rpx;
  margin-bottom: 12rpx; padding-bottom: 12rpx;
  border-bottom: 1rpx solid #f5f5f5;
}
.item:last-of-type { border-bottom: none; margin-bottom: 0; padding-bottom: 0; }
.cover {
  width: 120rpx; height: 120rpx; border-radius: 12rpx;
  background: #f4f4f5; flex-shrink: 0;
}
.item-meta { flex: 1; min-width: 0; }
.line-1 {
  font-size: 28rpx; display: block; color: #18181b;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.sub { color: #71717a; font-size: 22rpx; display: block; margin-top: 6rpx; }
.foot { display: flex; justify-content: space-between; align-items: center; margin-top: 16rpx; }
.price { color: #d83636; font-weight: 700; font-size: 30rpx; }
.link { font-size: 24rpx; color: #c4894a; }
.actions { display: flex; justify-content: flex-end; gap: 24rpx; margin-top: 16rpx; padding-top: 16rpx; border-top: 1rpx solid #f5f5f5; }
.act { color: #c8a876; font-size: 24rpx; }
.act.muted { color: #71717a; }
.empty, .more { text-align: center; color: #71717a; padding: 40rpx; font-size: 24rpx; }
</style>
