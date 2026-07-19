<template>
  <view class="page">
    <image class="hero" :src="shop.storefront_image || shop.logo || placeholder" mode="aspectFill" />
    <view class="panel">
      <view class="row">
        <image class="logo" :src="shop.logo || placeholder" mode="aspectFill" />
        <view class="meta">
          <text class="name">{{ shop.name || '店铺' }}</text>
          <text class="sub">{{ shop.category || '优选商户' }}{{ shop.city ? ` · ${shop.city}` : '' }}</text>
        </view>
      </view>
      <text class="desc">{{ shop.description || '暂无店铺介绍' }}</text>
    </view>

    <view class="section">
      <text class="sec-title">店铺商品</text>
      <view v-if="!products.length && !loading" class="empty">暂无商品</view>
      <view class="grid">
        <view
          v-for="p in products"
          :key="p.id"
          class="goods"
          @tap="goProduct(p.id)"
        >
          <image class="g-img" :src="p.main_image || placeholder" mode="aspectFill" />
          <text class="g-name line-2">{{ p.name }}</text>
          <text class="price">¥{{ p.sale_price }}</text>
        </view>
      </view>
      <view class="foot">
        <text v-if="loading">加载中...</text>
        <text v-else-if="finished && products.length">没有更多了</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { onLoad, onReachBottom } from '@dcloudio/uni-app'
import { getShop, listProducts } from '../../api/index'

const placeholder = 'https://picsum.photos/id/96/400/400'
const shopId = ref(0)
const shop = reactive({})
const products = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

function goProduct(id) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
}

async function loadShop() {
  try {
    const res = await getShop(shopId.value)
    Object.assign(shop, res || {})
    if (shop.name) uni.setNavigationBarTitle({ title: shop.name })
  } catch {
    uni.showToast({ title: '店铺不存在', icon: 'none' })
  }
}

async function loadProducts(reset = false) {
  if (!shopId.value || loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      products.value = []
    }
    const res = await listProducts({
      page: page.value,
      page_size: pageSize,
      shop_id: shopId.value,
    })
    const rows = res?.list || []
    products.value = reset ? rows : products.value.concat(rows)
    const total = res?.total || 0
    if (products.value.length >= total || rows.length < pageSize) {
      finished.value = true
    } else {
      page.value += 1
    }
  } catch {
    if (reset) products.value = []
  } finally {
    loading.value = false
  }
}

onLoad((q) => {
  shopId.value = Number(q?.id || 0)
  if (!shopId.value) {
    uni.showToast({ title: '参数错误', icon: 'none' })
    return
  }
  loadShop()
  loadProducts(true)
})

onReachBottom(() => loadProducts(false))
</script>

<style scoped>
.page { padding-bottom: 40rpx; background: #fafafa; min-height: 100vh; }
.hero { width: 100%; height: 360rpx; background: #e4e4e7; }
.panel {
  margin: -48rpx 24rpx 24rpx; background: #fff; border-radius: 16rpx;
  padding: 28rpx; position: relative; z-index: 1;
}
.row { display: flex; gap: 20rpx; align-items: center; margin-bottom: 16rpx; }
.logo { width: 96rpx; height: 96rpx; border-radius: 16rpx; background: #f4f4f5; }
.meta { flex: 1; min-width: 0; }
.name { display: block; font-size: 34rpx; font-weight: 700; color: #18181b; }
.sub { display: block; margin-top: 6rpx; font-size: 24rpx; color: #a1a1aa; }
.desc { font-size: 26rpx; color: #52525b; line-height: 1.5; }
.section { padding: 0 24rpx; }
.sec-title { font-size: 30rpx; font-weight: 600; color: #18181b; margin-bottom: 16rpx; display: block; }
.grid { display: flex; flex-wrap: wrap; gap: 16rpx; }
.goods {
  width: calc(50% - 8rpx); background: #fff; border-radius: 12rpx; overflow: hidden;
  padding-bottom: 16rpx;
}
.g-img { width: 100%; height: 280rpx; background: #f4f4f5; }
.g-name {
  display: block; margin: 12rpx 16rpx 8rpx; font-size: 26rpx; color: #27272a;
  min-height: 72rpx;
}
.price { margin: 0 16rpx; color: #d83636; font-size: 28rpx; font-weight: 600; }
.empty, .foot { text-align: center; color: #a1a1aa; font-size: 24rpx; padding: 32rpx 0; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2;
  overflow: hidden;
}
</style>
