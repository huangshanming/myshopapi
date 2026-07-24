<template>
  <view class="page">
    <view class="search-bar">
      <input
        class="search-input"
        v-model="keyword"
        confirm-type="search"
        placeholder="搜索本地商家"
        @confirm="reload"
      />
      <text class="search-btn" @tap="reload">搜索</text>
    </view>
    <view class="sort-row">
      <text class="sort" :class="{ on: sort === 'distance' }" @tap="setSort('distance')">距离优先</text>
      <text class="sort" :class="{ on: sort === 'default' }" @tap="setSort('default')">默认</text>
      <text class="loc" @tap="goLocate">📍 {{ locLabel }}</text>
    </view>

    <view v-for="s in list" :key="s.id" class="card" @tap="goDetail(s.id)">
      <image class="cover" :src="s.storefront_image || s.logo || placeholder" mode="aspectFill" />
      <view class="body">
        <view class="row">
          <text class="name line-1">{{ s.name }}</text>
          <text v-if="s.distance_km != null" class="dist">{{ formatDist(s.distance_km) }}</text>
        </view>
        <text class="sub">{{ s.category || '本地商家' }}{{ s.city ? ` · ${s.city}` : '' }}</text>
        <text class="desc line-2">{{ s.full_address || s.address || '暂无地址' }}</text>
      </view>
    </view>

    <view class="foot">
      <text v-if="!hasLoc">请先选择定位后再查看附近商家</text>
      <text v-else-if="loading">加载中...</text>
      <text v-else-if="!list.length">附近暂无本地商家</text>
      <text v-else-if="finished">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onReachBottom, onShow } from '@dcloudio/uni-app'
import { listLocalShops } from '../../api/index'
import { getCity, getCoords, hasCoords } from '../../stores/location'

const placeholder = 'https://picsum.photos/id/96/400/400'
const keyword = ref('')
const sort = ref('distance')
const list = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)
const hasLoc = ref(hasCoords())
const locLabel = computed(() => getCity() || '选择定位')

function formatDist(km) {
  if (km == null || Number.isNaN(Number(km))) return '—'
  const n = Number(km)
  if (n < 1) return `${Math.round(n * 1000)}m`
  return `${n.toFixed(1)}km`
}

function goDetail(id) {
  uni.navigateTo({ url: `/pages/shop/detail?id=${id}` })
}

function goLocate() {
  uni.navigateTo({ url: '/pages/location/city' })
}

function setSort(s) {
  if (sort.value === s) return
  sort.value = s
  reload()
}

function reload() {
  load(true)
}

async function load(reset = false) {
  hasLoc.value = hasCoords()
  if (!hasLoc.value) {
    list.value = []
    finished.value = true
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
    const coords = getCoords()
    const res = await listLocalShops({
      page: page.value,
      page_size: pageSize,
      lat: coords.latitude,
      lng: coords.longitude,
      keyword: keyword.value.trim(),
      sort: sort.value,
    })
    const rows = res?.list || []
    list.value = reset ? rows : list.value.concat(rows)
    const total = res?.total || 0
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

onLoad(() => load(true))
onShow(() => {
  hasLoc.value = hasCoords()
  load(true)
})
onReachBottom(() => load(false))
</script>

<style scoped>
.page { padding: 16rpx 24rpx 40rpx; }
.search-bar {
  display: flex; gap: 16rpx; align-items: center; margin-bottom: 16rpx;
}
.search-input {
  flex: 1; background: #fff; border-radius: 999rpx; padding: 16rpx 28rpx; font-size: 26rpx;
}
.search-btn {
  color: #c8a876; font-size: 28rpx; font-weight: 600; padding: 8rpx 12rpx;
}
.sort-row {
  display: flex; align-items: center; gap: 24rpx; margin-bottom: 20rpx;
}
.sort { font-size: 26rpx; color: #71717a; }
.sort.on { color: #c8a876; font-weight: 600; }
.loc { margin-left: auto; font-size: 24rpx; color: #a1a1aa; }
.card {
  display: flex; gap: 20rpx; background: #fff; border-radius: 16rpx;
  padding: 20rpx; margin-bottom: 16rpx;
}
.cover { width: 160rpx; height: 160rpx; border-radius: 12rpx; flex-shrink: 0; background: #f4f4f5; }
.body { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 8rpx; }
.row { display: flex; align-items: center; gap: 12rpx; }
.name { flex: 1; font-size: 30rpx; font-weight: 600; color: #18181b; }
.dist { font-size: 24rpx; color: #c8a876; font-weight: 600; }
.sub { font-size: 22rpx; color: #a1a1aa; }
.desc { font-size: 24rpx; color: #52525b; }
.foot { text-align: center; color: #a1a1aa; font-size: 24rpx; padding: 24rpx 0; }
.line-1 { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2;
  overflow: hidden;
}
</style>
