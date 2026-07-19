<template>
  <view class="page">
    <view class="card amount-card">
      <text class="label">应付金额</text>
      <text class="amount">¥{{ total.toFixed(2) }}</text>
    </view>

    <view class="card">
      <text class="sec">支付方式</text>
      <view class="pay-row">
        <view class="radio on" />
        <view class="pay-info">
          <text class="pay-name">余额支付</text>
          <text class="pay-sub">可用余额 ¥{{ balance }}</text>
        </view>
      </view>
    </view>

    <button class="btn" :loading="paying" @tap="pay">确认支付</button>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { createOrder, getUserWallet } from '../../api/index'
import { clearCartItems, clearCheckoutPayload, getCheckoutPayload } from '../../stores/cart'
import { isLoggedIn } from '../../stores/user'

const items = ref([])
const addressId = ref(0)
const userCouponId = ref(0)
const fromCart = ref(false)
const balance = ref(0)
const paying = ref(false)
const payAmount = ref(0)

const total = computed(() => {
  if (payAmount.value > 0) return payAmount.value
  return items.value.reduce((s, x) => s + Number(x.price) * Number(x.quantity), 0)
})

async function loadWallet() {
  try {
    const res = await getUserWallet()
    balance.value = Number(res?.balance) || 0
  } catch {
    balance.value = 0
  }
}

async function pay() {
  if (!addressId.value) {
    uni.showToast({ title: '请选择收货地址', icon: 'none' })
    return
  }
  if (balance.value + 0.0001 < total.value) {
    uni.showToast({ title: '余额不足', icon: 'none' })
    return
  }
  paying.value = true
  try {
    const orderItems = items.value.map((x) => {
      const it = {
        product_id: x.product_id,
        sku_id: x.sku_id || 0,
        quantity: x.quantity,
      }
      if (x.seckill_entry_id) it.seckill_entry_id = x.seckill_entry_id
      return it
    })
    const res = await createOrder(orderItems, addressId.value, userCouponId.value)
    const id = res?.id
    if (fromCart.value) {
      clearCartItems(items.value)
    }
    clearCheckoutPayload()
    uni.showToast({ title: '支付成功', icon: 'success' })
    setTimeout(() => {
      if (id) uni.redirectTo({ url: `/pages/order/detail?id=${id}` })
      else uni.redirectTo({ url: '/pages/order/list' })
    }, 500)
  } catch {
    /* toast in request */
  } finally {
    paying.value = false
  }
}

onLoad(() => {
  if (!isLoggedIn()) {
    uni.redirectTo({ url: '/pages/login/login' })
    return
  }
  const payload = getCheckoutPayload()
  if (!payload?.items?.length || !payload.address_id) {
    uni.showToast({ title: '请先确认订单', icon: 'none' })
    setTimeout(() => uni.navigateBack(), 500)
    return
  }
  items.value = payload.items
  addressId.value = Number(payload.address_id)
  userCouponId.value = Number(payload.user_coupon_id) || 0
  payAmount.value = Number(payload.pay_amount) || 0
  fromCart.value = payload.from === 'cart'
  loadWallet()
})
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; }
.card {
  background: #fff; border-radius: 24rpx; padding: 28rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.amount-card { text-align: center; padding: 48rpx 28rpx; }
.label { display: block; color: #71717a; font-size: 26rpx; }
.amount { display: block; margin-top: 12rpx; font-size: 56rpx; font-weight: 700; color: #d83636; }
.sec { font-weight: 600; font-size: 28rpx; display: block; margin-bottom: 20rpx; }
.pay-row { display: flex; align-items: center; gap: 20rpx; }
.radio {
  width: 40rpx; height: 40rpx; border-radius: 50%; border: 2rpx solid #d4d4d8;
}
.radio.on {
  border-color: #c8a876; background: #c8a876;
  box-shadow: inset 0 0 0 8rpx #fff;
}
.pay-name { display: block; font-size: 30rpx; font-weight: 500; }
.pay-sub { display: block; margin-top: 6rpx; color: #71717a; font-size: 24rpx; }
.btn {
  margin-top: 32rpx; height: 88rpx; line-height: 88rpx; border-radius: 999rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; font-size: 30rpx;
}
</style>
