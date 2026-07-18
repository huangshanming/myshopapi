<template>
  <view class="page" v-if="product">
    <image class="cover" :src="product.main_image || placeholder" mode="aspectFill" />
    <view class="card">
      <view class="price-row">
        <text class="price">¥{{ product.sale_price }}</text>
        <text v-if="product.market_price" class="price-old">¥{{ product.market_price }}</text>
      </view>
      <text class="title">{{ product.name }}</text>
      <text v-if="product.subtitle" class="sub">{{ product.subtitle }}</text>
      <view class="meta">
        <text>库存 {{ product.stock ?? '-' }}</text>
        <text>销量 {{ product.sold_count || 0 }}</text>
      </view>
    </view>
    <view class="card desc">
      <text class="desc-title">商品详情</text>
      <text class="desc-body">{{ product.description || '暂无详情介绍' }}</text>
    </view>
    <view class="bar">
      <button class="btn-ghost" @tap="goHome">首页</button>
      <button class="btn-buy" :loading="buying" @tap="buy">立即购买</button>
    </view>
  </view>
  <view v-else class="empty">{{ loading ? '加载中...' : '商品不存在' }}</view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { createOrder, getProductDetail } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/400/400'
const product = ref(null)
const loading = ref(false)
const buying = ref(false)
let productId = 0

onLoad((q) => {
  productId = Number(q.id || 0)
  load()
})

async function load() {
  if (!productId) return
  loading.value = true
  try {
    const res = await getProductDetail(productId)
    product.value = res.data || null
  } catch {
    product.value = null
  } finally {
    loading.value = false
  }
}

function goHome() {
  uni.switchTab({ url: '/pages/index/index' })
}

async function buy() {
  if (!product.value) return
  if (!isLoggedIn()) {
    uni.navigateTo({
      url: `/pages/login/login?redirect=${encodeURIComponent(`/pages/product/detail?id=${productId}`)}`,
    })
    return
  }
  buying.value = true
  try {
    const res = await createOrder([
      { product_id: productId, sku_id: 0, quantity: 1 },
    ])
    const id = res.data?.id
    uni.showToast({ title: '下单成功', icon: 'success' })
    setTimeout(() => {
      if (id) uni.redirectTo({ url: `/pages/order/detail?id=${id}` })
      else uni.redirectTo({ url: '/pages/order/list' })
    }, 500)
  } catch {
    /* toast in request */
  } finally {
    buying.value = false
  }
}
</script>

<style scoped>
.page { padding-bottom: 140rpx; }
.cover { width: 100%; height: 750rpx; background: #f3f3f3; }
.card {
  margin: 24rpx 32rpx 0; padding: 28rpx; background: #fff; border-radius: 32rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.price { color: #d83636; font-size: 44rpx; font-weight: 700; }
.price-old { color: #71717a; font-size: 24rpx; text-decoration: line-through; margin-left: 12rpx; }
.price-row { display: flex; align-items: baseline; }
.title { display: block; font-size: 32rpx; font-weight: 600; margin-top: 16rpx; }
.sub { display: block; color: #71717a; font-size: 24rpx; margin-top: 8rpx; }
.meta { display: flex; gap: 32rpx; margin-top: 20rpx; color: #71717a; font-size: 24rpx; }
.desc-title { font-weight: 600; font-size: 28rpx; }
.desc-body { display: block; margin-top: 16rpx; color: #555; line-height: 1.6; white-space: pre-wrap; }
.bar {
  position: fixed; left: 0; right: 0; bottom: 0; display: flex; gap: 20rpx;
  padding: 16rpx 32rpx calc(16rpx + env(safe-area-inset-bottom)); background: #fff;
  border-top: 1rpx solid #f0f0f0;
}
.btn-ghost, .btn-buy {
  flex: 1; height: 80rpx; line-height: 80rpx; border-radius: 999rpx; font-size: 28rpx; margin: 0;
}
.btn-ghost { background: #f5f5f5; color: #333; }
.btn-buy { background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; }
.empty { text-align: center; padding: 120rpx; color: #71717a; }
</style>
