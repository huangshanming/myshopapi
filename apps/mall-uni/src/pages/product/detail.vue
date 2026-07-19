<template>
  <view class="page" v-if="product">
    <image class="cover" :src="product.main_image || placeholder" mode="aspectFill" />
    <view class="card">
      <view class="price-row">
        <text v-if="seckillActive" class="badge">秒杀</text>
        <text class="price">¥{{ displayPrice }}</text>
        <text v-if="strikePrice" class="price-old">¥{{ strikePrice }}</text>
        <view class="fav" @tap="toggleFav">
          <text :class="{ on: favorited }">{{ favorited ? '♥' : '♡' }}</text>
          <text class="fav-n">{{ collectCount }}</text>
        </view>
      </view>
      <text v-if="seckillTip" class="seckill-tip">{{ seckillTip }}</text>
      <text class="title">{{ product.name }}</text>
      <text v-if="product.subtitle" class="sub">{{ product.subtitle }}</text>
      <view class="meta">
        <text v-if="seckillActive">秒杀库存 {{ seckill.seckill_stock }}</text>
        <text v-else>库存 {{ product.stock ?? '-' }}</text>
        <text>销量 {{ product.sold_count || 0 }}</text>
        <text v-if="product.review_count">评分 {{ Number(product.avg_rating || 0).toFixed(1) }}</text>
      </view>
    </view>

    <view class="card reviews">
      <view class="rev-head">
        <text class="desc-title">用户评价</text>
        <text v-if="reviewTotal" class="rev-more">共 {{ reviewTotal }} 条</text>
      </view>
      <view v-if="!reviews.length" class="desc-empty">暂无评价</view>
      <view v-for="r in reviews" :key="r.id" class="rev-item">
        <view class="rev-top">
          <text class="rev-user">{{ r.is_anonymous ? '匿名用户' : (r.user_name || '用户') }}</text>
          <text class="rev-stars">{{ '★'.repeat(r.rating) }}{{ '☆'.repeat(5 - r.rating) }}</text>
        </view>
        <text class="rev-content">{{ r.content || '（无文字）' }}</text>
        <view v-if="r.images?.length" class="rev-imgs">
          <image v-for="img in r.images" :key="img.id || img.url" :src="img.url" mode="aspectFill" class="rev-img" />
        </view>
        <text v-if="r.merchant_reply" class="rev-reply">商家回复：{{ r.merchant_reply }}</text>
      </view>
    </view>

    <view class="card desc">
      <text class="desc-title">商品详情</text>
      <rich-text v-if="descHtml" class="desc-body" :nodes="descHtml" />
      <text v-else class="desc-empty">暂无详情介绍</text>
    </view>
    <view class="bar">
      <view class="cart-entry" @tap="goCart">
        <view class="cart-icon">
          <view class="cart-body" />
          <view class="cart-wheel l" />
          <view class="cart-wheel r" />
        </view>
        <text v-if="cartCount > 0" class="badge-num">{{ cartCount > 99 ? '99+' : cartCount }}</text>
        <text class="cart-label">购物车</text>
      </view>
      <button class="btn-ghost" @tap="addCart">加入购物车</button>
      <button class="btn-buy" @tap="buyNow">立即购买</button>
    </view>
  </view>
  <view v-else class="empty">{{ loading ? '加载中...' : '商品不存在' }}</view>
</template>

<script setup>
import { onLoad, onShow } from '@dcloudio/uni-app'
import { computed, ref } from 'vue'
import {
  addFavorite, getFavoriteStatus, getProductDetail, getSeckillEntry,
  listProductReviews, removeFavorite,
} from '../../api/index'
import { addToCart, getCartCount, setCheckoutPayload } from '../../stores/cart'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/400/400'
const product = ref(null)
const seckill = ref(null)
const loading = ref(false)
const cartCount = ref(0)
const favorited = ref(false)
const collectCount = ref(0)
const reviews = ref([])
const reviewTotal = ref(0)
let productId = 0
let seckillEntryId = 0

const seckillActive = computed(() => !!seckill.value?.seckill_available)

const displayPrice = computed(() => {
  if (seckillActive.value) return seckill.value.seckill_price
  return product.value?.sale_price
})

const strikePrice = computed(() => {
  if (seckillActive.value) {
    return seckill.value.origin_price || product.value?.market_price || product.value?.sale_price
  }
  return product.value?.market_price || null
})

const seckillTip = computed(() => {
  if (!seckillEntryId) return ''
  if (seckillActive.value) return `秒杀进行中，仅剩 ${seckill.value.seckill_stock} 件`
  return '秒杀库存不足，已按原价购买'
})

const descHtml = computed(() => {
  let html = (product.value?.description || '').trim()
  if (!html || html === '<p><br></p>' || html === '<p></p>') return ''
  html = html.replace(/<img(?![^>]*style=)/gi, '<img style="max-width:100%;height:auto;display:block;margin:8px 0;"')
  html = html.replace(/<img([^>]*?)style="/gi, '<img$1style="max-width:100%;height:auto;display:block;margin:8px 0;')
  return html
})

function refreshCartCount() {
  cartCount.value = getCartCount()
}

function buildItem() {
  const item = {
    product_id: productId,
    sku_id: 0,
    quantity: 1,
    name: product.value.name,
    image: product.value.main_image || '',
    price: Number(displayPrice.value) || 0,
    shop_id: Number(product.value.shop_id) || 0,
  }
  if (seckillActive.value) item.seckill_entry_id = seckillEntryId
  return item
}

onLoad((q) => {
  productId = Number(q.id || 0)
  seckillEntryId = Number(q.seckill_entry_id || 0)
  load()
})

onShow(refreshCartCount)

async function loadFav() {
  collectCount.value = Number(product.value?.collect_count || 0)
  if (!isLoggedIn()) {
    favorited.value = false
    return
  }
  try {
    const res = await getFavoriteStatus(productId)
    favorited.value = !!res?.favorited
  } catch {
    favorited.value = false
  }
}

async function loadReviews() {
  try {
    const res = await listProductReviews(productId, { page: 1, page_size: 5 })
    reviews.value = res?.list || []
    reviewTotal.value = res?.total || 0
  } catch {
    reviews.value = []
    reviewTotal.value = 0
  }
}

async function load() {
  if (!productId) return
  loading.value = true
  try {
    const res = await getProductDetail(productId)
    product.value = res || null
    if (seckillEntryId) {
      try {
        const sres = await getSeckillEntry(seckillEntryId)
        seckill.value = sres || null
        if (seckill.value && Number(seckill.value.product_id) !== productId) {
          seckill.value = null
        }
        if (seckill.value && !seckill.value.seckill_available) {
          uni.showToast({ title: '秒杀库存不足，已按原价展示', icon: 'none' })
        }
      } catch {
        seckill.value = null
      }
    }
    await Promise.all([loadFav(), loadReviews()])
  } catch {
    product.value = null
  } finally {
    loading.value = false
  }
}

async function toggleFav() {
  if (!isLoggedIn()) {
    uni.navigateTo({
      url: `/pages/login/login?redirect=${encodeURIComponent(`/pages/product/detail?id=${productId}`)}`,
    })
    return
  }
  try {
    if (favorited.value) {
      await removeFavorite(productId)
      favorited.value = false
      collectCount.value = Math.max(0, collectCount.value - 1)
    } else {
      await addFavorite(productId)
      favorited.value = true
      collectCount.value += 1
    }
  } catch { /* handled */ }
}

function goCart() {
  uni.switchTab({ url: '/pages/cart/index' })
}

function addCart() {
  if (!product.value) return
  addToCart(buildItem())
  refreshCartCount()
  uni.showToast({ title: '已加入购物车', icon: 'success' })
}

function buyNow() {
  if (!product.value) return
  const redirect = seckillEntryId
    ? `/pages/product/detail?id=${productId}&seckill_entry_id=${seckillEntryId}`
    : `/pages/product/detail?id=${productId}`
  if (!isLoggedIn()) {
    uni.navigateTo({
      url: `/pages/login/login?redirect=${encodeURIComponent(redirect)}`,
    })
    return
  }
  setCheckoutPayload({
    from: 'buy_now',
    items: [buildItem()],
  })
  uni.navigateTo({ url: '/pages/order/confirm' })
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
.price-row { display: flex; align-items: baseline; flex-wrap: wrap; gap: 8rpx; }
.fav { margin-left: auto; display: flex; flex-direction: column; align-items: center; color: #a1a1aa; font-size: 40rpx; line-height: 1; }
.fav .on { color: #d83636; }
.fav-n { font-size: 20rpx; color: #71717a; margin-top: 4rpx; }
.badge {
  background: #d83636; color: #fff; font-size: 20rpx; padding: 2rpx 10rpx;
  border-radius: 6rpx; margin-right: 8rpx; font-weight: 600;
}
.seckill-tip { display: block; margin-top: 10rpx; font-size: 22rpx; color: #d97706; }
.title { display: block; font-size: 32rpx; font-weight: 600; margin-top: 16rpx; }
.sub { display: block; color: #71717a; font-size: 24rpx; margin-top: 8rpx; }
.meta { display: flex; gap: 32rpx; margin-top: 20rpx; color: #71717a; font-size: 24rpx; flex-wrap: wrap; }
.desc-title { font-weight: 600; font-size: 28rpx; }
.desc-body { display: block; margin-top: 16rpx; color: #555; font-size: 28rpx; line-height: 1.7; word-break: break-word; }
.desc-empty { display: block; margin-top: 16rpx; color: #a1a1aa; font-size: 26rpx; }
.rev-head { display: flex; justify-content: space-between; align-items: center; }
.rev-more { font-size: 22rpx; color: #a1a1aa; }
.rev-item { padding: 20rpx 0; border-top: 1rpx solid #f5f5f5; margin-top: 12rpx; }
.rev-top { display: flex; justify-content: space-between; margin-bottom: 8rpx; }
.rev-user { font-size: 24rpx; color: #71717a; }
.rev-stars { color: #f5a623; font-size: 22rpx; }
.rev-content { font-size: 26rpx; color: #3f3f46; line-height: 1.5; }
.rev-imgs { display: flex; gap: 8rpx; margin-top: 12rpx; flex-wrap: wrap; }
.rev-img { width: 120rpx; height: 120rpx; border-radius: 8rpx; }
.rev-reply { display: block; margin-top: 12rpx; font-size: 24rpx; color: #71717a; background: #faf8f4; padding: 12rpx; border-radius: 8rpx; }
.bar {
  position: fixed; left: 0; right: 0; bottom: 0; display: flex; gap: 16rpx; align-items: center;
  padding: 16rpx 24rpx calc(16rpx + env(safe-area-inset-bottom)); background: #fff;
  border-top: 1rpx solid #f0f0f0;
}
.cart-entry {
  position: relative; width: 96rpx; display: flex; flex-direction: column;
  align-items: center; justify-content: center; flex-shrink: 0;
}
.cart-icon { position: relative; width: 40rpx; height: 34rpx; margin-bottom: 2rpx; }
.cart-body {
  position: absolute; left: 4rpx; top: 0; right: 0; bottom: 10rpx;
  border: 3rpx solid #52525b; border-radius: 4rpx 4rpx 6rpx 6rpx;
}
.cart-wheel {
  position: absolute; bottom: 0; width: 10rpx; height: 10rpx; border-radius: 50%;
  background: #52525b;
}
.cart-wheel.l { left: 8rpx; }
.cart-wheel.r { right: 2rpx; }
.cart-label { font-size: 20rpx; color: #71717a; margin-top: 4rpx; }
.badge-num {
  position: absolute; top: -8rpx; right: 8rpx; min-width: 28rpx; height: 28rpx;
  padding: 0 6rpx; border-radius: 999rpx; background: #d83636; color: #fff;
  font-size: 18rpx; line-height: 28rpx; text-align: center;
}
.btn-ghost, .btn-buy {
  flex: 1; height: 80rpx; line-height: 80rpx; border-radius: 999rpx; font-size: 26rpx; margin: 0;
}
.btn-ghost { background: #f5f5f5; color: #333; }
.btn-buy { background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; }
.empty { text-align: center; padding: 120rpx; color: #71717a; }
</style>
