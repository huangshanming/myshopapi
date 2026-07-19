<template>
  <view class="page">
    <view v-if="!list.length" class="empty">
      <text>购物车是空的</text>
      <view class="empty-btn" @tap="goHome">去逛逛</view>
    </view>

    <view v-for="it in list" :key="rowKey(it)" class="card">
      <view class="check" @tap.stop="toggle(it)">
        <text class="box" :class="{ on: it.selected }">{{ it.selected ? '✓' : '' }}</text>
      </view>
      <image
        class="cover"
        :src="it.image || placeholder"
        mode="aspectFill"
        @tap.stop="goDetail(it)"
      />
      <view class="info">
        <text class="name" @tap.stop="goDetail(it)">{{ it.name }}</text>
        <text v-if="it.seckill_entry_id" class="tag">秒杀</text>
        <text class="price">¥{{ it.price }}</text>
        <view class="qty-row">
          <text class="qty-btn" @tap.stop="changeQty(it, -1)">−</text>
          <text class="qty">{{ it.quantity }}</text>
          <text class="qty-btn" @tap.stop="changeQty(it, 1)">+</text>
          <text class="del" @tap.stop="remove(it)">删除</text>
        </view>
      </view>
    </view>

    <view v-if="list.length" class="bar">
      <view class="check" @tap="toggleAll">
        <text class="box" :class="{ on: allSelected }">{{ allSelected ? '✓' : '' }}</text>
        <text class="all">全选</text>
      </view>
      <view class="sum-wrap">
        <text class="sum-label">合计</text>
        <text class="sum">¥{{ total.toFixed(2) }}</text>
      </view>
      <view class="btn" :class="{ disabled: !selected.length }" @tap="checkout">
        结算({{ selected.length }})
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import {
  getCartItems, removeCartItem, setAllSelected, setCartSelected, setCheckoutPayload, updateCartQty,
} from '../../stores/cart'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/200/200'
const list = ref([])

const selected = computed(() => list.value.filter((x) => x.selected))
const allSelected = computed(() => list.value.length > 0 && selected.value.length === list.value.length)
const total = computed(() =>
  selected.value.reduce((s, x) => s + Number(x.price) * Number(x.quantity), 0),
)

function rowKey(it) {
  return `${it.product_id}_${it.sku_id || 0}_${it.seckill_entry_id || 0}`
}

function reload() {
  list.value = getCartItems().map((x) => ({ ...x, selected: x.selected !== false }))
}

function toggle(it) {
  list.value = setCartSelected(it.product_id, it.sku_id, it.seckill_entry_id, !it.selected)
}

function toggleAll() {
  list.value = setAllSelected(!allSelected.value)
}

function changeQty(it, delta) {
  const q = Number(it.quantity) + delta
  list.value = updateCartQty(it.product_id, it.sku_id, it.seckill_entry_id, q)
}

function remove(it) {
  list.value = removeCartItem(it.product_id, it.sku_id, it.seckill_entry_id)
}

function goHome() {
  uni.switchTab({ url: '/pages/index/index' })
}

function goDetail(it) {
  if (!it?.product_id) return
  let url = `/pages/product/detail?id=${it.product_id}`
  if (it.seckill_entry_id) url += `&seckill_entry_id=${it.seckill_entry_id}`
  uni.navigateTo({ url })
}

function checkout() {
  if (!selected.value.length) {
    uni.showToast({ title: '请先选择商品', icon: 'none' })
    return
  }
  const shops = new Set(selected.value.map((x) => Number(x.shop_id) || 0))
  if (shops.size > 1) {
    uni.showToast({ title: '请选择同一店铺商品结算', icon: 'none' })
    return
  }
  setCheckoutPayload({
    from: 'cart',
    items: selected.value.map((x) => ({
      product_id: x.product_id,
      sku_id: x.sku_id || 0,
      quantity: x.quantity,
      name: x.name,
      image: x.image,
      price: x.price,
      shop_id: x.shop_id,
      seckill_entry_id: x.seckill_entry_id || 0,
    })),
  })
  if (!isLoggedIn()) {
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/order/confirm'),
    })
    return
  }
  uni.navigateTo({ url: '/pages/order/confirm' })
}

onShow(reload)
</script>

<style scoped>
.page { padding: 24rpx 32rpx calc(200rpx + 100px + env(safe-area-inset-bottom)); }
.empty {
  text-align: center; color: #a1a1aa; padding: 120rpx 0;
  display: flex; flex-direction: column; align-items: center; gap: 24rpx;
}
.empty-btn {
  padding: 16rpx 40rpx; border-radius: 999rpx; font-size: 26rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff;
}
.card {
  display: flex; gap: 20rpx; align-items: center;
  background: #fff; border-radius: 24rpx; padding: 24rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.check { display: flex; align-items: center; gap: 8rpx; flex-shrink: 0; }
.box {
  width: 40rpx; height: 40rpx; border-radius: 50%; border: 2rpx solid #d4d4d8;
  display: flex; align-items: center; justify-content: center; font-size: 22rpx; color: #fff;
}
.box.on { background: #c8a876; border-color: #c8a876; }
.cover { width: 160rpx; height: 160rpx; border-radius: 16rpx; background: #f3f3f3; flex-shrink: 0; }
.info { flex: 1; min-width: 0; }
.name {
  display: block; font-size: 28rpx; font-weight: 500; color: #18181b;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.tag {
  display: inline-block; margin-top: 8rpx; font-size: 20rpx; color: #ef4444;
  background: #fef2f2; padding: 2rpx 10rpx; border-radius: 6rpx;
}
.price { display: block; color: #d83636; font-size: 30rpx; font-weight: 700; margin-top: 12rpx; }
.qty-row { display: flex; align-items: center; gap: 16rpx; margin-top: 16rpx; }
.qty-btn {
  width: 48rpx; height: 48rpx; line-height: 48rpx; text-align: center;
  background: #f5f5f5; border-radius: 8rpx; font-size: 28rpx;
}
.qty { min-width: 40rpx; text-align: center; }
.del { margin-left: auto; color: #a1a1aa; font-size: 24rpx; }
.bar {
  position: fixed; left: 0; right: 0;
  bottom: calc(50px + env(safe-area-inset-bottom));
  z-index: 20;
  display: flex; align-items: center; gap: 16rpx;
  padding: 16rpx 32rpx;
  background: #fff; border-top: 1rpx solid #f0f0f0;
  box-shadow: 0 -4rpx 16rpx rgba(0,0,0,.04);
}
.all { font-size: 26rpx; }
.sum-wrap { flex: 1; text-align: right; display: flex; align-items: baseline; justify-content: flex-end; gap: 8rpx; }
.sum-label { font-size: 24rpx; color: #71717a; }
.sum { font-size: 32rpx; font-weight: 700; color: #d83636; }
.btn {
  height: 72rpx; line-height: 72rpx; padding: 0 32rpx; border-radius: 999rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; font-size: 26rpx;
  flex-shrink: 0;
}
.btn.disabled { opacity: 0.45; }
</style>
