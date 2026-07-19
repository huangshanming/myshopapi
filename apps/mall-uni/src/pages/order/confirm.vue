<template>
  <view class="page">
    <view class="addr-card" @tap="goAddress">
      <view class="addr-icon">
        <view class="pin" />
      </view>
      <view v-if="address" class="addr-body">
        <view class="addr-top">
          <text class="name">{{ address.receiver_name }}</text>
          <text class="phone">{{ address.receiver_phone }}</text>
          <text v-if="address.is_default" class="tag">默认</text>
        </view>
        <text class="line">{{ fullAddr(address) }}</text>
      </view>
      <view v-else class="addr-empty">
        <text class="empty-title">请添加收货地址</text>
        <text class="empty-sub">配送地址未选择，无法提交订单</text>
      </view>
      <text class="arrow">›</text>
    </view>

    <view class="shop-card">
      <view class="shop-head">
        <image class="shop-logo" :src="shop.logo || shopPlaceholder" mode="aspectFill" />
        <view class="shop-meta">
          <text class="shop-name">{{ shop.name || '店铺' }}</text>
          <text class="shop-sub">{{ shopLine }}</text>
        </view>
      </view>

      <view v-for="it in items" :key="rowKey(it)" class="goods">
        <image class="cover" :src="it.image || placeholder" mode="aspectFill" />
        <view class="goods-info">
          <text class="goods-name">{{ it.name }}</text>
          <text v-if="it.seckill_entry_id" class="sku-tip">秒杀价</text>
          <view class="goods-foot">
            <text class="price">¥{{ Number(it.price).toFixed(2) }}</text>
            <text class="qty">x{{ it.quantity }}</text>
          </view>
        </view>
        <text class="line-amt">¥{{ (Number(it.price) * Number(it.quantity)).toFixed(2) }}</text>
      </view>

      <view class="fee-row">
        <text>配送方式</text>
        <text class="fee-val">快递 免邮</text>
      </view>
      <view class="fee-row">
        <text>商品金额</text>
        <text class="fee-val">¥{{ goodsAmount.toFixed(2) }}</text>
      </view>
      <view class="fee-row">
        <text>运费</text>
        <text class="fee-val">¥0.00</text>
      </view>
      <view class="fee-row total-row">
        <text>小计</text>
        <text class="total">¥{{ total.toFixed(2) }}</text>
      </view>
    </view>

    <view class="tip-card">
      <text class="tip-title">订单说明</text>
      <text class="tip-line">· 支付方式：余额支付（下一步选择）</text>
      <text class="tip-line">· 下单后将冻结对应余额，确认发货后实扣</text>
      <text class="tip-line">· 共 {{ itemCount }} 件商品</text>
    </view>

    <view class="safe" />

    <view class="bar">
      <view class="bar-left">
        <text class="bar-label">合计</text>
        <text class="bar-total">¥{{ total.toFixed(2) }}</text>
      </view>
      <button class="bar-btn" @tap="goPay">去支付</button>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getShop, listAddresses } from '../../api/index'
import { getCheckoutPayload, setCheckoutPayload } from '../../stores/cart'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/200/200'
const shopPlaceholder = 'https://picsum.photos/id/20/100/100'
const items = ref([])
const address = ref(null)
const addressId = ref(0)
const shop = ref({})

const goodsAmount = computed(() =>
  items.value.reduce((s, x) => s + Number(x.price) * Number(x.quantity), 0),
)
const total = computed(() => goodsAmount.value)
const itemCount = computed(() =>
  items.value.reduce((n, x) => n + (Number(x.quantity) || 0), 0),
)
const shopLine = computed(() => {
  const parts = [shop.value.province, shop.value.city, shop.value.district].filter(Boolean)
  if (parts.length) return parts.join(' · ')
  return shop.value.category ? `主营 ${shop.value.category}` : '官方店铺'
})

function rowKey(it) {
  return `${it.product_id}_${it.sku_id || 0}_${it.seckill_entry_id || 0}`
}

function fullAddr(a) {
  if (!a) return ''
  return `${a.province || ''}${a.city || ''}${a.district || ''}${a.detail || ''}`
}

async function loadShop() {
  const shopID = Number(items.value[0]?.shop_id || 0)
  if (!shopID) {
    shop.value = { name: '店铺' }
    return
  }
  try {
    const res = await getShop(shopID)
    shop.value = res.data || { name: `店铺 #${shopID}` }
  } catch {
    shop.value = { name: `店铺 #${shopID}` }
  }
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

onShow(async () => {
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
  await Promise.all([loadAddresses(), loadShop()])
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 20rpx 24rpx 180rpx;
  background: linear-gradient(180deg, #f3ebe0 0%, #fafafa 180rpx);
}
.addr-card {
  display: flex; align-items: flex-start; gap: 16rpx;
  background: #fff; border-radius: 20rpx; padding: 28rpx 24rpx;
  margin-bottom: 20rpx; box-shadow: 0 4rpx 20rpx rgba(0,0,0,.04);
  position: relative;
}
.addr-card::after {
  content: ''; position: absolute; left: 0; right: 0; bottom: 0; height: 6rpx;
  background: repeating-linear-gradient(90deg, #f87171 0 16rpx, #60a5fa 16rpx 32rpx);
  border-radius: 0 0 20rpx 20rpx;
}
.addr-icon { padding-top: 8rpx; }
.pin {
  width: 28rpx; height: 28rpx; border-radius: 50% 50% 50% 0; background: #ef4444;
  transform: rotate(-45deg);
}
.addr-body { flex: 1; min-width: 0; padding-bottom: 8rpx; }
.addr-top { display: flex; align-items: center; flex-wrap: wrap; gap: 12rpx; margin-bottom: 10rpx; }
.name { font-weight: 700; font-size: 32rpx; color: #18181b; }
.phone { color: #52525b; font-size: 28rpx; }
.tag {
  font-size: 20rpx; color: #fff; background: #c8a876;
  padding: 2rpx 10rpx; border-radius: 6rpx;
}
.line { display: block; color: #3f3f46; font-size: 26rpx; line-height: 1.55; }
.addr-empty { flex: 1; padding-bottom: 8rpx; }
.empty-title { display: block; font-size: 30rpx; font-weight: 600; color: #18181b; }
.empty-sub { display: block; margin-top: 8rpx; font-size: 24rpx; color: #a1a1aa; }
.arrow { color: #a1a1aa; font-size: 36rpx; padding-top: 4rpx; }

.shop-card {
  background: #fff; border-radius: 20rpx; padding: 24rpx;
  margin-bottom: 20rpx; box-shadow: 0 4rpx 20rpx rgba(0,0,0,.04);
}
.shop-head {
  display: flex; align-items: center; gap: 16rpx;
  padding-bottom: 20rpx; border-bottom: 1rpx solid #f4f4f5;
}
.shop-logo { width: 64rpx; height: 64rpx; border-radius: 12rpx; background: #f4f4f5; }
.shop-meta { flex: 1; min-width: 0; }
.shop-name { display: block; font-size: 28rpx; font-weight: 700; color: #18181b; }
.shop-sub { display: block; margin-top: 4rpx; font-size: 22rpx; color: #a1a1aa; }

.goods {
  display: flex; gap: 20rpx; padding: 24rpx 0;
  border-bottom: 1rpx solid #fafafa;
}
.cover { width: 160rpx; height: 160rpx; border-radius: 12rpx; background: #f4f4f5; flex-shrink: 0; }
.goods-info { flex: 1; min-width: 0; display: flex; flex-direction: column; }
.goods-name {
  font-size: 28rpx; color: #27272a; line-height: 1.4;
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2; overflow: hidden;
}
.sku-tip {
  align-self: flex-start; margin-top: 10rpx;
  font-size: 20rpx; color: #ef4444; background: #fef2f2;
  padding: 2rpx 10rpx; border-radius: 6rpx;
}
.goods-foot { margin-top: auto; display: flex; justify-content: space-between; align-items: baseline; }
.price { color: #18181b; font-size: 28rpx; font-weight: 600; }
.qty { color: #a1a1aa; font-size: 24rpx; }
.line-amt { font-size: 26rpx; color: #52525b; padding-top: 4rpx; flex-shrink: 0; }

.fee-row {
  display: flex; justify-content: space-between; align-items: center;
  padding: 14rpx 0; font-size: 26rpx; color: #52525b;
}
.fee-val { color: #18181b; }
.total-row { padding-top: 18rpx; border-top: 1rpx solid #f4f4f5; margin-top: 4rpx; }
.total { color: #ef4444; font-size: 34rpx; font-weight: 700; }

.tip-card {
  background: #fff; border-radius: 20rpx; padding: 24rpx;
  box-shadow: 0 4rpx 20rpx rgba(0,0,0,.04);
}
.tip-title { display: block; font-size: 26rpx; font-weight: 600; margin-bottom: 12rpx; color: #3f3f46; }
.tip-line { display: block; font-size: 24rpx; color: #71717a; line-height: 1.7; }
.safe { height: 20rpx; }

.bar {
  position: fixed; left: 0; right: 0; bottom: 0;
  display: flex; align-items: center; justify-content: space-between; gap: 24rpx;
  padding: 16rpx 28rpx calc(16rpx + env(safe-area-inset-bottom));
  background: #fff; border-top: 1rpx solid #f0f0f0;
  box-shadow: 0 -4rpx 20rpx rgba(0,0,0,.04);
}
.bar-left { display: flex; align-items: baseline; gap: 10rpx; }
.bar-label { font-size: 26rpx; color: #52525b; }
.bar-total { font-size: 40rpx; font-weight: 700; color: #ef4444; }
.bar-btn {
  margin: 0; height: 80rpx; line-height: 80rpx; padding: 0 48rpx; border-radius: 999rpx;
  background: linear-gradient(135deg, #f59e0b, #ef4444); color: #fff; font-size: 30rpx; font-weight: 600;
}
</style>
