<template>
  <view class="page">
    <view class="bg" />
    <view class="bg-orb o1" />
    <view class="bg-orb o2" />
    <view class="bg-orb o3" />

    <!-- 周边氛围粒子（固定数量，节奏错开） -->
    <view class="fx">
      <view v-for="i in 10" :key="'p' + i" class="dot" :class="'d' + i" />
      <view v-for="i in 5" :key="'c' + i" class="coin" :class="'c' + i">✦</view>
    </view>

    <view class="header">
      <text class="title">{{ activity?.title || '九宫格抽奖' }}</text>
      <text class="subtitle">幸运九宫格 · 今日超级开奖</text>
      <view class="meta">
        <view class="chip">耗 {{ activity?.cost_points ?? '-' }} 积分</view>
        <view class="chip">
          今日剩
          {{ activity?.today_remaining === -1 ? '不限' : (activity?.today_remaining ?? '-') }}
        </view>
        <view class="chip chip-gold">积分 {{ points }}</view>
      </view>
    </view>

    <view class="stage" :class="{ drawing: drawing }">
      <view class="aura a1" />
      <view class="aura a2" />

      <view class="board-wrap">
        <view class="board-glow" />
        <view class="board">
          <view
            v-for="slot in 9"
            :key="slot - 1"
            class="cell"
            :class="{
              on: highlight === slot - 1,
              win: celebrating && highlight === slot - 1,
            }"
          >
            <view v-if="slot - 1 === 4" class="draw" :class="{ busy: drawing }" @tap="onDraw">
              <view class="draw-ring r-slow" />
              <view class="draw-ring r-fast" />
              <view class="draw-core">
                <text class="draw-txt">{{ drawing ? '抽奖中' : '立即抽奖' }}</text>
              </view>
            </view>

            <view v-else class="prize">
              <view class="prize-shine" />
              <view class="prize-icon" :class="iconClass(prizeAt(slot - 1))">
                <image
                  v-if="prizeAt(slot - 1).cover_url"
                  class="prize-img"
                  :src="prizeAt(slot - 1).cover_url"
                  mode="aspectFit"
                />
                <text v-else class="prize-emoji">{{ prizeIcon(prizeAt(slot - 1)) }}</text>
              </view>
              <text class="prize-name">{{ prizeAt(slot - 1).name || '神秘奖' }}</text>
              <text
                v-if="prizeAt(slot - 1).prize_type === 'points'"
                class="prize-val"
              >+{{ prizeAt(slot - 1).points_amount }}</text>
              <text
                v-else-if="prizeAt(slot - 1).prize_type === 'physical'"
                class="prize-val"
              >实物</text>
              <text v-else class="prize-val mute">—</text>
            </view>
          </view>
        </view>
      </view>

      <view class="base">
        <view class="base-ring b1" />
        <view class="base-ring b2" />
        <view class="base-ring b3" />
      </view>
    </view>

    <text class="hint">结果以服务端为准</text>

    <view v-if="celebrating" class="mask" @tap="closeCelebrate">
      <view class="modal" @tap.stop>
        <text class="modal-emoji">{{ winTip.emoji }}</text>
        <text class="modal-title">{{ winTip.title }}</text>
        <text class="modal-sub">{{ winTip.sub }}</text>
        <view class="modal-btn" @tap="closeCelebrate">知道了</view>
      </view>
    </view>

    <view v-if="addrVisible" class="mask" @tap="closeAddr">
      <view class="modal addr-modal" @tap.stop>
        <text class="modal-title">填写收货地址</text>
        <text class="modal-sub">奖品：{{ claimPrizeName }}</text>
        <view v-if="!addresses.length" class="addr-empty">
          <text>暂无收货地址</text>
          <view class="modal-btn" @tap="goAddAddress">去添加</view>
        </view>
        <scroll-view v-else scroll-y class="addr-list">
          <view
            v-for="a in addresses"
            :key="a.id"
            class="addr-item"
            :class="{ on: selectedAddressId === a.id }"
            @tap="selectedAddressId = a.id"
          >
            <text class="addr-name">{{ a.receiver_name }} {{ a.receiver_phone }}</text>
            <text class="addr-detail">{{ fullAddress(a) }}</text>
          </view>
        </scroll-view>
        <view class="addr-actions">
          <view class="modal-btn ghost" @tap="closeAddr">稍后填写</view>
          <view class="modal-btn" :class="{ disabled: !selectedAddressId || claiming }" @tap="submitAddress">
            {{ claiming ? '提交中…' : '确认地址' }}
          </view>
        </view>
      </view>
    </view>

    <view class="panel">
      <text class="panel-title">我的战绩</text>
      <view v-for="r in records" :key="r.id" class="row rec-card">
        <view class="rec-main">
          <text class="row-name">{{ r.prize_name }}</text>
          <text class="rec-tag">{{ recordTag(r) }}</text>
          <text v-if="r.prize_type === 'physical' && r.fulfill_status === 'shipped'" class="rec-ship">
            物流：{{ r.ship_company || '' }} {{ r.ship_no }}
          </text>
          <text v-else-if="r.prize_type === 'physical' && r.fulfill_status === 'pending'" class="rec-ship">
            待发货 · {{ r.receiver_name }} {{ r.receiver_phone }}
          </text>
        </view>
        <view class="rec-right">
          <text class="row-time">{{ r.created_at }}</text>
          <text
            v-if="r.prize_type === 'physical' && r.fulfill_status === 'need_address'"
            class="rec-claim"
            @tap="openClaim(r)"
          >填地址</text>
        </view>
      </view>
      <view v-if="!records.length" class="empty">暂无记录</view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { drawLottery, claimLotteryAddress, getLotteryActivity, getUserPoints, listAddresses, listLotteryRecords } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const RING_ORDER = [0, 1, 2, 5, 8, 7, 6, 3]

const activity = ref(null)
const points = ref(0)
const records = ref([])
const drawing = ref(false)
const highlight = ref(-1)
const celebrating = ref(false)
const winTip = ref({ emoji: '✨', title: '', sub: '' })

const addrVisible = ref(false)
const addresses = ref([])
const selectedAddressId = ref(0)
const claimRecordId = ref(0)
const claimPrizeName = ref('')
const claiming = ref(false)

const prizeMap = computed(() => {
  const map = {}
  for (let i = 0; i < 9; i++) {
    map[i] = { slot: i, name: '', prize_type: 'thanks', points_amount: 0 }
  }
  for (const p of activity.value?.prizes || []) {
    if (p.slot >= 0 && p.slot < 9) map[p.slot] = p
  }
  return map
})

function prizeAt(slot) {
  return prizeMap.value[slot] || { name: '', prize_type: 'thanks', points_amount: 0 }
}

function prizeIcon(p) {
  if (!p) return '🎁'
  if (p.prize_type === 'physical') return '📦'
  if (p.prize_type !== 'points') return '🎁'
  const n = Number(p.points_amount) || 0
  if (n >= 200) return '💎'
  if (n >= 100) return '🏆'
  if (n >= 50) return '📱'
  if (n >= 20) return '🧧'
  return '🪙'
}

function iconClass(p) {
  if (!p) return 'tone-soft'
  if (p.prize_type === 'physical') return 'tone-hot'
  if (p.prize_type !== 'points') return 'tone-soft'
  const n = Number(p.points_amount) || 0
  if (n >= 100) return 'tone-hot'
  if (n >= 20) return 'tone-mid'
  return 'tone-soft'
}

function fullAddress(a) {
  if (!a) return ''
  return `${a.province || ''}${a.city || ''}${a.district || ''}${a.detail || ''}`
}

function recordTag(r) {
  if (r.prize_type === 'points') return `+${r.points_amount} 积分`
  if (r.prize_type === 'physical') {
    if (r.fulfill_status === 'need_address') return '待填地址'
    if (r.fulfill_status === 'pending') return '待发货'
    if (r.fulfill_status === 'shipped') return '已发货'
    return '实物'
  }
  return '谢谢参与'
}

async function loadAddresses() {
  try {
    const res = await listAddresses()
    const list = res?.list || res?.data?.list || []
    addresses.value = list
    const def = list.find((a) => Number(a.is_default) === 1)
    selectedAddressId.value = def?.id || list[0]?.id || 0
  } catch (_) {
    addresses.value = []
    selectedAddressId.value = 0
  }
}

async function openClaim(row) {
  claimRecordId.value = row.id
  claimPrizeName.value = row.prize_name || '实物奖品'
  celebrating.value = false
  await loadAddresses()
  addrVisible.value = true
}

function closeAddr() {
  addrVisible.value = false
}

function closeCelebrate() {
  celebrating.value = false
}

function goAddAddress() {
  addrVisible.value = false
  uni.navigateTo({ url: '/pages/address/list' })
}

async function submitAddress() {
  if (!selectedAddressId.value || claiming.value) return
  claiming.value = true
  try {
    await claimLotteryAddress(claimRecordId.value, selectedAddressId.value)
    uni.showToast({ title: '地址已提交', icon: 'success' })
    addrVisible.value = false
    claimRecordId.value = 0
    await load()
  } catch (e) {
    uni.showToast({ title: e.message || '提交失败', icon: 'none' })
  } finally {
    claiming.value = false
  }
}

async function load() {
  try {
    const res = await getLotteryActivity()
    activity.value = res?.data || res
  } catch (e) {
    activity.value = null
    uni.showToast({ title: e.message || '加载失败', icon: 'none' })
  }
  if (isLoggedIn()) {
    try {
      const p = await getUserPoints()
      points.value = p?.points ?? p?.data?.points ?? 0
    } catch (_) {}
    try {
      const r = await listLotteryRecords({ page: 1, page_size: 20 })
      records.value = r?.list || []
    } catch (_) {
      records.value = []
    }
  }
}

function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms))
}

async function animateTo(targetSlot) {
  const target = Number(targetSlot)
  const landOnCenter = target === 4
  const endIdx = landOnCenter ? RING_ORDER.length - 1 : RING_ORDER.indexOf(target)
  const loops = 3
  const total = loops * RING_ORDER.length + Math.max(endIdx, 0)

  for (let i = 0; i <= total; i++) {
    highlight.value = RING_ORDER[i % RING_ORDER.length]
    const t = i / Math.max(total, 1)
    await sleep(36 + Math.pow(t, 2.2) * 200)
  }

  highlight.value = landOnCenter ? 4 : target
  for (let k = 0; k < 3; k++) {
    await sleep(100)
    highlight.value = -1
    await sleep(80)
    highlight.value = landOnCenter ? 4 : target
  }
}

async function onDraw() {
  if (drawing.value) return
  if (!isLoggedIn()) {
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/lottery/index'),
    })
    return
  }
  celebrating.value = false
  addrVisible.value = false
  drawing.value = true
  try {
    const res = await drawLottery()
    const data = res?.data || res
    await animateTo(Number(data.slot))
    if (data.prize_type === 'points') {
      winTip.value = {
        emoji: '💎',
        title: `+${data.points_amount} 积分`,
        sub: data.prize_name || '恭喜中奖',
      }
      celebrating.value = true
    } else if (data.prize_type === 'physical') {
      winTip.value = {
        emoji: '📦',
        title: data.prize_name || '恭喜获得实物',
        sub: '请选择收货地址',
      }
      celebrating.value = true
      claimRecordId.value = data.record_id
      claimPrizeName.value = data.prize_name || '实物奖品'
      setTimeout(async () => {
        celebrating.value = false
        await loadAddresses()
        addrVisible.value = true
      }, 900)
    } else {
      winTip.value = {
        emoji: '🎁',
        title: data.prize_name || '谢谢参与',
        sub: '再接再厉，下次更好运',
      }
      celebrating.value = true
    }
    await load()
  } catch (e) {
    uni.showToast({ title: e.message || '抽奖失败', icon: 'none' })
  } finally {
    drawing.value = false
  }
}

onShow(() => load())
</script>

<style scoped>
.page {
  position: relative;
  min-height: 100vh;
  padding: 48rpx 32rpx 72rpx;
  color: #f3e8ff;
  background: #5b21b6;
  overflow: hidden;
}

/* 中紫主调：既不发黑也不过浅 */
.bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse 90% 50% at 50% -8%, rgba(233, 213, 255, 0.55) 0%, transparent 55%),
    radial-gradient(ellipse 70% 45% at 100% 40%, rgba(192, 132, 252, 0.45), transparent 50%),
    radial-gradient(ellipse 60% 40% at 0% 70%, rgba(244, 114, 182, 0.28), transparent 50%),
    linear-gradient(180deg, #7c3aed 0%, #6d28d9 42%, #5b21b6 100%);
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
  filter: blur(2rpx);
  animation: orb-drift 7s ease-in-out infinite;
}
.bg-orb.o1 {
  width: 280rpx;
  height: 280rpx;
  top: 120rpx;
  right: -60rpx;
  background: radial-gradient(circle, rgba(244, 114, 182, 0.45), transparent 70%);
}
.bg-orb.o2 {
  width: 320rpx;
  height: 320rpx;
  bottom: 180rpx;
  left: -100rpx;
  background: radial-gradient(circle, rgba(129, 140, 248, 0.4), transparent 70%);
  animation-delay: -2.5s;
}
.bg-orb.o3 {
  width: 200rpx;
  height: 200rpx;
  top: 42%;
  left: 36%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.35), transparent 70%);
  animation-delay: -4s;
}
@keyframes orb-drift {
  0%, 100% { transform: translate3d(0, 0, 0) scale(1); opacity: 0.7; }
  50% { transform: translate3d(18rpx, -24rpx, 0) scale(1.08); opacity: 1; }
}

.fx {
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
}
.dot {
  position: absolute;
  width: 8rpx;
  height: 8rpx;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 0 12rpx #f0abfc, 0 0 22rpx #a855f7;
  animation: twinkle 2s ease-in-out infinite;
}
.d1 { top: 12%; left: 10%; }
.d2 { top: 18%; left: 78%; animation-delay: .2s; width: 6rpx; height: 6rpx; }
.d3 { top: 28%; left: 22%; animation-delay: .4s; }
.d4 { top: 34%; left: 88%; animation-delay: .6s; width: 5rpx; height: 5rpx; }
.d5 { top: 48%; left: 8%; animation-delay: .8s; }
.d6 { top: 52%; left: 92%; animation-delay: 1s; }
.d7 { top: 66%; left: 16%; animation-delay: .3s; width: 6rpx; height: 6rpx; }
.d8 { top: 72%; left: 70%; animation-delay: .7s; }
.d9 { top: 84%; left: 40%; animation-delay: 1.1s; }
.d10 { top: 22%; left: 48%; animation-delay: .5s; width: 4rpx; height: 4rpx; }
@keyframes twinkle {
  0%, 100% { opacity: 0.2; transform: scale(0.6); }
  50% { opacity: 1; transform: scale(1.35); }
}

.coin {
  position: absolute;
  font-size: 28rpx;
  color: #fbbf24;
  text-shadow: 0 0 12rpx rgba(251, 191, 36, 0.8);
  animation: float-up 5s ease-in-out infinite;
  opacity: 0;
}
.c1 { left: 12%; animation-delay: 0s; }
.c2 { left: 28%; animation-delay: 1s; font-size: 22rpx; }
.c3 { left: 55%; animation-delay: 2s; }
.c4 { left: 72%; animation-delay: .6s; font-size: 24rpx; }
.c5 { left: 86%; animation-delay: 3s; font-size: 20rpx; }
@keyframes float-up {
  0% { top: 70%; opacity: 0; transform: translateY(0) rotate(0deg) scale(0.8); }
  20% { opacity: 0.9; }
  100% { top: 12%; opacity: 0; transform: translateY(-20rpx) rotate(40deg) scale(1.2); }
}

.header {
  position: relative;
  z-index: 2;
  text-align: center;
  margin-bottom: 28rpx;
}
.title {
  display: block;
  font-size: 46rpx;
  font-weight: 800;
  letter-spacing: 2rpx;
  color: #fff;
  text-shadow:
    0 2rpx 0 rgba(255, 255, 255, 0.25),
    0 0 28rpx rgba(233, 213, 255, 0.85);
  animation: title-glow 2.6s ease-in-out infinite;
}
@keyframes title-glow {
  0%, 100% { filter: brightness(1); }
  50% { filter: brightness(1.12); }
}
.subtitle {
  display: block;
  margin-top: 12rpx;
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.78);
  letter-spacing: 3rpx;
}
.meta {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-top: 24rpx;
}
.chip {
  padding: 10rpx 22rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  font-weight: 600;
  color: #f5e8ff;
  background: rgba(255, 255, 255, 0.16);
  border: 1rpx solid rgba(233, 213, 255, 0.4);
  box-shadow: 0 6rpx 16rpx rgba(76, 29, 149, 0.25);
}
.chip-gold {
  color: #78350f;
  background: linear-gradient(135deg, #fde68a, #f59e0b);
  border-color: transparent;
  box-shadow: 0 6rpx 18rpx rgba(245, 158, 11, 0.35);
}

.stage {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.aura {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
  border: 2rpx solid rgba(255, 255, 255, 0.35);
  box-shadow: 0 0 24rpx rgba(232, 121, 249, 0.35);
}
.aura.a1 {
  width: 108%;
  height: 420rpx;
  top: 4%;
  animation: aura-spin 12s linear infinite;
  opacity: 0.55;
}
.aura.a2 {
  width: 92%;
  height: 360rpx;
  top: 10%;
  border-color: rgba(34, 211, 238, 0.35);
  animation: aura-spin 9s linear infinite reverse;
  opacity: 0.45;
}
@keyframes aura-spin {
  to { transform: rotate(360deg); }
}

.stage.drawing .draw-ring.r-fast {
  animation-duration: 0.7s;
}
.stage.drawing .board-glow {
  opacity: 1;
  animation-duration: 1s;
}
.stage.drawing .aura {
  opacity: 0.85;
}

.board-wrap {
  position: relative;
  width: 100%;
  max-width: 640rpx;
}
.board-glow {
  position: absolute;
  inset: -18rpx;
  border-radius: 40rpx;
  pointer-events: none;
  background: radial-gradient(circle at 50% 40%, rgba(244, 114, 182, 0.45), transparent 65%);
  animation: glow-pulse 2.4s ease-in-out infinite;
  opacity: 0.75;
}
@keyframes glow-pulse {
  0%, 100% { transform: scale(0.96); opacity: 0.55; }
  50% { transform: scale(1.04); opacity: 1; }
}

.board {
  position: relative;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14rpx;
  padding: 20rpx;
  border-radius: 32rpx;
  background: linear-gradient(165deg, rgba(167, 139, 250, 0.95), rgba(124, 58, 237, 0.92) 48%, rgba(109, 40, 217, 0.95));
  border: 3rpx solid transparent;
  background-clip: padding-box;
  box-shadow:
    0 0 0 3rpx rgba(165, 243, 252, 0.7),
    0 0 36rpx rgba(192, 132, 252, 0.55),
    0 20rpx 40rpx rgba(76, 29, 149, 0.35),
    inset 0 2rpx 0 rgba(255, 255, 255, 0.35);
  animation: board-border 4s linear infinite;
}
@keyframes board-border {
  0%, 100% { box-shadow: 0 0 0 3rpx rgba(165, 243, 252, 0.7), 0 0 36rpx rgba(192, 132, 252, 0.55), 0 20rpx 40rpx rgba(76, 29, 149, 0.35), inset 0 2rpx 0 rgba(255, 255, 255, 0.35); }
  50% { box-shadow: 0 0 0 3rpx rgba(240, 171, 252, 0.9), 0 0 48rpx rgba(232, 121, 249, 0.65), 0 20rpx 40rpx rgba(76, 29, 149, 0.35), inset 0 2rpx 0 rgba(255, 255, 255, 0.35); }
}

.cell {
  aspect-ratio: 1;
  min-width: 0;
}
.cell.on .prize,
.cell.on .draw-core {
  border-color: #fbbf24;
  box-shadow:
    0 0 0 3rpx rgba(251, 191, 36, 0.75),
    0 0 32rpx rgba(250, 204, 21, 0.7),
    inset 0 0 18rpx rgba(253, 224, 71, 0.25);
  transform: scale(1.06);
}
.cell.win .prize,
.cell.win .draw-core {
  animation: win-pulse 0.45s ease-in-out 3;
}
@keyframes win-pulse {
  0%, 100% { transform: scale(1.06); }
  50% { transform: scale(1.12); }
}

.prize {
  position: relative;
  height: 100%;
  overflow: hidden;
  border-radius: 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 8rpx 4rpx 6rpx;
  background: linear-gradient(160deg, rgba(255, 255, 255, 0.28), rgba(237, 233, 254, 0.22));
  border: 2rpx solid rgba(233, 213, 255, 0.45);
  box-shadow:
    0 8rpx 18rpx rgba(76, 29, 149, 0.25),
    inset 0 1rpx 0 rgba(255, 255, 255, 0.35);
  transition: transform 0.15s ease, box-shadow 0.15s ease, border-color 0.15s ease;
}
.prize-shine {
  position: absolute;
  top: 0;
  left: -60%;
  width: 40%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.45), transparent);
  transform: skewX(-18deg);
  animation: shine-sweep 3.2s ease-in-out infinite;
  pointer-events: none;
}
.cell:nth-child(2) .prize-shine { animation-delay: .4s; }
.cell:nth-child(3) .prize-shine { animation-delay: .8s; }
.cell:nth-child(4) .prize-shine { animation-delay: 1.2s; }
.cell:nth-child(6) .prize-shine { animation-delay: 1.6s; }
.cell:nth-child(7) .prize-shine { animation-delay: .2s; }
.cell:nth-child(8) .prize-shine { animation-delay: .6s; }
.cell:nth-child(9) .prize-shine { animation-delay: 1s; }
@keyframes shine-sweep {
  0% { left: -60%; opacity: 0; }
  25% { opacity: 1; }
  55%, 100% { left: 130%; opacity: 0; }
}

.prize-icon {
  width: 100rpx;
  height: 100rpx;
  margin-bottom: 6rpx;
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: icon-bob 2.4s ease-in-out infinite;
}
.cell:nth-child(odd) .prize-icon { animation-delay: .35s; }
@keyframes icon-bob {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6rpx); }
}
.prize-icon.tone-soft { background: rgba(255, 255, 255, 0.18); }
.prize-icon.tone-mid { background: rgba(244, 114, 182, 0.28); }
.prize-icon.tone-hot { background: rgba(251, 191, 36, 0.28); }
.prize-emoji {
  font-size: 64rpx;
  line-height: 1;
  filter: drop-shadow(0 6rpx 10rpx rgba(76, 29, 149, 0.35));
}
.prize-img {
  width: 88rpx;
  height: 88rpx;
}
.prize-name {
  max-width: 100%;
  font-size: 24rpx;
  font-weight: 700;
  color: #fff;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-shadow: 0 1rpx 4rpx rgba(76, 29, 149, 0.45);
}
.prize-val {
  margin-top: 2rpx;
  font-size: 22rpx;
  font-weight: 800;
  color: #fde68a;
}
.prize-val.mute {
  color: rgba(255, 255, 255, 0.4);
  font-weight: 500;
}

.draw {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.draw.busy { pointer-events: none; }
.draw-ring {
  position: absolute;
  border-radius: 50%;
  border: 3rpx solid transparent;
}
.draw-ring.r-slow {
  width: 92%;
  height: 92%;
  border-top-color: #e879f9;
  border-right-color: rgba(34, 211, 238, 0.8);
  animation: spin 3.2s linear infinite;
  filter: drop-shadow(0 0 8rpx rgba(232, 121, 249, 0.7));
}
.draw-ring.r-fast {
  width: 76%;
  height: 76%;
  border-bottom-color: #22d3ee;
  border-left-color: #c084fc;
  animation: spin 1.8s linear infinite reverse;
  filter: drop-shadow(0 0 8rpx rgba(34, 211, 238, 0.6));
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.draw-core {
  position: relative;
  z-index: 1;
  width: 66%;
  height: 66%;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(circle at 35% 28%, #f5d0fe, #c084fc 42%, #9333ea 100%);
  border: 2rpx solid rgba(255, 255, 255, 0.75);
  box-shadow:
    0 0 24rpx rgba(232, 121, 249, 0.65),
    0 10rpx 22rpx rgba(147, 51, 234, 0.4),
    inset 0 2rpx 0 rgba(255, 255, 255, 0.55);
  transition: transform 0.15s ease, box-shadow 0.15s ease, border-color 0.15s ease;
  animation: core-breathe 2.2s ease-in-out infinite;
}
@keyframes core-breathe {
  0%, 100% { transform: scale(1); box-shadow: 0 0 24rpx rgba(232, 121, 249, 0.65), 0 10rpx 22rpx rgba(147, 51, 234, 0.4), inset 0 2rpx 0 rgba(255, 255, 255, 0.55); }
  50% { transform: scale(1.05); box-shadow: 0 0 36rpx rgba(232, 121, 249, 0.9), 0 10rpx 22rpx rgba(147, 51, 234, 0.45), inset 0 2rpx 0 rgba(255, 255, 255, 0.55); }
}
.draw:active .draw-core {
  transform: scale(0.95);
  animation: none;
}
.draw-txt {
  font-size: 22rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 1rpx;
  text-align: center;
  line-height: 1.25;
  padding: 0 6rpx;
  text-shadow: 0 2rpx 8rpx rgba(91, 33, 182, 0.45);
}

.base {
  position: relative;
  width: 82%;
  height: 72rpx;
  margin-top: 2rpx;
  transform: scaleY(0.42);
}
.base-ring {
  position: absolute;
  left: 50%;
  top: 50%;
  border-radius: 50%;
  transform: translate(-50%, -50%);
  border: 3rpx solid rgba(192, 132, 252, 0.55);
  box-shadow: 0 0 20rpx rgba(168, 85, 247, 0.45);
  animation: base-glow 2.2s ease-in-out infinite;
}
.base-ring.b1 { width: 100%; height: 100%; }
.base-ring.b2 {
  width: 74%;
  height: 74%;
  border-color: rgba(34, 211, 238, 0.5);
  animation-delay: .35s;
}
.base-ring.b3 {
  width: 48%;
  height: 48%;
  border-color: rgba(244, 114, 182, 0.55);
  animation-delay: .7s;
}
@keyframes base-glow {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 1; }
}

.hint {
  position: relative;
  z-index: 2;
  display: block;
  margin: 28rpx 0 8rpx;
  text-align: center;
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.45);
}

.mask {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(49, 10, 101, 0.55);
  padding: 48rpx;
}
.modal {
  width: 100%;
  max-width: 560rpx;
  padding: 48rpx 40rpx 36rpx;
  border-radius: 28rpx;
  text-align: center;
  background: linear-gradient(165deg, #ede9fe, #c4b5fd);
  border: 2rpx solid rgba(192, 132, 252, 0.55);
  box-shadow: 0 24rpx 48rpx rgba(76, 29, 149, 0.35);
  animation: modal-in 0.35s cubic-bezier(0.22, 1, 0.36, 1);
}
@keyframes modal-in {
  from { opacity: 0; transform: translateY(24rpx) scale(0.96); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.modal-emoji { display: block; font-size: 64rpx; margin-bottom: 8rpx; }
.modal-title {
  display: block;
  font-size: 36rpx;
  font-weight: 800;
  color: #5b21b6;
}
.modal-sub {
  display: block;
  margin-top: 12rpx;
  font-size: 24rpx;
  color: rgba(91, 33, 182, 0.7);
}
.modal-btn {
  margin-top: 36rpx;
  height: 80rpx;
  line-height: 80rpx;
  border-radius: 999rpx;
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(135deg, #c084fc, #9333ea);
  box-shadow: 0 10rpx 24rpx rgba(147, 51, 234, 0.4);
}
.modal-btn.ghost {
  background: rgba(91, 33, 182, 0.15);
  color: #5b21b6;
  box-shadow: none;
}
.modal-btn.disabled { opacity: 0.55; }

.addr-modal { text-align: left; max-height: 78vh; }
.addr-empty {
  margin-top: 24rpx;
  text-align: center;
  color: rgba(91, 33, 182, 0.7);
  font-size: 26rpx;
}
.addr-list {
  margin-top: 20rpx;
  max-height: 420rpx;
}
.addr-item {
  padding: 20rpx;
  margin-bottom: 12rpx;
  border-radius: 16rpx;
  background: rgba(255, 255, 255, 0.65);
  border: 2rpx solid transparent;
}
.addr-item.on {
  border-color: #a855f7;
  background: #f5e8ff;
}
.addr-name {
  display: block;
  font-size: 26rpx;
  font-weight: 700;
  color: #5b21b6;
}
.addr-detail {
  display: block;
  margin-top: 6rpx;
  font-size: 22rpx;
  color: rgba(91, 33, 182, 0.7);
}
.addr-actions {
  display: flex;
  gap: 16rpx;
  margin-top: 20rpx;
}
.addr-actions .modal-btn {
  flex: 1;
  margin-top: 0;
  text-align: center;
}

.panel {
  position: relative;
  z-index: 2;
  margin-top: 32rpx;
  padding: 28rpx 24rpx;
  border-radius: 24rpx;
  background: rgba(255, 255, 255, 0.16);
  border: 1rpx solid rgba(233, 213, 255, 0.35);
  box-shadow: 0 12rpx 28rpx rgba(76, 29, 149, 0.2);
}
.panel-title {
  display: block;
  margin-bottom: 8rpx;
  font-size: 28rpx;
  font-weight: 800;
  color: #fff;
}
.row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20rpx 0;
  border-bottom: 1rpx solid rgba(233, 213, 255, 0.18);
}
.row:last-child { border-bottom: none; }
.rec-main { flex: 1; min-width: 0; padding-right: 16rpx; }
.rec-right { display: flex; flex-direction: column; align-items: flex-end; gap: 8rpx; }
.row-name {
  display: block;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.95);
  font-weight: 700;
}
.rec-tag {
  display: inline-block;
  margin-top: 8rpx;
  padding: 2rpx 12rpx;
  border-radius: 999rpx;
  font-size: 20rpx;
  color: #fef3c7;
  background: rgba(251, 191, 36, 0.25);
}
.rec-ship {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: rgba(233, 213, 255, 0.75);
  line-height: 1.4;
}
.rec-claim {
  font-size: 22rpx;
  font-weight: 700;
  color: #fde68a;
  padding: 6rpx 14rpx;
  border-radius: 999rpx;
  background: rgba(251, 191, 36, 0.2);
}
.row-time {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.5);
}
.empty {
  padding: 40rpx 0 16rpx;
  text-align: center;
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.45);
}
</style>
