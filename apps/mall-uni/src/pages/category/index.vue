<template>
  <view class="page">
    <view class="layout">
      <scroll-view scroll-y class="side scroll-hide">
        <view
          class="side-item"
          :class="{ active: activeRoot === 0 }"
          @tap="selectRoot(0)"
        >全部</view>
        <view
          v-for="c in roots"
          :key="c.id"
          class="side-item"
          :class="{ active: activeRoot === c.id }"
          @tap="selectRoot(c.id)"
        >{{ c.name }}</view>
      </scroll-view>

      <view class="main">
        <view v-if="children.length" class="sub-grid">
          <view
            v-for="c in children"
            :key="c.id"
            class="sub-item"
            :class="{ on: activeCat === c.id }"
            @tap="selectChild(c.id)"
          >
            <image class="sub-icon" :src="c.icon || c.image || placeholder" mode="aspectFill" />
            <text class="sub-name line-1">{{ c.name }}</text>
          </view>
        </view>

        <scroll-view
          scroll-y
          class="goods-scroll scroll-hide"
          @scrolltolower="loadMore"
        >
          <view class="goods-grid">
            <view
              v-for="p in products"
              :key="p.id"
              class="goods-card"
              @tap="goDetail(p.id)"
            >
              <image class="goods-img" :src="p.main_image || placeholder" mode="aspectFill" />
              <text class="line-2 goods-name">{{ p.name }}</text>
              <view class="price-row">
                <text class="price">¥{{ p.sale_price }}</text>
                <text class="sold">月销{{ p.sold_count || 0 }}</text>
              </view>
            </view>
          </view>
          <view class="foot">
            <text v-if="loading">加载中...</text>
            <text v-else-if="!products.length">暂无商品</text>
            <text v-else-if="finished">没有更多了</text>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { listCategories, listProducts } from '../../api/index'

const placeholder = 'https://picsum.photos/id/96/200/200'
const flatCats = ref([])
const activeRoot = ref(0)
const activeCat = ref(0)
const products = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

const roots = computed(() => flatCats.value.filter((c) => !c.parent_id))

const children = computed(() => {
  if (!activeRoot.value) return []
  return flatCats.value.filter((c) => c.parent_id === activeRoot.value)
})

function goDetail(id) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
}

function selectRoot(id) {
  activeRoot.value = id
  activeCat.value = id
  loadProducts(true)
}

function selectChild(id) {
  activeCat.value = id
  loadProducts(true)
}

async function loadCats() {
  try {
    const res = await listCategories({ page: 1, page_size: 200 })
    flatCats.value = (res?.list || []).map((c) => ({
      ...c,
      parent_id: c.parent_id ?? c.ParentId ?? 0,
    }))
  } catch {
    flatCats.value = []
  }
}

async function loadProducts(reset = false) {
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      products.value = []
    }
    const params = {
      page: page.value,
      page_size: pageSize,
      order_by: 'sold_count_desc',
    }
    if (activeCat.value) params.category_id = activeCat.value
    const res = await listProducts(params)
    const rows = res?.list || []
    products.value = reset ? rows : products.value.concat(rows)
    const total = res?.total || 0
    if (products.value.length >= total || rows.length < pageSize) finished.value = true
    else page.value += 1
  } catch {
    if (reset) products.value = []
  } finally {
    loading.value = false
  }
}

function loadMore() {
  loadProducts(false)
}

function applyFocus() {
  const focus = Number(uni.getStorageSync('category_focus_id') || 0)
  uni.removeStorageSync('category_focus_id')
  if (!focus) {
    activeRoot.value = 0
    activeCat.value = 0
    return
  }
  const node = flatCats.value.find((c) => c.id === focus)
  if (!node) {
    activeRoot.value = 0
    activeCat.value = 0
    return
  }
  if (!node.parent_id) {
    activeRoot.value = node.id
    activeCat.value = node.id
  } else {
    activeRoot.value = node.parent_id
    activeCat.value = node.id
  }
}

onShow(async () => {
  if (!flatCats.value.length) await loadCats()
  applyFocus()
  loadProducts(true)
})
</script>

<style scoped>
.page { height: 100vh; background: #fafafa; }
.layout { display: flex; height: 100%; }
.side {
  width: 180rpx; background: #f4f4f5; height: 100%;
}
.side-item {
  padding: 28rpx 16rpx; font-size: 24rpx; color: #52525b; text-align: center;
  position: relative;
}
.side-item.active {
  background: #fff; color: #c8a876; font-weight: 600;
}
.side-item.active::before {
  content: ''; position: absolute; left: 0; top: 24rpx; bottom: 24rpx; width: 6rpx;
  background: #c8a876; border-radius: 0 6rpx 6rpx 0;
}
.main { flex: 1; min-width: 0; display: flex; flex-direction: column; background: #fff; }
.sub-grid {
  display: flex; flex-wrap: wrap; padding: 16rpx 12rpx 8rpx; gap: 8rpx;
  border-bottom: 1rpx solid #f4f4f5;
}
.sub-item {
  width: calc(33.33% - 8rpx); display: flex; flex-direction: column; align-items: center;
  padding: 12rpx 4rpx; border-radius: 12rpx;
}
.sub-item.on { background: rgba(200,168,118,.12); }
.sub-icon { width: 72rpx; height: 72rpx; border-radius: 16rpx; background: #f4f4f5; }
.sub-name { font-size: 22rpx; color: #3f3f46; margin-top: 8rpx; max-width: 100%; }
.goods-scroll { flex: 1; height: 0; padding: 16rpx; box-sizing: border-box; }
.goods-grid { display: flex; flex-wrap: wrap; gap: 16rpx; }
.goods-card {
  width: calc(50% - 8rpx); background: #fafafa; border-radius: 16rpx; padding: 12rpx;
  box-sizing: border-box;
}
.goods-img { width: 100%; height: 220rpx; border-radius: 12rpx; background: #eee; }
.goods-name { display: block; font-size: 24rpx; margin-top: 10rpx; color: #18181b; min-height: 68rpx; }
.price-row { display: flex; align-items: baseline; justify-content: space-between; margin-top: 8rpx; }
.price { color: #d83636; font-weight: 700; font-size: 28rpx; }
.sold { font-size: 20rpx; color: #a1a1aa; }
.foot { text-align: center; color: #a1a1aa; font-size: 22rpx; padding: 24rpx 0 40rpx; }
.line-1, .line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; overflow: hidden;
}
.line-1 { -webkit-line-clamp: 1; }
.line-2 { -webkit-line-clamp: 2; }
.scroll-hide::-webkit-scrollbar { display: none; width: 0; height: 0; }
</style>
