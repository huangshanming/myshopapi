<template>
  <view class="page">
    <view class="head">
      <view class="points-row">
        <text class="points-label">我的积分</text>
        <text class="points-num">{{ points }}</text>
      </view>
      <view class="search-row">
        <input
          class="search"
          v-model="keyword"
          confirm-type="search"
          placeholder="搜索积分好物"
          @confirm="onSearch"
        />
        <text class="search-btn" @tap="onSearch">搜索</text>
      </view>
      <scroll-view class="chips" scroll-x :show-scrollbar="false">
        <view
          v-for="t in tiers"
          :key="t.key"
          class="chip"
          :class="{ on: tier === t.key }"
          @tap="onTier(t.key)"
        >{{ t.label }}</view>
      </scroll-view>
    </view>

    <view class="grid">
      <view
        v-for="p in displayItems"
        :key="p.id"
        class="card"
        @tap="goDetail(p)"
      >
        <image class="cover" :src="p.cover_url || placeholder" mode="aspectFill" />
        <view class="body">
          <text class="name line-2">{{ p.name }}</text>
          <view class="meta">
            <text class="price">{{ p.points_price }}</text>
            <text class="unit">积分</text>
          </view>
          <text class="stock">库存 {{ p.stock ?? 0 }}</text>
        </view>
      </view>
    </view>

    <view v-if="!loading && !displayItems.length" class="empty">暂无积分商品</view>
    <view class="foot">
      <text v-if="loading">加载中...</text>
      <text v-else-if="finished">没有更多了</text>
      <text v-else @tap="loadMore">上拉加载更多</text>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onReachBottom, onShow } from '@dcloudio/uni-app'
import { getUserPoints, listPointsMallProducts } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/400/400'
const items = ref([])
const points = ref(0)
const keyword = ref('')
const appliedKeyword = ref('')
const tier = ref('all')
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

const tiers = [
  { key: 'all', label: '全部' },
  { key: 'low', label: '100以内', min: 0, max: 100 },
  { key: 'mid', label: '101-500', min: 101, max: 500 },
  { key: 'high', label: '501-2000', min: 501, max: 2000 },
  { key: 'top', label: '2000以上', min: 2001, max: Infinity },
]

const displayItems = computed(() => {
  const t = tiers.find((x) => x.key === tier.value)
  if (!t || t.key === 'all') return items.value
  return items.value.filter((p) => {
    const price = Number(p.points_price) || 0
    return price >= t.min && price <= t.max
  })
})

async function loadPoints() {
  if (!isLoggedIn()) {
    points.value = 0
    return
  }
  try {
    const res = await getUserPoints()
    points.value = Number(res?.points) || 0
  } catch {
    points.value = 0
  }
}

async function load(reset = false) {
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      items.value = []
    }
    const res = await listPointsMallProducts({
      page: page.value,
      page_size: pageSize,
      status: 'on',
      keyword: appliedKeyword.value || undefined,
    })
    const list = res?.list || []
    items.value = reset ? list : items.value.concat(list)
    const total = Number(res?.total) || 0
    if (items.value.length >= total || list.length < pageSize) {
      finished.value = true
    } else {
      page.value += 1
    }
  } catch {
    if (reset) items.value = []
  } finally {
    loading.value = false
  }
}

function loadMore() {
  load(false)
}

function onSearch() {
  appliedKeyword.value = String(keyword.value || '').trim()
  load(true)
}

function onTier(key) {
  if (tier.value === key) return
  tier.value = key
}

function goDetail(p) {
  uni.navigateTo({ url: `/pages/points-mall/detail?id=${p.id}` })
}

onReachBottom(() => loadMore())

onShow(() => {
  loadPoints()
  if (!items.value.length) load(true)
})
</script>

<style scoped>
.page { padding: 0 0 48rpx; min-height: 100vh; background: #f7f3ec; }
.head {
  padding: 24rpx 32rpx 8rpx;
  background: linear-gradient(180deg, #efe4d2 0%, #f7f3ec 100%);
}
.points-row { display: flex; align-items: baseline; gap: 12rpx; margin-bottom: 20rpx; }
.points-label { font-size: 24rpx; color: #8b7355; }
.points-num { font-size: 44rpx; font-weight: 700; color: #5c4a32; }
.search-row { display: flex; gap: 16rpx; align-items: center; margin-bottom: 20rpx; }
.search {
  flex: 1; height: 72rpx; background: #fff; border-radius: 999rpx;
  padding: 0 28rpx; font-size: 26rpx;
}
.search-btn {
  font-size: 26rpx; color: #fff; background: #c8a876;
  padding: 16rpx 28rpx; border-radius: 999rpx;
}
.chips { white-space: nowrap; margin-bottom: 8rpx; }
.chip {
  display: inline-block; margin-right: 16rpx; padding: 10rpx 24rpx;
  border-radius: 999rpx; font-size: 24rpx; color: #8b7355;
  background: rgba(255,255,255,.7);
}
.chip.on { background: #c8a876; color: #fff; }
.grid {
  display: flex; flex-wrap: wrap; gap: 20rpx;
  padding: 16rpx 32rpx 0;
}
.card {
  width: calc(50% - 10rpx); background: #fff; border-radius: 24rpx; overflow: hidden;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.cover { width: 100%; height: 280rpx; background: #f3f3f3; display: block; }
.body { padding: 16rpx 18rpx 20rpx; }
.name { font-size: 26rpx; font-weight: 600; color: #18181b; min-height: 72rpx; }
.meta { margin-top: 12rpx; display: flex; align-items: baseline; gap: 6rpx; }
.price { font-size: 34rpx; font-weight: 700; color: #b8860b; }
.unit { font-size: 22rpx; color: #8b7355; }
.stock { display: block; margin-top: 6rpx; font-size: 22rpx; color: #a1a1aa; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2; overflow: hidden;
}
.empty { text-align: center; color: #a1a1aa; padding: 80rpx 0 24rpx; font-size: 26rpx; }
.foot { text-align: center; padding: 24rpx; color: #a1a1aa; font-size: 24rpx; }
</style>
