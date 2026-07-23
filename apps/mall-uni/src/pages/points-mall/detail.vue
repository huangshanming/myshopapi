<template>
  <view class="page" v-if="product">
    <image class="cover" :src="product.cover_url || placeholder" mode="aspectFill" />
    <view class="panel">
      <view class="price-row">
        <text class="price">{{ product.points_price }}</text>
        <text class="unit">积分</text>
      </view>
      <text class="name">{{ product.name }}</text>
      <view class="tags">
        <text class="tag">库存 {{ product.stock ?? 0 }}</text>
        <text v-if="product.per_user_limit > 0" class="tag">每人限兑 {{ product.per_user_limit }}</text>
      </view>
      <text v-if="product.description" class="desc">{{ product.description }}</text>
    </view>

    <view class="panel addr" @tap="pickAddress">
      <view class="addr-top">
        <text class="addr-label">收货地址</text>
        <text class="addr-link">{{ address ? '更换' : '去选择' }} ›</text>
      </view>
      <template v-if="address">
        <text class="addr-name">{{ address.receiver_name }} {{ address.receiver_phone }}</text>
        <text class="addr-line">{{ fullAddr(address) }}</text>
      </template>
      <text v-else class="addr-empty">请选择收货地址</text>
    </view>

    <view class="safe" />
    <view class="bar">
      <view class="bar-points">
        <text class="bar-label">我的积分</text>
        <text class="bar-num">{{ points }}</text>
      </view>
      <button class="ex-btn" :disabled="exchanging" @tap="onExchange">立即兑换</button>
    </view>
  </view>
  <view v-else class="loading">{{ loadErr || '加载中...' }}</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import {
  exchangePointsMall,
  getPointsMallProduct,
  getUserPoints,
  listAddresses,
} from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/750/750'
const productId = ref(0)
const product = ref(null)
const points = ref(0)
const address = ref(null)
const exchanging = ref(false)
const loadErr = ref('')

function fullAddr(a) {
  return `${a.province || ''}${a.city || ''}${a.district || ''}${a.detail || ''}`
}

async function loadProduct() {
  if (!productId.value) return
  try {
    const res = await getPointsMallProduct(productId.value)
    product.value = res?.product || null
    if (!product.value) loadErr.value = '商品不存在'
  } catch (e) {
    loadErr.value = e.message || '加载失败'
    product.value = null
  }
}

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

async function syncAddress() {
  if (!isLoggedIn()) {
    address.value = null
    return
  }
  try {
    const picked = uni.getStorageSync('mymall_picked_address_id')
    const res = await listAddresses()
    const list = res?.list || []
    if (picked) {
      uni.removeStorageSync('mymall_picked_address_id')
      const found = list.find((a) => String(a.id) === String(picked))
      if (found) {
        address.value = found
        return
      }
    }
    if (address.value) {
      const keep = list.find((a) => Number(a.id) === Number(address.value.id))
      if (keep) {
        address.value = keep
        return
      }
    }
    address.value = list.find((a) => a.is_default) || list[0] || null
  } catch {
    /* keep current */
  }
}

function pickAddress() {
  if (!isLoggedIn()) {
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent(`/pages/points-mall/detail?id=${productId.value}`),
    })
    return
  }
  uni.navigateTo({ url: '/pages/address/list?from=confirm' })
}

async function onExchange() {
  if (!product.value) return
  if (!isLoggedIn()) {
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent(`/pages/points-mall/detail?id=${productId.value}`),
    })
    return
  }
  if (!address.value?.id) {
    uni.showToast({ title: '请先选择收货地址', icon: 'none' })
    return
  }
  if ((product.value.stock ?? 0) <= 0) {
    uni.showToast({ title: '库存不足', icon: 'none' })
    return
  }
  if (points.value < (product.value.points_price || 0)) {
    uni.showToast({ title: '积分不足', icon: 'none' })
    return
  }
  const ok = await new Promise((resolve) => {
    uni.showModal({
      title: '确认兑换',
      content: `将消耗 ${product.value.points_price} 积分兑换「${product.value.name}」`,
      success: (r) => resolve(r.confirm),
    })
  })
  if (!ok) return
  exchanging.value = true
  try {
    const res = await exchangePointsMall({
      product_id: Number(productId.value),
      address_id: Number(address.value.id),
    })
    const order = res?.data || res
    const oid = order?.id
    uni.showToast({ title: '兑换成功', icon: 'success' })
    if (oid) {
      setTimeout(() => {
        uni.redirectTo({ url: `/pages/points-mall/order-detail?id=${oid}` })
      }, 400)
    } else {
      uni.navigateTo({ url: '/pages/points-mall/orders' })
    }
  } catch (e) {
    uni.showToast({ title: e.message || '兑换失败', icon: 'none' })
  } finally {
    exchanging.value = false
  }
}

onLoad((q) => {
  productId.value = Number(q.id) || 0
  loadProduct()
})

onShow(() => {
  loadPoints()
  syncAddress()
})
</script>

<style scoped>
.page { padding-bottom: 140rpx; min-height: 100vh; background: #f7f3ec; }
.cover { width: 100%; height: 720rpx; background: #eee; display: block; }
.panel {
  margin: 24rpx 32rpx; background: #fff; border-radius: 24rpx; padding: 28rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.price-row { display: flex; align-items: baseline; gap: 8rpx; }
.price { font-size: 48rpx; font-weight: 700; color: #b8860b; }
.unit { font-size: 24rpx; color: #8b7355; }
.name { display: block; margin-top: 16rpx; font-size: 34rpx; font-weight: 700; color: #18181b; }
.tags { display: flex; flex-wrap: wrap; gap: 12rpx; margin-top: 16rpx; }
.tag {
  font-size: 22rpx; color: #8b7355; background: #f7f3ec;
  padding: 6rpx 14rpx; border-radius: 8rpx;
}
.desc {
  display: block; margin-top: 20rpx; font-size: 26rpx; color: #52525b; line-height: 1.6;
  white-space: pre-wrap;
}
.addr-top { display: flex; justify-content: space-between; margin-bottom: 12rpx; }
.addr-label { font-size: 28rpx; font-weight: 600; }
.addr-link { font-size: 24rpx; color: #c8a876; }
.addr-name { display: block; font-size: 28rpx; font-weight: 600; }
.addr-line { display: block; margin-top: 8rpx; font-size: 26rpx; color: #52525b; line-height: 1.5; }
.addr-empty { font-size: 26rpx; color: #a1a1aa; }
.safe { height: 24rpx; }
.bar {
  position: fixed; left: 0; right: 0; bottom: 0;
  display: flex; align-items: center; gap: 24rpx;
  padding: 16rpx 32rpx calc(16rpx + env(safe-area-inset-bottom));
  background: #fff; box-shadow: 0 -4rpx 24rpx rgba(0,0,0,.04);
}
.bar-points { flex: 1; }
.bar-label { display: block; font-size: 22rpx; color: #a1a1aa; }
.bar-num { font-size: 32rpx; font-weight: 700; color: #5c4a32; }
.ex-btn {
  margin: 0; min-width: 260rpx; height: 80rpx; line-height: 80rpx;
  border-radius: 999rpx; background: linear-gradient(135deg, #bfa472, #d4b890);
  color: #fff; font-size: 30rpx;
}
.ex-btn[disabled] { opacity: .6; }
.loading { text-align: center; padding: 120rpx 32rpx; color: #a1a1aa; }
</style>
