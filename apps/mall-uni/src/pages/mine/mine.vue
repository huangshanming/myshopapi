<template>
  <view class="page">
    <view class="hero">
      <view class="hero-glow" />
      <view class="profile" @tap="!user && goLogin()">
        <view class="avatar">
          <image v-if="user?.avatar" class="avatar-img" :src="user.avatar" mode="aspectFill" />
          <text v-else>{{ avatarText }}</text>
        </view>
        <view class="info">
          <template v-if="user">
            <text class="name">{{ user.nickname || user.mobile }}</text>
            <text class="sub">欢迎回来</text>
          </template>
          <template v-else>
            <text class="name">未登录</text>
            <text class="sub">登录后同步订单与资产</text>
          </template>
        </view>
        <view v-if="!user" class="login-btn" @tap.stop="goLogin">去登录</view>
      </view>

      <view v-if="user" class="asset">
        <view class="asset-col" @tap="goWalletLogs">
          <text class="a-label">可用余额</text>
          <view class="a-amount-row">
            <text class="a-yen">¥</text>
            <text class="a-num">{{ wallet.balance ?? '0.00' }}</text>
          </view>
        </view>
        <view class="asset-col" @tap="goWalletLogs">
          <text class="a-label">冻结金额</text>
          <view class="a-amount-row">
            <text class="a-yen">¥</text>
            <text class="a-num">{{ wallet.frozen_balance ?? '0.00' }}</text>
          </view>
        </view>
        <view class="asset-col" @tap="goCoupons">
          <text class="a-label">优惠券</text>
          <view class="a-amount-row">
            <text class="a-num">{{ couponCount }}</text>
            <text class="a-unit">张</text>
          </view>
        </view>
        <view class="asset-col" @tap="goTasks">
          <text class="a-label">积分</text>
          <view class="a-amount-row">
            <text class="a-num">{{ points }}</text>
          </view>
        </view>
      </view>
    </view>

    <view class="panel order-panel">
      <view class="panel-head" @tap="goOrders()">
        <text class="panel-title">我的订单</text>
        <text class="panel-more">全部 ›</text>
      </view>
      <view class="order-grid">
        <view
          v-for="item in orderEntries"
          :key="item.status"
          class="order-entry"
          @tap="goOrders(item.status)"
        >
          <view class="icon-wrap">
            <image class="order-ico" :src="item.icon" mode="aspectFit" />
            <view v-if="badge(item.status) > 0" class="badge">
              {{ badge(item.status) > 99 ? '99+' : badge(item.status) }}
            </view>
          </view>
          <text class="entry-label">{{ item.label }}</text>
        </view>
      </view>
    </view>

    <view class="panel menu-panel">
      <view
        v-for="m in menuItems"
        :key="m.key"
        class="menu-item"
        :class="{ danger: m.danger, hide: m.needLogin && !user }"
        @tap="m.onTap"
      >
        <view class="menu-left">
          <view class="menu-ico-wrap" :style="{ background: m.bg }">
            <image class="menu-ico" :src="m.icon" mode="aspectFit" />
          </view>
          <text class="menu-text">{{ m.label }}</text>
        </view>
        <text class="menu-arrow">›</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getOrderStatusCounts, getUserPoints, getUserWallet, listMyCoupons, reportTaskEvent } from '../../api/index'
import { clearAuth, getUser, isLoggedIn } from '../../stores/user'

const stroke = '#b8956a'
const strokeMuted = '#8b7355'
const strokeDanger = '#c45c5c'

function svgIcon(paths, color = stroke) {
  const xml = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="${color}" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">${paths}</svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(xml)}`
}

const ICO = {
  pay: svgIcon('<rect x="2.5" y="5" width="19" height="14" rx="2.5"/><path d="M2.5 10h19"/><path d="M7 15h3"/>'),
  box: svgIcon('<path d="M12 2.8 20.5 7.5v9L12 21.2 3.5 16.5v-9L12 2.8z"/><path d="M12 12.2V21"/><path d="M3.7 7.7 12 12.2l8.3-4.5"/>'),
  truck: svgIcon('<path d="M3 7.5h10.5V16H3z"/><path d="M13.5 10.5H18l2.5 3V16h-7"/><circle cx="7" cy="17.5" r="1.6"/><circle cx="17" cy="17.5" r="1.6"/>'),
  review: svgIcon('<path d="M12 3.2l2.1 4.3 4.7.7-3.4 3.3.8 4.7L12 14l-4.2 2.2.8-4.7-3.4-3.3 4.7-.7L12 3.2z"/>'),
  refund: svgIcon('<path d="M4 12a8 8 0 0 1 13.5-5.8"/><path d="M17.5 3.5v3.2H14"/><path d="M20 12a8 8 0 0 1-13.5 5.8"/><path d="M6.5 20.5v-3.2H10"/>'),
  coupon: svgIcon('<path d="M4 8.5A2.5 2.5 0 0 0 6.5 6h11A2.5 2.5 0 0 0 20 8.5v1.2a2 2 0 1 1 0 4v1.8A2.5 2.5 0 0 1 17.5 18h-11A2.5 2.5 0 0 1 4 15.5v-1.8a2 2 0 1 1 0-4V8.5z"/><path d="M12 6.5v11"/>', strokeMuted),
  gift: svgIcon('<rect x="4" y="9" width="16" height="11" rx="1.5"/><path d="M4 13h16"/><path d="M12 9v11"/><path d="M12 9c-2.2-3.5-5.5-2.2-5.5 0S12 9 12 9z"/><path d="M12 9c2.2-3.5 5.5-2.2 5.5 0S12 9 12 9z"/>', strokeMuted),
  heart: svgIcon('<path d="M12 20s-7-4.4-7-9.2A3.8 3.8 0 0 1 12 8.2a3.8 3.8 0 0 1 7 2.6C19 15.6 12 20 12 20z"/>', strokeMuted),
  like: svgIcon('<path d="M8 11v9H5.5A1.5 1.5 0 0 1 4 18.5v-6A1.5 1.5 0 0 1 5.5 11H8zm0 0 3.2-6.2A1.8 1.8 0 0 1 12.8 4a1.8 1.8 0 0 1 1.8 2.1L14 11h4.2a2 2 0 0 1 2 2.3l-1.1 6A2 2 0 0 1 17.1 21H8"/>', strokeMuted),
  pin: svgIcon('<path d="M12 21s6-5.2 6-10a6 6 0 1 0-12 0c0 4.8 6 10 6 10z"/><circle cx="12" cy="11" r="2.2"/>', strokeMuted),
  bill: svgIcon('<path d="M7 3.5h10A1.5 1.5 0 0 1 18.5 5v16L12 17.5 5.5 21V5A1.5 1.5 0 0 1 7 3.5z"/><path d="M9 8h6M9 11.5h6"/>', strokeMuted),
  exit: svgIcon('<path d="M10 4.5H6.5A2 2 0 0 0 4.5 6.5v11a2 2 0 0 0 2 2H10"/><path d="M14 12H8.5"/><path d="M14 12l-3-3m3 3-3 3"/><path d="M14 4.5h3.5a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H14"/>', strokeDanger),
}

const user = ref(null)
const wallet = ref({ balance: 0, frozen_balance: 0 })
const couponCount = ref(0)
const points = ref(0)
const counts = ref({})

const orderEntries = [
  { status: 'pending', label: '待付款', icon: ICO.pay },
  { status: 'confirmed', label: '待发货', icon: ICO.box },
  { status: 'shipped', label: '待收货', icon: ICO.truck },
  { status: 'completed', label: '待评价', icon: ICO.review },
  { status: 'after_sale', label: '退款/售后', icon: ICO.refund },
]

const avatarText = computed(() => {
  const n = user.value?.nickname || user.value?.mobile || '?'
  return String(n).slice(0, 1)
})

function badge(status) {
  return Number(counts.value?.[status]) || 0
}

function goLogin() {
  uni.navigateTo({ url: '/pages/login/login' })
}

function goOrders(status = '') {
  if (!isLoggedIn()) {
    goLogin()
    return
  }
  const q = status ? `?status=${status}` : ''
  uni.navigateTo({ url: `/pages/order/list${q}` })
}

function goWalletLogs() {
  if (!isLoggedIn()) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/wallet/logs' })
}

function goAddresses() {
  if (!isLoggedIn()) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/address/list' })
}

function goFavorites() {
  if (!isLoggedIn()) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/favorite/list' })
}

function goCoupons() {
  if (!isLoggedIn()) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/coupon/mine' })
}

function goCouponCenter() {
  uni.navigateTo({ url: '/pages/coupon/center' })
}

function goLikes() {
  if (!isLoggedIn()) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/community/liked' })
}

function goTasks() {
  if (!isLoggedIn()) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/task/index' })
}

function goMyNotes() {
  if (!isLoggedIn()) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/community/mine' })
}

function logout() {
  clearAuth()
  user.value = null
  wallet.value = { balance: 0, frozen_balance: 0 }
  couponCount.value = 0
  points.value = 0
  counts.value = {}
  uni.showToast({ title: '已退出', icon: 'none' })
}

const menuItems = computed(() => [
  { key: 'tasks', label: '任务中心', icon: ICO.gift, bg: '#faf6ef', needLogin: true, onTap: goTasks },
  { key: 'notes', label: '我的笔记', icon: ICO.like, bg: '#f7f3ec', needLogin: true, onTap: goMyNotes },
  { key: 'coupons', label: '我的优惠券', icon: ICO.coupon, bg: '#faf6ef', needLogin: true, onTap: goCoupons },
  { key: 'center', label: '领券中心', icon: ICO.gift, bg: '#f7f3ec', needLogin: false, onTap: goCouponCenter },
  { key: 'fav', label: '我的收藏', icon: ICO.heart, bg: '#faf6ef', needLogin: true, onTap: goFavorites },
  { key: 'like', label: '我的点赞', icon: ICO.like, bg: '#f7f3ec', needLogin: true, onTap: goLikes },
  { key: 'addr', label: '收货地址', icon: ICO.pin, bg: '#faf6ef', needLogin: true, onTap: goAddresses },
  { key: 'wallet', label: '余额明细', icon: ICO.bill, bg: '#f7f3ec', needLogin: true, onTap: goWalletLogs },
  { key: 'out', label: '退出登录', icon: ICO.exit, bg: '#fdf2f2', needLogin: true, danger: true, onTap: logout },
])

async function loadWallet() {
  if (!isLoggedIn()) {
    wallet.value = { balance: 0, frozen_balance: 0 }
    return
  }
  try {
    const res = await getUserWallet()
    wallet.value = res || { balance: 0, frozen_balance: 0 }
  } catch {
    wallet.value = { balance: 0, frozen_balance: 0 }
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
  const u = user.value
  if (u?.nickname && u?.avatar) {
    reportTaskEvent({ task_code: 'first_profile', ref_type: 'profile', ref_id: u.id }).catch(() => {})
  }
}

async function loadCouponCount() {
  if (!isLoggedIn()) {
    couponCount.value = 0
    return
  }
  try {
    const res = await listMyCoupons({ page: 1, page_size: 1, status: 'unused' })
    couponCount.value = Number(res?.total) || 0
  } catch {
    couponCount.value = 0
  }
}

async function loadOrderCounts() {
  if (!isLoggedIn()) {
    counts.value = {}
    return
  }
  try {
    counts.value = (await getOrderStatusCounts()) || {}
  } catch {
    counts.value = {}
  }
}

onShow(() => {
  user.value = getUser()
  loadWallet()
  loadCouponCount()
  loadPoints()
  loadOrderCounts()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding-bottom: 48rpx;
  background:
    radial-gradient(ellipse 80% 40% at 10% -10%, rgba(212, 184, 144, 0.35), transparent 55%),
    radial-gradient(ellipse 60% 35% at 100% 5%, rgba(191, 164, 114, 0.2), transparent 50%),
    linear-gradient(180deg, #f3ebe0 0%, #f7f3ec 28%, #faf8f5 100%);
}

.hero {
  position: relative;
  margin: 0 0 20rpx;
  padding: 40rpx 32rpx 28rpx;
  overflow: hidden;
}
.hero-glow {
  position: absolute;
  inset: 0;
  background: linear-gradient(145deg, rgba(191, 164, 114, 0.22), rgba(255, 255, 255, 0.05) 55%);
  pointer-events: none;
}

.profile {
  position: relative;
  display: flex;
  align-items: center;
  gap: 24rpx;
  z-index: 1;
}
.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background: linear-gradient(145deg, #c8a876, #a88755);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48rpx;
  font-weight: 700;
  box-shadow: 0 10rpx 28rpx rgba(168, 135, 85, 0.28);
  border: 4rpx solid rgba(255, 255, 255, 0.65);
  overflow: hidden;
  flex-shrink: 0;
}
.avatar-img { width: 100%; height: 100%; }
.name {
  display: block;
  font-size: 38rpx;
  font-weight: 700;
  color: #2c2416;
  letter-spacing: 0.5rpx;
}
.sub {
  display: block;
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #8a7d6b;
}
.login-btn {
  margin-left: auto;
  padding: 12rpx 28rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890);
  color: #fff;
  font-size: 24rpx;
  font-weight: 600;
  box-shadow: 0 6rpx 16rpx rgba(191, 164, 114, 0.35);
}

.asset {
  position: relative;
  z-index: 1;
  margin-top: 28rpx;
  display: flex;
  align-items: flex-start;
  padding: 28rpx 8rpx;
  border-radius: 24rpx;
  background: linear-gradient(135deg, #fffdf9 0%, #f8f0e4 100%);
  border: 1rpx solid rgba(200, 168, 118, 0.28);
  box-shadow: 0 8rpx 24rpx rgba(168, 135, 85, 0.08);
}
.asset-col {
  flex: 1;
  min-width: 0;
  padding: 0 12rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  position: relative;
}
.asset-col + .asset-col::before {
  content: '';
  position: absolute;
  left: 0;
  top: 8rpx;
  bottom: 8rpx;
  width: 1rpx;
  background: rgba(184, 149, 106, 0.2);
}
.a-label {
  display: block;
  width: 100%;
  font-size: 22rpx;
  color: #9a8b76;
  margin-bottom: 14rpx;
  line-height: 1;
}
.a-amount-row {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 4rpx;
  min-height: 44rpx;
}
.a-yen {
  font-size: 24rpx;
  font-weight: 600;
  color: #8b7355;
}
.a-num {
  font-size: 32rpx;
  font-weight: 700;
  color: #2c2416;
  line-height: 1;
  letter-spacing: 0.5rpx;
}
.a-unit {
  font-size: 22rpx;
  color: #8b7355;
  margin-left: 2rpx;
}

.panel {
  margin: 0 28rpx 20rpx;
  padding: 28rpx 24rpx 24rpx;
  background: rgba(255, 255, 255, 0.92);
  border-radius: 28rpx;
  box-shadow: 0 8rpx 28rpx rgba(140, 110, 70, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.8);
}
.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 8rpx 24rpx;
}
.panel-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #2c2416;
}
.panel-more { font-size: 24rpx; color: #a89880; }

.order-grid {
  display: flex;
  justify-content: space-between;
  padding: 4rpx 0 8rpx;
}
.order-entry {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14rpx;
}
.icon-wrap {
  position: relative;
  width: 80rpx;
  height: 80rpx;
  border-radius: 24rpx;
  background: linear-gradient(160deg, #fbf7f0, #f3ebe0);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: inset 0 1rpx 0 rgba(255, 255, 255, 0.9);
}
.order-ico { width: 44rpx; height: 44rpx; }
.badge {
  position: absolute;
  top: -10rpx;
  right: -12rpx;
  min-width: 30rpx;
  height: 30rpx;
  padding: 0 8rpx;
  border-radius: 30rpx;
  background: #d83636;
  color: #fff;
  font-size: 18rpx;
  line-height: 30rpx;
  text-align: center;
  box-sizing: border-box;
  border: 2rpx solid #fff;
}
.entry-label {
  font-size: 22rpx;
  color: #5c5348;
}

.menu-panel { padding: 8rpx 12rpx; }
.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 22rpx 16rpx;
  border-radius: 20rpx;
}
.menu-item.hide { display: none; }
.menu-item:active { background: rgba(191, 164, 114, 0.08); }
.menu-left { display: flex; align-items: center; gap: 20rpx; }
.menu-ico-wrap {
  width: 64rpx;
  height: 64rpx;
  border-radius: 18rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
.menu-ico { width: 34rpx; height: 34rpx; }
.menu-text {
  font-size: 28rpx;
  color: #2c2416;
}
.menu-item.danger .menu-text { color: #b4534b; }
.menu-arrow {
  font-size: 34rpx;
  color: #c8b89a;
  line-height: 1;
}
</style>
