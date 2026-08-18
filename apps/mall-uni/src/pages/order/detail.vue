<template>
  <view class="page" v-if="order">
    <template v-if="!showAfterSaleForm">
      <view class="card shop-card" @tap="goShop">
        <image class="shop-logo" :src="shop.logo || shopPlaceholder" mode="aspectFill" />
        <view class="shop-meta">
          <text class="shop-name">{{ shopDisplayName }}</text>
          <text class="shop-sub">{{ shopLine }}</text>
        </view>
        <text class="shop-arrow">进店 ›</text>
      </view>

      <view class="card">
        <view class="row">
          <text class="label">订单号</text>
          <text>{{ order.order_no }}</text>
        </view>
        <view class="row">
          <text class="label">状态</text>
          <text class="gold">{{ statusText(order.status) }}</text>
        </view>
        <view class="row">
          <text class="label">金额</text>
          <text class="price">¥{{ order.total_amount }}</text>
        </view>
        <view class="row">
          <text class="label">下单时间</text>
          <text>{{ order.created_at }}</text>
        </view>
        <view v-if="order.receiver_name || order.receiver_address" class="row col">
          <text class="label">收货信息</text>
          <text class="addr">{{ order.receiver_name }} {{ order.receiver_phone }}</text>
          <text class="addr">{{ order.receiver_address }}</text>
        </view>
        <view v-if="order.ship_company" class="row">
          <text class="label">物流</text>
          <text>{{ order.ship_company }} {{ order.ship_no }}</text>
        </view>
        <view v-if="afterSale.deadline" class="row">
          <text class="label">售后截止</text>
          <text>{{ afterSale.deadline }}</text>
        </view>
      </view>

      <view class="card">
        <text class="sec">商品</text>
        <view v-for="it in order.items || []" :key="it.id" class="item">
          <text class="name">{{ it.product_name }}</text>
          <text class="sub">×{{ it.quantity }} · ¥{{ it.price }}</text>
        </view>
      </view>

      <view class="actions">
        <button v-if="canCancel" class="btn outline" :loading="busy" @tap="onCancel">取消订单</button>
        <button v-if="canConfirm" class="btn primary" :loading="busy" @tap="onConfirm">确认收货</button>
        <button v-if="canReview" class="btn primary" @tap="goReview">去评价</button>
        <button v-if="canViewReview" class="btn outline" @tap="goViewReview">查看评价</button>
        <button v-if="canAfterSale" class="btn outline" @tap="openAfterSale">申请售后</button>
      </view>
    </template>

    <view v-else class="after-sale">
      <view class="card form">
        <text class="sec">申请售后</text>
        <text class="hint">订单 {{ order.order_no }} · 最多退 ¥{{ payAmount }}</text>
        <view class="field">
          <text class="flabel">类型</text>
          <view class="seg">
            <text
              class="seg-item"
              :class="{ on: afterForm.type === 'refund' }"
              @tap="afterForm.type = 'refund'"
            >仅退款</text>
            <text
              class="seg-item"
              :class="{ on: afterForm.type === 'return_refund' }"
              @tap="afterForm.type = 'return_refund'"
            >退货退款</text>
          </view>
        </view>
        <view class="field">
          <text class="flabel">退款金额</text>
          <view class="amount-row">
            <text class="amount-prefix">¥</text>
            <input
              class="input amount-input"
              type="digit"
              v-model="afterForm.amount"
              :placeholder="payAmount"
            />
          </view>
        </view>
        <view class="field">
          <text class="flabel">原因</text>
          <textarea
            class="textarea"
            v-model="afterForm.reason"
            maxlength="200"
            placeholder="请填写售后原因"
            :auto-height="true"
          />
        </view>
      </view>
      <view class="foot-actions">
        <button class="btn outline" :disabled="busy" @tap="showAfterSaleForm = false">取消</button>
        <button class="btn primary" :loading="busy" @tap="onSubmitAfterSale">提交申请</button>
      </view>
    </view>
  </view>
  <view v-else class="empty">{{ loading ? '加载中...' : '订单不存在' }}</view>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import {
  cancelOrder,
  confirmReceive,
  createAfterSale,
  getAfterSaleEligible,
  getOrder,
  getShop,
  ORDER_STATUS,
} from '../../api/index'

const shopPlaceholder = 'https://picsum.photos/id/20/100/100'
const order = ref(null)
const shop = ref({})
const afterSale = ref({ eligible: false, reason: '', days: 7, deadline: '' })
const loading = ref(false)
const busy = ref(false)
const showAfterSaleForm = ref(false)
const afterForm = reactive({
  type: 'refund',
  amount: '',
  reason: '',
})
let orderId = 0

const canCancel = computed(() => {
  const s = order.value?.status
  return s === 'pending' || s === 'confirmed'
})
const canConfirm = computed(() => order.value?.status === 'shipped')
const canReview = computed(() => order.value?.status === 'completed')
const canViewReview = computed(() => order.value?.status === 'reviewed')
const canAfterSale = computed(() => !!afterSale.value?.eligible)
const payAmount = computed(() => {
  const o = order.value
  if (!o) return '0'
  const n = Number(o.pay_amount > 0 ? o.pay_amount : o.total_amount)
  return Number.isFinite(n) ? n.toFixed(2) : '0'
})

const shopDisplayName = computed(() =>
  shop.value.name || order.value?.shop_name || (order.value?.shop_id ? `店铺 #${order.value.shop_id}` : '店铺'),
)

const shopLine = computed(() => {
  const parts = [shop.value.province, shop.value.city, shop.value.district].filter(Boolean)
  if (parts.length) return parts.join(' · ')
  return shop.value.category ? `主营 ${shop.value.category}` : '官方店铺'
})

function statusText(s) {
  return ORDER_STATUS[s] || s
}

onLoad((q) => {
  orderId = Number(q.id || 0)
  load()
})

async function loadShop(shopID) {
  if (!shopID) {
    shop.value = {}
    return
  }
  try {
    const res = await getShop(shopID)
    shop.value = res || { name: order.value?.shop_name || `店铺 #${shopID}` }
  } catch {
    shop.value = { name: order.value?.shop_name || `店铺 #${shopID}` }
  }
}

async function loadEligible() {
  if (!orderId) return
  try {
    const res = await getAfterSaleEligible(orderId)
    afterSale.value = {
      eligible: !!res?.eligible,
      reason: res?.reason || '',
      days: res?.days ?? 7,
      deadline: res?.deadline || '',
    }
  } catch {
    afterSale.value = { eligible: false, reason: '', days: 7, deadline: '' }
  }
}

async function load() {
  if (!orderId) return
  loading.value = true
  showAfterSaleForm.value = false
  try {
    const res = await getOrder(orderId)
    order.value = res || null
    if (order.value?.shop_id) {
      await loadShop(order.value.shop_id)
    } else {
      shop.value = { name: order.value?.shop_name || '店铺' }
    }
    await loadEligible()
  } catch {
    order.value = null
    shop.value = {}
    afterSale.value = { eligible: false, reason: '', days: 7, deadline: '' }
  } finally {
    loading.value = false
  }
}

function goShop() {
  const id = order.value?.shop_id
  if (!id) return
  uni.navigateTo({ url: `/pages/shop/detail?id=${id}` })
}

function openAfterSale() {
  afterForm.type = 'refund'
  afterForm.amount = payAmount.value
  afterForm.reason = ''
  showAfterSaleForm.value = true
  uni.pageScrollTo({ scrollTop: 0, duration: 0 })
}

async function onSubmitAfterSale() {
  const reason = (afterForm.reason || '').trim()
  if (!reason) {
    uni.showToast({ title: '请填写原因', icon: 'none' })
    return
  }
  const amount = Number(afterForm.amount)
  if (!Number.isFinite(amount) || amount <= 0) {
    uni.showToast({ title: '请输入有效金额', icon: 'none' })
    return
  }
  busy.value = true
  try {
    await createAfterSale(orderId, {
      type: afterForm.type,
      reason,
      amount,
    })
    uni.showToast({ title: '已提交', icon: 'success' })
    showAfterSaleForm.value = false
    load()
  } catch {
    /* handled */
  } finally {
    busy.value = false
  }
}

async function onCancel() {
  const ok = await new Promise((resolve) => {
    uni.showModal({
      title: '取消订单',
      content: '确认取消该订单？',
      success: (r) => resolve(r.confirm),
    })
  })
  if (!ok) return
  busy.value = true
  try {
    await cancelOrder(orderId)
    uni.showToast({ title: '已取消', icon: 'success' })
    load()
  } catch {
    /* handled */
  } finally {
    busy.value = false
  }
}

async function onConfirm() {
  const ok = await new Promise((resolve) => {
    uni.showModal({
      title: '确认收货',
      content: '确认已收到商品？',
      success: (r) => resolve(r.confirm),
    })
  })
  if (!ok) return
  busy.value = true
  try {
    await confirmReceive(orderId)
    uni.showToast({ title: '已确认收货', icon: 'success' })
    load()
  } catch {
    /* handled */
  } finally {
    busy.value = false
  }
}

function goReview() {
  uni.navigateTo({ url: `/pages/order/review?id=${orderId}` })
}

function goViewReview() {
  uni.navigateTo({ url: `/pages/order/review-view?id=${orderId}` })
}
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; padding-bottom: calc(48rpx + env(safe-area-inset-bottom)); }
.card {
  background: #fff; border-radius: 24rpx; padding: 28rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.shop-card {
  display: flex; align-items: center; gap: 16rpx; padding: 24rpx 28rpx;
}
.shop-logo { width: 72rpx; height: 72rpx; border-radius: 14rpx; background: #f4f4f5; flex-shrink: 0; }
.shop-meta { flex: 1; min-width: 0; }
.shop-name {
  display: block; font-size: 28rpx; font-weight: 700; color: #18181b;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.shop-sub { display: block; margin-top: 4rpx; font-size: 22rpx; color: #a1a1aa; }
.shop-arrow { font-size: 24rpx; color: #a1a1aa; flex-shrink: 0; }
.row { display: flex; justify-content: space-between; padding: 12rpx 0; font-size: 26rpx; }
.row.col { flex-direction: column; align-items: flex-start; gap: 8rpx; }
.addr { color: #3f3f46; line-height: 1.4; }
.label { color: #71717a; }
.gold { color: #c8a876; }
.price { color: #d83636; font-weight: 700; }
.sec { font-weight: 600; font-size: 28rpx; display: block; margin-bottom: 16rpx; }
.item { padding: 12rpx 0; border-top: 1rpx solid #f5f5f5; }
.name { display: block; font-size: 28rpx; }
.sub { color: #71717a; font-size: 22rpx; }
.actions { display: flex; flex-direction: column; gap: 16rpx; margin-top: 8rpx; }
.btn {
  border-radius: 999rpx; height: 80rpx; line-height: 80rpx; font-size: 28rpx; margin: 0;
}
.btn.outline { background: #fff; color: #d83636; border: 2rpx solid #d83636; }
.btn.primary { background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; border: none; }
.empty { text-align: center; padding: 120rpx; color: #71717a; }
.after-sale { padding-bottom: 160rpx; }
.hint { display: block; font-size: 22rpx; color: #a1a1aa; margin: -8rpx 0 20rpx; }
.form .field { margin-bottom: 20rpx; }
.flabel { display: block; font-size: 24rpx; color: #71717a; margin-bottom: 10rpx; }
.seg { display: flex; gap: 12rpx; }
.seg-item {
  flex: 1; text-align: center; padding: 16rpx 0; border-radius: 12rpx;
  background: #f4f4f5; font-size: 26rpx; color: #52525b;
}
.seg-item.on { background: rgba(200,168,118,.18); color: #9a7b4f; font-weight: 600; }
.input, .textarea {
  width: 100%; box-sizing: border-box; background: #fafafa; border-radius: 12rpx;
  padding: 18rpx 20rpx; font-size: 26rpx; color: #18181b;
}
.amount-row {
  display: flex; align-items: center; gap: 8rpx;
  background: #fafafa; border-radius: 12rpx; padding: 0 20rpx;
  min-height: 88rpx; box-sizing: border-box;
}
.amount-prefix {
  font-size: 32rpx; font-weight: 700; color: #d83636; flex-shrink: 0; line-height: 1;
}
.amount-input {
  flex: 1; min-width: 0; width: auto; height: 88rpx; line-height: 88rpx;
  padding: 0; background: transparent; font-size: 32rpx; font-weight: 600;
  color: #18181b;
}
.textarea { min-height: 160rpx; width: 100%; }
.foot-actions {
  position: fixed; left: 0; right: 0; bottom: 0; z-index: 20;
  display: flex; gap: 16rpx; padding: 16rpx 32rpx calc(16rpx + env(safe-area-inset-bottom));
  background: rgba(255,255,255,.96); box-shadow: 0 -4rpx 24rpx rgba(0,0,0,.06);
}
.foot-actions .btn { flex: 1; }
</style>
