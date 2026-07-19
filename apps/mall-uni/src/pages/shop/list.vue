<template>
  <view class="page">
    <view class="tabs">
      <text
        class="tab"
        :class="{ active: slotType === 'brand_shop' }"
        @tap="switchType('brand_shop')"
      >品牌商户</text>
      <text
        class="tab"
        :class="{ active: slotType === 'quality_shop' }"
        @tap="switchType('quality_shop')"
      >优质商户</text>
    </view>

    <view
      v-for="s in list"
      :key="s.id"
      class="card"
      @tap="goDetail(s.id)"
    >
      <image class="cover" :src="s.storefront_image || s.logo || placeholder" mode="aspectFill" />
      <view class="body">
        <view class="row">
          <text class="name line-1">{{ s.name }}</text>
          <text v-if="s.paid" class="badge">推广</text>
        </view>
        <text class="sub">{{ s.category || '优选商户' }}{{ s.city ? ` · ${s.city}` : '' }}</text>
        <text class="desc line-2">{{ s.description || '品质好店，欢迎逛逛' }}</text>
      </view>
    </view>

    <view class="foot">
      <text v-if="loading">加载中...</text>
      <text v-else-if="!list.length">暂无店铺</text>
      <text v-else-if="finished">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onReachBottom } from '@dcloudio/uni-app'
import { listShops } from '../../api/index'

const placeholder = 'https://picsum.photos/id/96/400/400'
const slotType = ref('quality_shop')
const list = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

function goDetail(id) {
  uni.navigateTo({ url: `/pages/shop/detail?id=${id}` })
}

function switchType(t) {
  if (slotType.value === t) return
  slotType.value = t
  uni.setNavigationBarTitle({
    title: t === 'brand_shop' ? '品牌商户' : '优质商户',
  })
  load(true)
}

async function load(reset = false) {
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      list.value = []
    }
    const res = await listShops({
      page: page.value,
      page_size: pageSize,
      slot_type: slotType.value,
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

onLoad((q) => {
  const t = q?.slot_type === 'brand_shop' ? 'brand_shop' : 'quality_shop'
  slotType.value = t
  uni.setNavigationBarTitle({
    title: t === 'brand_shop' ? '品牌商户' : '优质商户',
  })
  load(true)
})

onReachBottom(() => load(false))
</script>

<style scoped>
.page { padding: 16rpx 24rpx 40rpx; }
.tabs {
  display: flex; gap: 24rpx; margin-bottom: 20rpx;
}
.tab {
  font-size: 28rpx; color: #71717a; padding: 8rpx 4rpx;
}
.tab.active { color: #c8a876; font-weight: 600; border-bottom: 4rpx solid #c8a876; }
.card {
  display: flex; gap: 20rpx; background: #fff; border-radius: 16rpx;
  padding: 20rpx; margin-bottom: 16rpx;
}
.cover { width: 160rpx; height: 160rpx; border-radius: 12rpx; flex-shrink: 0; background: #f4f4f5; }
.body { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 8rpx; }
.row { display: flex; align-items: center; gap: 12rpx; }
.name { flex: 1; font-size: 30rpx; font-weight: 600; color: #18181b; }
.badge {
  font-size: 20rpx; color: #c8a876; background: rgba(200,168,118,.15);
  padding: 2rpx 10rpx; border-radius: 6rpx;
}
.sub { font-size: 22rpx; color: #a1a1aa; }
.desc { font-size: 24rpx; color: #52525b; }
.foot { text-align: center; color: #a1a1aa; font-size: 24rpx; padding: 24rpx 0; }
.line-1 { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2;
  overflow: hidden;
}
</style>
