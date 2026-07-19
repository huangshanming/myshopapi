<template>
  <view class="page">
    <view class="card addr" @tap="goAddress">
      <view v-if="address" class="addr-body">
        <view class="addr-top">
          <text class="name">{{ address.receiver_name }}</text>
          <text class="phone">{{ address.receiver_phone }}</text>
          <text v-if="address.is_default" class="tag">默认</text>
        </view>
        <text class="line">{{ fullAddr(address) }}</text>
      </view>
      <view v-else class="addr-empty">
        <text>请选择收货地址</text>
      </view>
      <text class="arrow">›</text>
    </view>

    <view class="card">
      <text class="sec">商品</text>
      <view v-for="it in items" :key="rowKey(it)" class="item">
        <image class="cover" :src="it.image || placeholder" mode="aspectFill" />
        <view class="info">
          <text class="name">{{ it.name }}</text>
          <view class="row">
            <text class="price">¥{{ it.price }}</text>
            <text class="qty">×{{ it.quantity }}</text>
          </view>
        </view>
      </view>
    </view>

    <view class="card total-card">
      <text>合计</text>
      <text class="total">¥{{ total.toFixed(2) }}</text>
    </view>

    <button class="btn" @tap="goPay">去支付</button>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { listAddresses } from '../../api/index'
import { getCheckoutPayload, setCheckoutPayload } from '../../stores/cart'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/200/200'
const items = ref([])
const address = ref(null)
const addressId = ref(0)

const total = computed(() =>
  items.value.reduce((s, x) => s + Number(x.price) * Number(x.quantity), 0),
)

function rowKey(it) {
  return `${it.product_id}_${it.sku_id || 0}_${it.seckill_entry_id || 0}`
}

function fullAddr(a) {
  if (!a) return ''
  return `${a.province || ''}${a.city || ''}${a.district || ''}${a.detail || ''}`
}

async function loadAddresses() {
  try {
    const res = await listAddresses()
    const list = res.data || []
    const picked = uni.getStorageSync('mymall_picked_address_id')
    if (picked) {
      addressId.value = Number(picked)
      uni.removeStorageSync('mymall_picked_address_id')
    }
    if (addressId.value) {
      address.value = list.find((x) => Number(x.id) === addressId.value) || null
    }
    if (!address.value) {
      address.value = list.find((x) => x.is_default === 1) || list[0] || null
      addressId.value = address.value?.id || 0
    }
  } catch {
    address.value = null
  }
}

function goAddress() {
  uni.navigateTo({ url: '/pages/address/list?from=confirm' })
}

function goPay() {
  if (!addressId.value || !address.value) {
    uni.showToast({ title: '请选择收货地址', icon: 'none' })
    return
  }
  if (!items.value.length) {
    uni.showToast({ title: '没有可结算商品', icon: 'none' })
    return
  }
  const payload = getCheckoutPayload() || {}
  setCheckoutPayload({
    ...payload,
    address_id: addressId.value,
    items: items.value,
  })
  uni.navigateTo({ url: '/pages/order/pay' })
}

onShow(() => {
  if (!isLoggedIn()) {
    uni.redirectTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/order/confirm') })
    return
  }
  const payload = getCheckoutPayload()
  if (!payload?.items?.length) {
    uni.showToast({ title: '结算信息已失效', icon: 'none' })
    setTimeout(() => uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/cart/index' }) }), 500)
    return
  }
  items.value = payload.items
  if (payload.address_id) addressId.value = Number(payload.address_id)
  loadAddresses()
})
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; }
.card {
  background: #fff; border-radius: 24rpx; padding: 28rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.addr { display: flex; align-items: center; gap: 16rpx; }
.addr-body { flex: 1; min-width: 0; }
.addr-top { display: flex; align-items: center; gap: 16rpx; margin-bottom: 8rpx; }
.name { font-weight: 600; font-size: 30rpx; }
.phone { color: #71717a; font-size: 26rpx; }
.tag {
  font-size: 20rpx; color: #c8a876; border: 1rpx solid #c8a876;
  padding: 0 8rpx; border-radius: 6rpx;
}
.line { display: block; color: #52525b; font-size: 26rpx; line-height: 1.5; }
.addr-empty { flex: 1; color: #a1a1aa; font-size: 28rpx; }
.arrow { color: #c8a876; font-size: 40rpx; }
.sec { font-weight: 600; font-size: 28rpx; display: block; margin-bottom: 16rpx; }
.item { display: flex; gap: 20rpx; padding: 16rpx 0; border-top: 1rpx solid #f5f5f5; }
.item:first-of-type { border-top: none; }
.cover { width: 140rpx; height: 140rpx; border-radius: 12rpx; background: #f3f3f3; }
.info { flex: 1; }
.item .name { display: block; font-size: 28rpx; }
.row { display: flex; justify-content: space-between; margin-top: 20rpx; }
.price { color: #d83636; font-weight: 600; }
.qty { color: #71717a; }
.total-card { display: flex; justify-content: space-between; align-items: center; font-size: 28rpx; }
.total { color: #d83636; font-size: 36rpx; font-weight: 700; }
.btn {
  margin-top: 24rpx; height: 88rpx; line-height: 88rpx; border-radius: 999rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; font-size: 30rpx;
}
</style>
