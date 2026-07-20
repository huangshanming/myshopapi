<template>
  <view class="page">
    <view class="hero">
      <text class="hero-label">我的积分</text>
      <text class="hero-points">{{ points }}</text>
      <text class="hero-sub">完成任务后记得领取奖励</text>
    </view>

    <view class="card checkin" :class="{ done: checkinDone }" @tap="onCheckinCard">
      <view>
        <text class="c-title">每日签到</text>
        <text class="c-desc">{{ checkinDesc }}</text>
      </view>
      <text class="c-btn" :class="checkinBtnClass">{{ checkinBtnText }}</text>
    </view>

    <view v-for="t in otherTasks" :key="t.code" class="task">
      <view class="t-main">
        <text class="t-title">{{ t.title }}</text>
        <text class="t-desc">{{ t.description }}</text>
        <text class="t-meta">
          进度 {{ t.progress }}/{{ t.target_count }}
          · 奖励 {{ t.reward_points }} 积分
          <text v-if="t.period === 'daily' && t.daily_limit"> · 每日上限 {{ t.daily_limit }}</text>
        </text>
      </view>
      <text
        class="t-btn"
        :class="t.status"
        @tap="onAction(t)"
      >{{ btnText(t) }}</text>
    </view>

    <view v-if="!loading && !tasks.length" class="empty">暂无任务</view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { checkinTask, claimTask, getUserPoints, listUserTasks } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const points = ref(0)
const tasks = ref([])
const loading = ref(false)

const checkinTaskItem = computed(() => tasks.value.find((t) => t.code === 'daily_checkin') || null)
const otherTasks = computed(() => tasks.value.filter((t) => t.code !== 'daily_checkin'))
const checkinStatus = computed(() => checkinTaskItem.value?.status || 'ongoing')
const checkinDone = computed(() => checkinStatus.value === 'claimed' || checkinStatus.value === 'claimable')

const checkinBtnText = computed(() => {
  if (checkinStatus.value === 'claimed') return '已签到'
  if (checkinStatus.value === 'claimable') return '领取'
  return '签到'
})

const checkinBtnClass = computed(() => {
  if (checkinStatus.value === 'claimed') return 'disabled'
  if (checkinStatus.value === 'claimable') return 'claimable'
  return ''
})

const checkinDesc = computed(() => {
  if (checkinStatus.value === 'claimed') return '今日已签到并领取，明天再来吧'
  if (checkinStatus.value === 'claimable') return '签到完成，点击领取积分'
  return '点此签到，完成后可领取积分'
})

function btnText(t) {
  if (t.status === 'claimable') return '领取'
  if (t.status === 'claimed') return '已领'
  if (t.code === 'publish_article') return '去发文'
  if (t.code === 'browse_products') return '去逛逛'
  if (t.code === 'place_order') return '去下单'
  return '进行中'
}

function onCheckinCard() {
  if (checkinStatus.value === 'claimed') {
    uni.showToast({ title: '今日已签到', icon: 'none' })
    return
  }
  if (checkinStatus.value === 'claimable') {
    doClaim('daily_checkin')
    return
  }
  doCheckin()
}

function onAction(t) {
  if (t.status === 'claimable') {
    doClaim(t.code)
    return
  }
  if (t.status === 'claimed') return
  if (t.code === 'publish_article') {
    uni.navigateTo({ url: '/pages/community/publish' })
    return
  }
  if (t.code === 'browse_products' || t.code === 'place_order') {
    uni.switchTab({ url: '/pages/index/index' })
    return
  }
  if (t.code === 'comment_article' || t.code === 'like_article' || t.code === 'favorite_article') {
    uni.navigateTo({ url: '/pages/community/list' })
  }
}

async function doCheckin() {
  if (!ensureLogin()) return
  if (checkinDone.value) {
    uni.showToast({ title: '今日已签到', icon: 'none' })
    return
  }
  try {
    await checkinTask()
    uni.showToast({ title: '签到成功，可领取', icon: 'none' })
    load()
  } catch (e) {
    uni.showToast({ title: e.message || '签到失败', icon: 'none' })
  }
}

async function doClaim(code) {
  if (!ensureLogin()) return
  try {
    const res = await claimTask(code)
    points.value = res?.points ?? points.value
    uni.showToast({ title: '领取成功', icon: 'none' })
    load()
  } catch (e) {
    uni.showToast({ title: e.message || '领取失败', icon: 'none' })
  }
}

function ensureLogin() {
  if (isLoggedIn()) return true
  uni.navigateTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/task/index') })
  return false
}

async function load() {
  if (!ensureLogin()) return
  loading.value = true
  try {
    const [p, t] = await Promise.all([getUserPoints(), listUserTasks()])
    points.value = p?.points || 0
    tasks.value = t?.list || []
  } catch {
    tasks.value = []
  } finally {
    loading.value = false
  }
}

onShow(load)
</script>

<style scoped>
.page { padding: 24rpx 24rpx 60rpx; min-height: 100vh; background: #fafafa; }
.hero {
  background: linear-gradient(135deg, #e8d5b5, #c8a876);
  border-radius: 24rpx; padding: 40rpx 32rpx; margin-bottom: 24rpx; color: #fff;
}
.hero-label { font-size: 24rpx; opacity: .9; }
.hero-points { display: block; font-size: 72rpx; font-weight: 700; margin: 8rpx 0; }
.hero-sub { font-size: 22rpx; opacity: .85; }
.checkin {
  display: flex; justify-content: space-between; align-items: center;
  background: #fff; border-radius: 16rpx; padding: 28rpx; margin-bottom: 20rpx;
}
.checkin.done { opacity: 0.92; }
.c-title { display: block; font-size: 30rpx; font-weight: 600; color: #18181b; }
.c-desc { display: block; font-size: 22rpx; color: #a1a1aa; margin-top: 6rpx; }
.c-btn {
  background: #c8a876; color: #fff; font-size: 26rpx; padding: 12rpx 28rpx; border-radius: 32rpx;
}
.c-btn.claimable { background: #c8a876; color: #fff; }
.c-btn.disabled {
  background: #e4e4e7; color: #a1a1aa;
}
.task {
  display: flex; gap: 16rpx; align-items: center;
  background: #fff; border-radius: 16rpx; padding: 28rpx; margin-bottom: 16rpx;
}
.t-main { flex: 1; min-width: 0; }
.t-title { display: block; font-size: 28rpx; font-weight: 600; color: #18181b; }
.t-desc { display: block; font-size: 22rpx; color: #71717a; margin-top: 6rpx; }
.t-meta { display: block; font-size: 20rpx; color: #a1a1aa; margin-top: 10rpx; }
.t-btn {
  flex-shrink: 0; font-size: 24rpx; padding: 12rpx 24rpx; border-radius: 32rpx;
  background: #f4f4f5; color: #71717a;
}
.t-btn.claimable { background: #c8a876; color: #fff; }
.t-btn.claimed { background: #f4f4f5; color: #a1a1aa; }
.empty { text-align: center; color: #a1a1aa; padding: 80rpx 0; font-size: 26rpx; }
</style>
