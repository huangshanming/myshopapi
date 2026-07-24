<template>
  <view class="page">
    <!-- 优惠券通栏 -->
    <view class="coupon-bar gold-gradient" @tap="goCouponCenter">
      <text class="coupon-text">{{ couponBarText }}</text>
      <text class="coupon-btn">去领取</text>
    </view>

    <!-- 顶栏 -->
    <view class="header">
      <view class="loc" @tap="goChooseCity">
        <text class="loc-icon">📍</text>
        <text>{{ locCity }}</text>
      </view>
      <view class="search" @tap="toast('搜索即将开放')">
        <text class="search-icon">🔍</text>
        <text class="search-ph">搜索店铺 / 好物 / 种草笔记</text>
      </view>
      <view class="bell-wrap" @tap="goMessages">
        <text class="bell" :class="{ shake: unreadCount > 0 }">🔔</text>
        <view v-if="unreadCount > 0" class="badge">{{ unreadCount > 99 ? '99+' : unreadCount }}</view>
      </view>
    </view>

    <!-- Banner -->
    <view class="banner-wrap">
      <swiper class="banner" circular autoplay interval="3500" indicator-dots indicator-active-color="#C8A876">
        <swiper-item v-for="(b, i) in banners" :key="b.id || i">
          <image class="banner-img" :src="b.image_url || b" mode="aspectFill" @tap="onBannerTap(b)" />
        </swiper-item>
      </swiper>
      <view class="corner-tag"><text class="corner-text sale-gradient">限时特惠</text></view>
    </view>

    <!-- 分类入口：两行 4×2 -->
    <view class="cat-grid section-card">
      <view v-for="(c, i) in categoryEntries" :key="c.id || c.name + '-' + i" class="cat-item" @tap="onCategory(c)">
        <view class="cat-icon" :style="{ background: c.bg }" :class="{ bounce: c.bounce }">
          <image class="cat-img" :src="c.icon" mode="aspectFit" />
        </view>
        <text class="cat-name">{{ c.name }}</text>
      </view>
    </view>

    <!-- 限时秒杀 -->
    <view class="section-card seckill">
      <view class="sec-head">
        <view class="sec-left">
          <text class="sec-title">⚡ 限时秒杀</text>
          <view class="countdown">
            <text class="cd-label">距结束</text>
            <text class="count-num">{{ cd.h }}</text>:
            <text class="count-num">{{ cd.m }}</text>:
            <text class="count-num">{{ cd.s }}</text>
          </view>
        </view>
        <text class="sec-more" @tap="goSeckillList">更多 ›</text>
      </view>
      <scroll-view scroll-x :show-scrollbar="false" class="scroll-hide">
        <view class="seckill-row">
          <view
            v-for="(p, i) in seckillItems"
            :key="p.id || i"
            class="seckill-item"
            @tap="p.product_id && goSeckillDetail(p)"
          >
            <image class="seckill-img" :src="p.img || placeholder" mode="aspectFill" />
            <text class="line-2 seckill-name">{{ p.name }}</text>
            <view class="price-row">
              <text class="price">¥{{ p.price }}</text>
              <text class="price-old">¥{{ p.old }}</text>
            </view>
            <text class="stock-left">仅剩{{ p.left }}件</text>
          </view>
          <view v-if="!seckillItems.length" class="seckill-empty">
            <text class="sub">本场暂无秒杀商品</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 积分商城入口：热闹舞台风，配色贴合首页金棕 -->
    <view class="points-entry" @tap="goPointsMall">
      <view class="pe-glow" />
      <view class="pe-ring" />
      <view class="pe-beam pe-beam-a" />
      <view class="pe-beam pe-beam-b" />
      <view class="pe-floor" />

      <view class="pe-gift">
        <view class="pe-gift-lid" />
        <view class="pe-gift-body" />
        <view class="pe-gift-bow" />
      </view>
      <view class="pe-ticket">
        <view class="pe-ticket-cut" />
      </view>
      <view class="pe-coin c1">¥</view>
      <view class="pe-coin c2">¥</view>
      <view class="pe-coin c3">会</view>
      <view class="pe-spark s1" />
      <view class="pe-spark s2" />
      <view class="pe-spark s3" />
      <view class="pe-spark s4" />
      <view class="pe-bolt" />

      <text class="pe-chip pe-chip-a">新人专享</text>
      <text class="pe-chip pe-chip-b">热门兑换</text>
      <text class="pe-chip pe-chip-c">每日福利</text>

      <view class="pe-main">
        <text class="pe-title">积分商城</text>
        <view class="pe-sub">积分翻倍兑 · 福利限时抢</view>
      </view>
      <view class="pe-cta">立即逛逛</view>
    </view>

    <!-- 优惠购入口 -->
    <view class="cps-entry" @tap="goCpsList">
      <view class="cps-entry-text">
        <text class="cps-entry-title">优惠购</text>
        <text class="cps-entry-sub">外卖打车省更多 · 复制链接去第三方购买</text>
      </view>
      <text class="cps-entry-cta">去看看</text>
    </view>

    <!-- 头部品牌商户 -->
    <view class="block">
      <view class="block-head">
        <text class="sec-title">👑 头部品牌商户</text>
        <text class="sec-more" @tap="goShopList('brand_shop')">全部品牌店 ›</text>
      </view>
      <scroll-view scroll-x :show-scrollbar="false" class="scroll-hide">
        <view class="brand-row">
          <view
            v-for="s in brandShops"
            :key="s.id || s.name"
            class="brand-card"
            @tap="goShopDetail(s.id)"
          >
            <image class="brand-cover" :src="s.img" mode="aspectFill" />
            <view class="brand-body">
              <text class="line-1 brand-name">{{ s.name }}</text>
              <view class="brand-meta">
                <text class="gold">{{ s.tag }}</text>
                <text v-if="s.city" class="sub">{{ s.city }}</text>
                <text v-if="s.paid" class="sub">推广</text>
              </view>
            </view>
          </view>
          <view v-if="!brandShops.length" class="seckill-empty">
            <text class="sub">暂无品牌商户</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 种草社区 -->
    <view class="block">
      <view class="block-head">
        <text class="sec-title">✍️ 好物种草社区</text>
        <text class="sec-more" @tap="goCommunityList">更多精选 ›</text>
      </view>
      <scroll-view scroll-x :show-scrollbar="false" class="scroll-hide">
        <view class="note-row">
          <view
            v-for="(n, i) in notes"
            :key="n.id || i"
            class="note-card"
            :class="'w-' + n.w"
            @tap="goArticle(n.id)"
          >
            <image class="note-cover" :src="n.imgs[0] || placeholder" mode="aspectFill" />
            <view class="note-body">
              <text class="note-label">{{ n.label }}</text>
              <text class="line-2 note-title">{{ n.title }}</text>
              <view class="note-foot">
                <text>阅 {{ n.reads }}</text>
                <text>❤ {{ n.likes }}</text>
              </view>
            </view>
          </view>
          <view v-if="!notes.length" class="seckill-empty">
            <text class="sub">暂无种草内容</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 主题集市 -->
    <view class="section-card theme">
      <text class="sec-title theme-title">▦ 主题好物集市</text>
      <view class="theme-grid">
        <view
          v-for="t in themes"
          :key="t.slot_id || t.position"
          class="theme-item"
          @tap="onThemeTap(t)"
        >
          <image class="theme-img" :src="t.cover_url || placeholder" mode="aspectFill" />
          <view class="theme-mask">
            <text class="theme-name">{{ t.name }}</text>
            <text class="theme-desc">{{ t.desc }}</text>
          </view>
          <text v-if="t.paid" class="theme-paid">推广</text>
        </view>
      </view>
    </view>

    <!-- 优质商户 -->
    <view class="block">
      <view class="block-head">
        <text class="sec-title">🏪 优质入驻商户</text>
        <text class="sec-more" @tap="goShopList('quality_shop')">全部店铺 ›</text>
      </view>
      <scroll-view scroll-x :show-scrollbar="false" class="scroll-hide">
        <view class="shop-row">
          <view
            v-for="s in shops"
            :key="s.id || s.name"
            class="shop-card"
            @tap="goShopDetail(s.id)"
          >
            <image class="shop-cover" :src="s.img" mode="aspectFill" />
            <view class="shop-body">
              <text class="line-1 brand-name">{{ s.name }}</text>
              <text class="sub shop-score">{{ s.category || '优选商户' }}{{ s.city ? ` · ${s.city}` : '' }}{{ s.paid ? ' · 推广' : '' }}</text>
              <view class="shop-promo">
                <text class="sub line-1">{{ s.desc || '品质好店' }}</text>
              </view>
            </view>
          </view>
          <view v-if="!shops.length" class="seckill-empty">
            <text class="sub">暂无优质商户</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 本地商家 -->
    <view class="block">
      <view class="block-head">
        <text class="sec-title">📍 本地商家</text>
        <text class="sec-more" @tap="goLocalShops">更多 ›</text>
      </view>
      <scroll-view scroll-x :show-scrollbar="false" class="scroll-hide">
        <view class="shop-row">
          <view
            v-for="s in localShops"
            :key="s.id"
            class="shop-card"
            @tap="goShopDetail(s.id)"
          >
            <image class="shop-cover" :src="s.img" mode="aspectFill" />
            <view class="shop-body">
              <text class="line-1 brand-name">{{ s.name }}</text>
              <text class="sub shop-score">
                {{ s.city || '本地' }}{{ s.distance_km != null ? ` · ${formatDist(s.distance_km)}` : '' }}
              </text>
              <view class="shop-promo">
                <text class="sub line-1">{{ s.address || s.desc || '附近好店' }}</text>
              </view>
            </view>
          </view>
          <view v-if="!localShops.length" class="seckill-empty">
            <text class="sub">{{ localHint }}</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 销量榜 -->
    <view class="section-card rank">
      <view class="block-head no-pad">
        <text class="sec-title">🏆 今日必买销量榜</text>
        <text class="sec-more" @tap="goSalesRank">完整榜单 ›</text>
      </view>
      <view v-if="!rankList.length" class="rank-empty">
        <text class="sub">暂无上榜商品</text>
      </view>
      <view
        v-for="(r, i) in rankList"
        :key="r.id"
        class="rank-item"
        @tap="goDetail(r.id)"
      >
        <view class="rank-no" :class="'r' + (i + 1)">{{ i + 1 }}</view>
        <image class="rank-img" :src="r.main_image || placeholder" mode="aspectFill" />
        <view class="rank-info">
          <text class="line-1 rank-name">{{ rankDisplayName(r) }}</text>
          <view class="rank-foot">
            <text class="price">¥{{ r.sale_price }}</text>
            <text class="sub">{{ rankSoldText(r) }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 热销商品双列 -->
    <view class="section-card goods">
      <view class="block-head no-pad">
        <text class="sec-title">🛍 全城商户热销好物</text>
        <text class="sec-more" @tap="goCategoryPage()">全部商品 ›</text>
      </view>
      <view class="goods-grid">
        <view v-for="p in products" :key="p.id" class="goods-card" @tap="goDetail(p.id)">
          <image class="goods-img" :src="p.main_image || placeholder" mode="aspectFill" />
          <text class="line-2 goods-name">{{ p.name }}</text>
          <view class="price-row">
            <text class="price goods-price">¥{{ p.sale_price }}</text>
            <text v-if="p.market_price" class="price-old">¥{{ p.market_price }}</text>
          </view>
          <text class="sub sold">月销{{ p.sold_count || 0 }}+</text>
        </view>
      </view>
      <view class="load-more">
        <text v-if="loading">加载中...</text>
        <text v-else-if="finished">没有更多了</text>
        <text v-else @tap="loadMore">上拉加载更多</text>
      </view>
    </view>

    <!-- 任务中心悬浮入口：礼盒风格 + 画面动效 -->
    <view class="task-fab" @tap="goTaskCenter">
      <view class="task-fab-glow" />
      <image class="task-fab-art" src="/static/task-fab.png" mode="aspectFit" />
      <view class="task-fab-coin c1">$</view>
      <view class="task-fab-coin c2">$</view>
      <view class="task-fab-spark s1" />
      <view class="task-fab-spark s2" />
      <view class="task-fab-spark s3" />
      <view class="task-fab-spark s4" />
      <text class="task-fab-badge">任务</text>
    </view>
  </view>
</template>

<script setup>
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { onReachBottom, onShow } from '@dcloudio/uni-app'
import {
  getNotificationUnreadCount, getSeckillCurrent, listAddresses, listArticles, listBanners, listCategories,
  listCouponCenter, listHomeSlots, listLocalShops, listProducts, listSalesRank, listThemeTiles,
} from '../../api/index'
import { isLoggedIn } from '../../stores/user'
import { applyAddressCityIfNeeded, getCity, getCoords, hasCoords } from '../../stores/location'

const placeholder = 'https://picsum.photos/id/96/400/400'
const locCity = ref(getCity())
const fallbackBanners = [
  { id: 'f1', image_url: '/static/banner/banner-1.png', link_type: 'none' },
  { id: 'f2', image_url: '/static/banner/banner-2.png', link_type: 'none' },
  { id: 'f3', image_url: '/static/banner/banner-3.png', link_type: 'none' },
]
const banners = ref([...fallbackBanners])

const CAT_ICON = {
  shop: '/static/cat/shop.png',
  fashion: '/static/cat/fashion.png',
  snack: '/static/cat/snack.png',
  fresh: '/static/cat/fresh.png',
  beauty: '/static/cat/beauty.png',
  digital: '/static/cat/digital.png',
  community: '/static/cat/community.png',
  seckill: '/static/cat/seckill.png',
  coupon: '/static/cat/coupon.png',
  points: '/static/cat/points.png',
  orders: '/static/cat/orders.png',
  messages: '/static/cat/messages.png',
  more: '/static/cat/more.png',
  brand: '/static/cat/brand.png',
}
const apiCatIconPool = [
  CAT_ICON.fashion, CAT_ICON.snack, CAT_ICON.fresh, CAT_ICON.beauty, CAT_ICON.digital,
]

function resolveApiCatIcon(name, index) {
  const n = String(name || '')
  if (/服|鞋|包|衣|穿搭|时装/.test(n)) return CAT_ICON.fashion
  if (/食|零食|休|糕|饮|茶|酒/.test(n)) return CAT_ICON.snack
  if (/生鲜|烘焙|果|菜|肉|海鲜/.test(n)) return CAT_ICON.fresh
  if (/美妆|妆|护肤|化妆/.test(n)) return CAT_ICON.beauty
  if (/数码|电|手机|家电/.test(n)) return CAT_ICON.digital
  return apiCatIconPool[index % apiCatIconPool.length]
}

const fixedLead = {
  name: '全部商户',
  icon: CAT_ICON.shop,
  bg: 'rgba(230,213,188,.45)',
}
const fixedFeatures = [
  { name: '本地商家', icon: CAT_ICON.shop, bg: '#E8F0FE', local: true },
  { name: '优惠购', icon: CAT_ICON.coupon, bg: '#FEF3C7', cps: true },
  { name: '积分商城', icon: CAT_ICON.points, bg: '#FFF1E0', pointsMall: true },
]
const fixedTail = [
  { name: '种草社区', icon: CAT_ICON.community, bg: 'rgba(230,213,188,.45)' },
  { name: '限时秒杀', icon: CAT_ICON.seckill, bg: '#FEF2F2', bounce: true },
]
const fillEntries = [
  { name: '领券中心', icon: CAT_ICON.coupon, bg: '#FEF3C7', coupon: true },
  { name: '我的订单', icon: CAT_ICON.orders, bg: '#DBEAFE', orders: true },
  { name: '消息中心', icon: CAT_ICON.messages, bg: '#FCE7F3', messages: true },
  { name: '更多分类', icon: CAT_ICON.more, bg: '#F4F4F5', more: true },
  { name: '品牌好店', icon: CAT_ICON.brand, bg: 'rgba(230,213,188,.45)', brand: true },
]
const apiCats = ref([])
// 10 格：全部商户 + 4 分类 + 本地商家/优惠购/积分 + 种草社区 + 限时秒杀
const categoryEntries = computed(() => {
  const colors = ['#E8F8EF', '#E8F0FE', '#FDE8F2', '#FFF6E0', '#EEF0FF']
  const mid = apiCats.value.slice(0, 4).map((c, i) => ({
    name: c.name,
    icon: resolveApiCatIcon(c.name, i),
    id: c.id,
    bg: colors[i % colors.length],
  }))
  let fi = 0
  while (mid.length < 4 && fi < fillEntries.length) {
    mid.push(fillEntries[fi++])
  }
  return [fixedLead, ...mid, ...fixedFeatures, ...fixedTail]
})

const seckillItems = ref([])
const brandShops = ref([])
const shops = ref([])
const localShops = ref([])
const localHint = ref('加载中...')
const notes = ref([])

const themes = ref([])
const couponBarText = ref('领券中心有好券，点击立即领取')
const unreadCount = ref(0)
let unreadTimer

const rankList = ref([])

const products = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

const endAt = ref(Date.now() + 3600e3)
const cd = reactive({ h: '00', m: '00', s: '00' })
let timer

function pad(n) {
  return String(n).padStart(2, '0')
}
function tick() {
  let left = Math.max(0, endAt.value - Date.now())
  const h = Math.floor(left / 3600e3)
  left -= h * 3600e3
  const m = Math.floor(left / 60e3)
  left -= m * 60e3
  const s = Math.floor(left / 1e3)
  cd.h = pad(h)
  cd.m = pad(m)
  cd.s = pad(s)
}

function parseEndAt(v) {
  if (!v) return Date.now() + 3600e3
  if (typeof v === 'number') return v
  const t = new Date(String(v).replace(/-/g, '/')).getTime()
  return Number.isFinite(t) ? t : Date.now() + 3600e3
}

async function loadSeckill() {
  try {
    const res = await getSeckillCurrent()
    const data = res || {}
    endAt.value = parseEndAt(data.end_at)
    seckillItems.value = data.items || []
    tick()
  } catch {
    seckillItems.value = []
  }
}

function toast(title) {
  uni.showToast({ title, icon: 'none' })
}

function onCategory(c) {
  if (c.name === '全部商户' || c.brand) {
    goShopList(c.brand ? 'brand_shop' : 'quality_shop')
    return
  }
  if (c.name === '种草社区') {
    goCommunityList()
    return
  }
  if (c.name === '限时秒杀') {
    goSeckillList()
    return
  }
  if (c.pointsMall || c.name === '积分商城') {
    goPointsMall()
    return
  }
  if (c.cps || c.name === '优惠购') {
    goCpsList()
    return
  }
  if (c.local || c.name === '本地商家') {
    goLocalShops()
    return
  }
  if (c.coupon) {
    goCouponCenter()
    return
  }
  if (c.orders) {
    uni.navigateTo({ url: '/pages/order/list' })
    return
  }
  if (c.messages) {
    goMessages()
    return
  }
  if (c.more) {
    goCategoryPage()
    return
  }
  if (c.id) {
    goCategoryPage(c.id)
    return
  }
  toast(`${c.name}即将开放`)
}

function goCategoryPage(categoryId) {
  if (categoryId) {
    uni.setStorageSync('category_focus_id', categoryId)
  } else {
    uni.removeStorageSync('category_focus_id')
  }
  uni.switchTab({ url: '/pages/category/index' })
}

function onBannerTap(b) {
  if (!b || typeof b === 'string') return
  if (b.link_type === 'product' && b.link_id) {
    goDetail(b.link_id)
    return
  }
  if (b.link_type === 'article' && b.link_id) {
    goArticle(b.link_id)
  }
}

function goShopList(slotType) {
  uni.navigateTo({ url: `/pages/shop/list?slot_type=${slotType || 'quality_shop'}` })
}

function goShopDetail(id) {
  if (!id) return
  uni.navigateTo({ url: `/pages/shop/detail?id=${id}` })
}

function goCommunityList() {
  uni.navigateTo({ url: '/pages/community/list' })
}

function goArticle(id) {
  if (!id) return
  uni.navigateTo({ url: `/pages/community/detail?id=${id}` })
}

function goDetail(id) {
  uni.navigateTo({ url: `/pages/product/detail?id=${id}` })
}

function goSeckillDetail(p) {
  uni.navigateTo({
    url: `/pages/product/detail?id=${p.product_id}&seckill_entry_id=${p.id}`,
  })
}

function goSeckillList() {
  uni.navigateTo({ url: '/pages/seckill/list' })
}

function goSalesRank() {
  uni.navigateTo({ url: '/pages/product/rank' })
}

function onThemeTap(t) {
  if (!t) return
  if (t.link_type === 'shop' && t.link_id) {
    goShopDetail(t.link_id)
    return
  }
  if (t.link_type === 'product' && t.link_id) {
    goDetail(t.link_id)
    return
  }
  if (t.link_type === 'category' && t.link_id) {
    goCategoryPage(t.link_id)
  }
}

async function loadThemes() {
  try {
    const res = await listThemeTiles()
    themes.value = res?.list || []
  } catch {
    themes.value = []
  }
}

async function loadCouponBar() {
  try {
    const res = await listCouponCenter()
    const first = (res?.list || [])[0]
    if (first) {
      const face = first.coupon_type === 'discount'
        ? `${(Number(first.discount_rate) * 10).toFixed(1)}折券`
        : `满${first.threshold_amount || 0}减${first.discount_amount}`
      couponBarText.value = `${first.name} · ${face}，立即领取`
    }
  } catch { /* keep default */ }
}

function goCouponCenter() {
  uni.navigateTo({ url: '/pages/coupon/center' })
}

function goPointsMall() {
  uni.navigateTo({ url: '/pages/points-mall/list' })
}

function goCpsList() {
  uni.navigateTo({ url: '/pages/cps/index' })
}

function goMessages() {
  if (!isLoggedIn()) {
    uni.navigateTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/message/list') })
    return
  }
  uni.navigateTo({ url: '/pages/message/list' })
}

function goTaskCenter() {
  if (!isLoggedIn()) {
    uni.navigateTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/task/index') })
    return
  }
  uni.navigateTo({ url: '/pages/task/index' })
}

async function loadUnread() {
  if (!isLoggedIn()) {
    unreadCount.value = 0
    return
  }
  try {
    const res = await getNotificationUnreadCount()
    unreadCount.value = Number(res?.count) || 0
  } catch {
    unreadCount.value = 0
  }
}

function rankDisplayName(r) {
  if (r.shop_name) return `${r.name} | ${r.shop_name}`
  return r.name
}

function rankSoldText(r) {
  if (r.today_sold > 0) return `今日售出${r.today_sold}`
  return `总销量${r.sold_count || 0}`
}

async function loadSalesRank() {
  try {
    const res = await listSalesRank({ page: 1, page_size: 3 })
    rankList.value = res?.list || []
  } catch {
    rankList.value = []
  }
}

async function loadCats() {
  try {
    const res = await listCategories({ page: 1, page_size: 20 })
    apiCats.value = res?.list || []
  } catch {
    apiCats.value = []
  }
}

function mapShop(s) {
  return {
    id: s.id,
    name: s.name,
    tag: s.category || '优选商户',
    category: s.category,
    desc: s.description || '',
    paid: !!s.paid,
    city: s.city || '',
    img: s.storefront_image || s.logo || placeholder,
  }
}

function mapNote(a, i) {
  return {
    id: a.id,
    w: i % 3 === 0 ? 64 : 44,
    label: a.paid ? '推广精选' : '种草笔记',
    likes: a.like_count || 0,
    reads: a.read_count || 0,
    title: a.title,
    imgs: [a.cover_url || placeholder],
  }
}

async function loadHomeSlots() {
  try {
    const [brandRes, qualityRes] = await Promise.all([
      listHomeSlots('brand_shop'),
      listHomeSlots('quality_shop'),
    ])
    brandShops.value = (Array.isArray(brandRes) ? brandRes : []).map(mapShop)
    shops.value = (Array.isArray(qualityRes) ? qualityRes : []).map(mapShop)
  } catch {
    brandShops.value = []
    shops.value = []
  }
}

function formatDist(km) {
  if (km == null || Number.isNaN(Number(km))) return '—'
  const n = Number(km)
  if (n < 1) return `${Math.round(n * 1000)}m`
  return `${n.toFixed(1)}km`
}

async function loadLocalShops() {
  if (!hasCoords()) {
    localShops.value = []
    localHint.value = '请先选择定位'
    return
  }
  try {
    const coords = getCoords()
    const res = await listLocalShops({
      page: 1,
      page_size: 10,
      lat: coords.latitude,
      lng: coords.longitude,
      sort: 'distance',
    })
    const rows = res?.list || []
    localShops.value = rows.map((s) => ({
      id: s.id,
      name: s.name,
      city: s.city,
      address: s.full_address || s.address,
      desc: s.description || '',
      distance_km: s.distance_km,
      img: s.storefront_image || s.logo || placeholder,
    }))
    localHint.value = rows.length ? '' : '附近暂无本地商家'
  } catch {
    localShops.value = []
    localHint.value = '本地商家加载失败'
  }
}

function goLocalShops() {
  uni.navigateTo({ url: '/pages/shop/local' })
}

async function syncCityFromAddress() {
  if (!isLoggedIn()) return
  try {
    const res = await listAddresses({ silent: true })
    const list = res?.list || res || []
    const rows = Array.isArray(list) ? list : []
    const addr = rows.find((a) => a.is_default) || rows[0]
    if (applyAddressCityIfNeeded(addr)) {
      locCity.value = getCity()
    }
  } catch { /* ignore */ }
}

function goChooseCity() {
  uni.navigateTo({ url: '/pages/location/city' })
}

async function loadNotes() {
  try {
    const res = await listArticles({ page: 1, page_size: 10, home: 1 })
    notes.value = (res?.list || []).map(mapNote)
  } catch {
    notes.value = []
  }
}

async function loadBanners() {
  try {
    const res = await listBanners()
    const rows = res?.list || []
    banners.value = rows.length ? rows : [...fallbackBanners]
  } catch {
    banners.value = [...fallbackBanners]
  }
}

async function loadProducts(reset = false) {
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      products.value = []
    }
    const res = await listProducts({ page: page.value, page_size: pageSize, order_by: 'sold_count_desc' })
    const list = res?.list || []
    products.value = reset ? list : products.value.concat(list)
    const total = res?.total || 0
    if (products.value.length >= total || list.length < pageSize) {
      finished.value = true
    } else {
      page.value += 1
    }
  } catch {
    if (reset) products.value = []
  } finally {
    loading.value = false
  }
}

function loadMore() {
  loadProducts(false)
}

onReachBottom(() => loadMore())

onShow(async () => {
  loadUnread()
  await syncCityFromAddress()
  locCity.value = getCity()
  loadHomeSlots()
  loadLocalShops()
})

onMounted(() => {
  tick()
  timer = setInterval(tick, 1000)
  loadBanners()
  loadCats()
  loadHomeSlots()
  loadLocalShops()
  loadNotes()
  loadSeckill()
  loadThemes()
  loadCouponBar()
  loadSalesRank()
  loadProducts(true)
  loadUnread()
  unreadTimer = setInterval(loadUnread, 30000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
  if (unreadTimer) clearInterval(unreadTimer)
})
</script>

<style scoped>
.page { padding-bottom: 24rpx; }
.coupon-bar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16rpx 32rpx; color: #fff; font-size: 26rpx;
}
.coupon-btn {
  background: rgba(255,255,255,.2); padding: 4rpx 20rpx; border-radius: 999rpx; font-size: 22rpx;
}
.header {
  position: sticky; top: 0; z-index: 50;
  display: flex; align-items: center; gap: 16rpx;
  padding: 20rpx 32rpx; background: #fff;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.loc { font-size: 22rpx; color: var(--shop-sub); display: flex; align-items: center; gap: 6rpx; white-space: nowrap; }
.search {
  flex: 1; display: flex; align-items: center; gap: 12rpx;
  background: rgba(230,213,188,.2); border-radius: 999rpx; padding: 14rpx 24rpx;
}
.search-ph { font-size: 22rpx; color: rgba(113,113,122,.6); }
.bell-wrap { position: relative; width: 48rpx; height: 48rpx; display: flex; align-items: center; justify-content: center; }
.bell { font-size: 36rpx; display: inline-block; transform-origin: top center; }
.bell.shake {
  animation: bell-shake 1.2s ease-in-out infinite;
}
@keyframes bell-shake {
  0%, 100% { transform: rotate(0deg) scale(1); }
  10% { transform: rotate(14deg) scale(1.05); }
  20% { transform: rotate(-12deg) scale(1.05); }
  30% { transform: rotate(10deg); }
  40% { transform: rotate(-8deg); }
  50% { transform: rotate(6deg); }
  60% { transform: rotate(-4deg); }
  70% { transform: rotate(2deg); }
  80% { transform: rotate(0deg) scale(1); }
}
.badge {
  position: absolute; top: -6rpx; right: -10rpx; min-width: 28rpx; height: 28rpx; padding: 0 6rpx;
  border-radius: 14rpx; background: #ef4444; color: #fff; font-size: 18rpx; line-height: 28rpx; text-align: center;
}
.banner-wrap { margin: 24rpx 32rpx 0; border-radius: 32rpx; overflow: hidden; position: relative; box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08); }
.banner { height: 344rpx; }
.banner-img { width: 100%; height: 344rpx; }
.points-entry {
  position: relative;
  margin: 24rpx 32rpx 0;
  height: 268rpx;
  border-radius: 28rpx;
  overflow: hidden;
  background:
    radial-gradient(ellipse 70% 80% at 55% 40%, rgba(255, 214, 140, 0.55) 0%, transparent 55%),
    radial-gradient(circle at 20% 80%, rgba(255, 140, 70, 0.35) 0%, transparent 40%),
    linear-gradient(155deg, #6b4423 0%, #9a5c28 38%, #c4783a 68%, #e0a04a 100%);
  box-shadow: 0 12rpx 36rpx rgba(154, 92, 40, 0.28);
}
.pe-glow {
  position: absolute; inset: -20%;
  background: radial-gradient(circle at 50% 45%, rgba(255, 236, 180, 0.45) 0%, transparent 48%);
  animation: pe-pulse 2.8s ease-in-out infinite;
  pointer-events: none;
}
.pe-ring {
  position: absolute; left: 50%; top: 42%;
  width: 220rpx; height: 220rpx;
  margin: -110rpx 0 0 -110rpx;
  border-radius: 50%;
  border: 6rpx solid rgba(255, 214, 120, 0.55);
  box-shadow:
    0 0 0 10rpx rgba(255, 180, 80, 0.12),
    inset 0 0 40rpx rgba(255, 220, 140, 0.25),
    0 0 48rpx rgba(255, 190, 90, 0.35);
  pointer-events: none;
}
.pe-beam {
  position: absolute; left: 50%; top: 40%; width: 8rpx; height: 160%;
  margin-left: -4rpx; transform-origin: center top;
  background: linear-gradient(180deg, rgba(255, 230, 160, 0.55), transparent 70%);
  pointer-events: none;
}
.pe-beam-a { transform: rotate(-28deg); opacity: .7; }
.pe-beam-b { transform: rotate(32deg); opacity: .55; }
.pe-floor {
  position: absolute; left: -10%; right: -10%; bottom: -8%;
  height: 46%;
  background: linear-gradient(180deg, transparent 0%, rgba(60, 32, 12, 0.35) 100%);
  transform: perspective(200rpx) rotateX(18deg);
  pointer-events: none;
}
.pe-main {
  position: absolute; left: 0; right: 0; top: 58rpx;
  display: flex; flex-direction: column; align-items: center; z-index: 3;
}
.pe-title {
  font-size: 56rpx; font-weight: 800; letter-spacing: 6rpx; color: #fffaf0;
  text-shadow:
    0 2rpx 0 #8b5a20,
    0 4rpx 0 #6b4010,
    0 8rpx 0 #4a2c08,
    0 12rpx 20rpx rgba(40, 20, 0, 0.45);
}
.pe-sub {
  margin-top: 14rpx;
  padding: 8rpx 28rpx;
  border-radius: 999rpx;
  font-size: 22rpx; color: #fff8e8; font-weight: 600;
  background: linear-gradient(135deg, rgba(90, 48, 18, 0.88), rgba(140, 72, 28, 0.85));
  border: 1rpx solid rgba(255, 220, 160, 0.55);
  box-shadow: 0 4rpx 12rpx rgba(40, 20, 0, 0.25);
}
.pe-cta {
  position: absolute; right: 28rpx; bottom: 28rpx; z-index: 4;
  padding: 14rpx 34rpx; border-radius: 999rpx;
  font-size: 26rpx; font-weight: 700; color: #5c3a12;
  background: linear-gradient(135deg, #ffe9a8 0%, #ffc85a 45%, #f0a020 100%);
  border: 2rpx solid rgba(255, 245, 200, 0.9);
  box-shadow:
    0 0 0 4rpx rgba(255, 180, 60, 0.25),
    0 8rpx 24rpx rgba(255, 160, 40, 0.55),
    inset 0 2rpx 0 rgba(255, 255, 255, 0.55);
  animation: pe-cta-glow 1.8s ease-in-out infinite;
}
.pe-chip {
  position: absolute; z-index: 3;
  padding: 6rpx 16rpx; border-radius: 12rpx;
  font-size: 18rpx; font-weight: 600; color: #7a4a18;
  background: linear-gradient(180deg, #fff8e8, #ffe2a8);
  border: 1rpx solid rgba(255, 200, 120, 0.8);
  box-shadow: 0 4rpx 12rpx rgba(80, 40, 10, 0.2);
}
.pe-chip-a { left: 24rpx; top: 36rpx; transform: rotate(-8deg); }
.pe-chip-b { left: 36rpx; top: 118rpx; transform: rotate(4deg); }
.pe-chip-c { right: 36rpx; top: 48rpx; transform: rotate(7deg); }
.pe-gift {
  position: absolute; left: 28rpx; bottom: 32rpx; z-index: 2;
  width: 56rpx; height: 52rpx;
  animation: pe-float-a 2.2s ease-in-out infinite;
  filter: drop-shadow(0 6rpx 8rpx rgba(40, 20, 0, 0.3));
  pointer-events: none;
}
.pe-gift-body {
  position: absolute; left: 4rpx; bottom: 0; width: 48rpx; height: 34rpx;
  border-radius: 6rpx;
  background: linear-gradient(180deg, #ff9a4a, #e86a20);
}
.pe-gift-body::before {
  content: ''; position: absolute; left: 50%; top: 0; bottom: 0; width: 10rpx;
  margin-left: -5rpx; background: #ffe08a;
}
.pe-gift-lid {
  position: absolute; left: 0; top: 8rpx; width: 56rpx; height: 14rpx;
  border-radius: 4rpx;
  background: linear-gradient(180deg, #ffb06a, #f07828);
}
.pe-gift-bow {
  position: absolute; left: 50%; top: 0; width: 18rpx; height: 14rpx;
  margin-left: -9rpx; border-radius: 50% 50% 40% 40%;
  background: #ffe08a;
  box-shadow: -14rpx 2rpx 0 -2rpx #ffe08a, 14rpx 2rpx 0 -2rpx #ffe08a;
}
.pe-ticket {
  position: absolute; right: 118rpx; top: 30rpx; z-index: 2;
  width: 52rpx; height: 30rpx; border-radius: 6rpx;
  background: linear-gradient(135deg, #fff3c4, #ffd56a);
  border: 2rpx dashed rgba(180, 110, 20, 0.45);
  box-shadow: 0 4rpx 10rpx rgba(80, 40, 10, 0.25);
  animation: pe-float-b 2.6s ease-in-out infinite;
  pointer-events: none;
}
.pe-ticket-cut {
  position: absolute; left: -6rpx; top: 50%; width: 12rpx; height: 12rpx;
  margin-top: -6rpx; border-radius: 50%;
  background: transparent;
  box-shadow: 46rpx 0 0 0 transparent;
}
.pe-ticket::after {
  content: ''; position: absolute; right: 10rpx; top: 6rpx; bottom: 6rpx; width: 10rpx;
  border-radius: 2rpx; background: rgba(200, 120, 30, 0.25);
}
.pe-bolt {
  position: absolute; left: 42%; bottom: 40rpx; z-index: 2;
  width: 0; height: 0;
  border-left: 10rpx solid transparent;
  border-right: 10rpx solid transparent;
  border-bottom: 28rpx solid #ffe08a;
  transform: skewX(-12deg);
  filter: drop-shadow(0 0 8rpx rgba(255, 200, 80, 0.7));
  animation: pe-twinkle 1.6s ease-in-out infinite .2s;
  pointer-events: none;
}
.pe-coin {
  position: absolute; z-index: 2;
  width: 44rpx; height: 44rpx; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 22rpx; font-weight: 800; color: #8a5200;
  background: linear-gradient(145deg, #fff1b8, #f0b429 55%, #d4920a);
  border: 2rpx solid rgba(255, 240, 180, 0.9);
  box-shadow: 0 4rpx 10rpx rgba(120, 70, 10, 0.35);
  pointer-events: none;
}
.pe-coin.c1 { right: 48rpx; top: 100rpx; animation: pe-float-a 1.8s ease-in-out infinite; }
.pe-coin.c2 { left: 110rpx; bottom: 28rpx; width: 36rpx; height: 36rpx; font-size: 18rpx; animation: pe-float-b 2.1s ease-in-out infinite; }
.pe-coin.c3 { right: 150rpx; bottom: 52rpx; width: 40rpx; height: 40rpx; font-size: 18rpx; animation: pe-float-a 2.4s ease-in-out infinite .3s; }
.pe-spark {
  position: absolute; z-index: 2; width: 12rpx; height: 12rpx;
  background: #ffe08a; transform: rotate(45deg);
  box-shadow: 0 0 10rpx rgba(255, 220, 120, 0.9);
  animation: pe-twinkle 1.4s ease-in-out infinite;
  pointer-events: none;
}
.pe-spark.s1 { left: 48%; top: 28rpx; }
.pe-spark.s2 { left: 62%; top: 72rpx; width: 8rpx; height: 8rpx; animation-delay: .35s; background: #fff6d0; }
.pe-spark.s3 { left: 28%; bottom: 72rpx; width: 10rpx; height: 10rpx; animation-delay: .7s; }
.pe-spark.s4 { right: 28%; top: 40rpx; animation-delay: 1s; background: #ffb347; }
@keyframes pe-pulse {
  0%, 100% { opacity: .7; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.06); }
}
@keyframes pe-cta-glow {
  0%, 100% { box-shadow: 0 0 0 4rpx rgba(255, 180, 60, 0.25), 0 8rpx 24rpx rgba(255, 160, 40, 0.55), inset 0 2rpx 0 rgba(255, 255, 255, 0.55); }
  50% { box-shadow: 0 0 0 8rpx rgba(255, 180, 60, 0.35), 0 10rpx 32rpx rgba(255, 170, 50, 0.7), inset 0 2rpx 0 rgba(255, 255, 255, 0.65); }
}
@keyframes pe-float-a {
  0%, 100% { transform: translateY(0) rotate(-6deg); }
  50% { transform: translateY(-12rpx) rotate(4deg); }
}
@keyframes pe-float-b {
  0%, 100% { transform: translateY(0) rotate(8deg); }
  50% { transform: translateY(-10rpx) rotate(-4deg); }
}
@keyframes pe-twinkle {
  0%, 100% { opacity: .35; transform: rotate(45deg) scale(.8); }
  50% { opacity: 1; transform: rotate(45deg) scale(1.25); }
}
.cps-entry {
  margin: 24rpx 32rpx 0;
  padding: 32rpx;
  border-radius: 28rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #f3ebe0 0%, #e8d5b5 55%, #d4b890 100%);
  box-shadow: 0 8rpx 28rpx rgba(200,168,118,.18);
}
.cps-entry-text { display: flex; flex-direction: column; gap: 8rpx; flex: 1; min-width: 0; }
.cps-entry-title { font-size: 34rpx; font-weight: 700; color: #5c4a32; }
.cps-entry-sub { font-size: 22rpx; color: rgba(92,74,50,.82); }
.cps-entry-cta {
  flex-shrink: 0; font-size: 26rpx; color: #5c4a32; background: rgba(255,255,255,.78);
  padding: 14rpx 28rpx; border-radius: 999rpx; font-weight: 600;
}
.corner-tag { position: absolute; top: 0; right: 0; width: 140rpx; height: 140rpx; overflow: hidden; }
.corner-text {
  position: absolute; top: 24rpx; right: -44rpx; width: 180rpx; text-align: center;
  font-size: 20rpx; color: #fff; padding: 4rpx 0; transform: rotate(45deg);
}
.scroll-hide { width: 100%; white-space: nowrap; }
.cat-grid {
  display: flex; flex-wrap: wrap; padding: 24rpx 4rpx 12rpx;
}
.cat-item {
  width: 20%; box-sizing: border-box;
  display: flex; flex-direction: column; align-items: center;
  padding: 12rpx 0 18rpx;
}
.cat-icon {
  width: 100rpx; height: 100rpx; border-radius: 28rpx;
  display: flex; align-items: center; justify-content: center;
}
.cat-icon.bounce {
  animation: cat-float 1.2s ease-in-out infinite;
}
.cat-img { width: 80rpx; height: 80rpx; }
.cat-item:active .cat-icon { transform: scale(0.92); }
.cat-item:active .cat-icon.bounce { animation-play-state: paused; }
@keyframes cat-float {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-10rpx) scale(1.06); }
}
.cat-name {
  font-size: 20rpx; margin-top: 10rpx; color: var(--shop-text);
  width: 100%; text-align: center; padding: 0 4rpx;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.task-fab {
  position: fixed; right: 5rpx; bottom: calc(130rpx + env(safe-area-inset-bottom));
  width: 88rpx; height: 88rpx; z-index: 50;
  display: flex; align-items: center; justify-content: center;
}
.task-fab-glow {
  position: absolute; inset: -8rpx; border-radius: 50%;
  background: radial-gradient(circle, rgba(255, 180, 90, 0.55) 0%, rgba(255, 140, 80, 0.2) 50%, transparent 72%);
  animation: fab-glow 2s ease-in-out infinite;
  pointer-events: none;
}
.task-fab-art {
  position: relative; z-index: 1;
  width: 80rpx; height: 80rpx;
  border-radius: 50%;
  animation: fab-bob 2s ease-in-out infinite;
  filter: drop-shadow(0 10rpx 18rpx rgba(232, 140, 70, 0.4));
}
.task-fab-coin {
  position: absolute; z-index: 2;
  width: 28rpx; height: 28rpx; border-radius: 50%;
  background: linear-gradient(145deg, #ffe08a, #f0b429);
  color: #c47a12; font-size: 16rpx; font-weight: 800;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 4rpx 8rpx rgba(196, 122, 18, 0.35);
  pointer-events: none;
}
.task-fab-coin.c1 {
  top: 6rpx; right: 16rpx;
  animation: coin-float-a 1.6s ease-in-out infinite;
}
.task-fab-coin.c2 {
  top: 26rpx; left: 8rpx;
  animation: coin-float-b 1.9s ease-in-out infinite;
}
.task-fab-spark {
  position: absolute; z-index: 2; pointer-events: none;
  width: 10rpx; height: 10rpx; border-radius: 2rpx;
  background: #ffd76a;
  animation: spark-twinkle 1.4s ease-in-out infinite;
}
.task-fab-spark.s1 { top: 12rpx; left: 36rpx; animation-delay: 0s; }
.task-fab-spark.s2 { top: 40rpx; right: 8rpx; width: 8rpx; height: 8rpx; background: #fff6d6; animation-delay: .3s; }
.task-fab-spark.s3 { bottom: 36rpx; left: 14rpx; width: 7rpx; height: 7rpx; transform: rotate(25deg); animation-delay: .6s; }
.task-fab-spark.s4 { bottom: 48rpx; right: 22rpx; background: #ffb347; animation-delay: .9s; }
.task-fab-badge {
  position: absolute; z-index: 3; left: 50%; bottom: -10rpx;
  transform: translateX(-50%);
  padding: 4rpx 16rpx; border-radius: 20rpx;
  font-size: 18rpx; font-weight: 700; color: #fff;
  background: linear-gradient(90deg, #f0a04b, #e86b3a);
  box-shadow: 0 4rpx 10rpx rgba(232, 107, 58, 0.35);
  white-space: nowrap;
}
.task-fab:active .task-fab-art { transform: scale(0.92); animation-play-state: paused; }
@keyframes fab-bob {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-10rpx) scale(1.04); }
}
@keyframes fab-glow {
  0%, 100% { opacity: 0.65; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.1); }
}
@keyframes coin-float-a {
  0%, 100% { transform: translateY(0) rotate(-8deg); opacity: 1; }
  50% { transform: translateY(-14rpx) rotate(12deg); opacity: 1; }
}
@keyframes coin-float-b {
  0%, 100% { transform: translateY(0) rotate(10deg) scale(0.9); }
  50% { transform: translateY(-12rpx) rotate(-6deg) scale(1); }
}
@keyframes spark-twinkle {
  0%, 100% { opacity: 0.25; transform: scale(0.7) rotate(0deg); }
  50% { opacity: 1; transform: scale(1.2) rotate(20deg); }
}
.seckill { padding: 28rpx; }
.sec-head, .block-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20rpx; }
.block-head.no-pad { margin-bottom: 24rpx; }
.sec-left { display: flex; align-items: center; gap: 16rpx; flex-wrap: wrap; }
.countdown { display: flex; align-items: center; gap: 6rpx; font-size: 24rpx; color: var(--shop-sub); }
.count-num { background: #1d1d1f; color: #fff; padding: 2rpx 10rpx; border-radius: 8rpx; font-weight: 700; }
.seckill-row { display: inline-flex; gap: 20rpx; padding-bottom: 8rpx; }
.seckill-item { width: 220rpx; display: inline-block; }
.seckill-empty { padding: 24rpx 32rpx; display: inline-block; }
.seckill-img { width: 220rpx; height: 200rpx; border-radius: 16rpx; }
.seckill-name { font-size: 22rpx; margin-top: 12rpx; height: 64rpx; }
.price-row { display: flex; align-items: baseline; margin-top: 8rpx; }
.stock-left { font-size: 20rpx; color: var(--shop-sub); margin-top: 4rpx; }
.block { margin: 32rpx 32rpx 0; }
.brand-row, .shop-row, .note-row { display: inline-flex; gap: 24rpx; padding-bottom: 8rpx; }
.brand-card, .shop-card {
  width: 300rpx; background: #fff; border-radius: 32rpx; overflow: hidden;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08); display: inline-block;
}
.brand-cover, .shop-cover { width: 100%; height: 180rpx; }
.brand-body, .shop-body { padding: 20rpx; }
.brand-name { font-size: 26rpx; font-weight: 500; }
.brand-meta { display: flex; justify-content: space-between; margin-top: 12rpx; font-size: 22rpx; }
.gold { color: var(--shop-gold); }
.sub { color: var(--shop-sub); font-size: 22rpx; }
.tag-scroll { margin-bottom: 16rpx; }
.tag-row { display: inline-flex; gap: 12rpx; }
.note-tag {
  font-size: 20rpx; padding: 8rpx 20rpx; border-radius: 999rpx;
  background: rgba(230,213,188,.3); color: var(--shop-gold);
}
.note-tag.active { background: linear-gradient(135deg,#bfa472,#d4b890); color: #fff; }
.note-card {
  background: #fff; border-radius: 32rpx; overflow: hidden;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08); display: inline-block;
}
.w-64 { width: 480rpx; }
.w-44 { width: 320rpx; }
.note-cover { width: 100%; height: 260rpx; }
.note-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8rpx; padding: 12rpx 12rpx 0; height: 180rpx; }
.note-grid-img { width: 100%; height: 160rpx; border-radius: 12rpx; }
.note-body { padding: 16rpx 20rpx 20rpx; }
.note-label { font-size: 20rpx; color: var(--shop-sub); }
.note-title { font-size: 24rpx; margin-top: 8rpx; font-weight: 500; }
.note-foot { display: flex; justify-content: space-between; margin-top: 12rpx; font-size: 20rpx; color: var(--shop-sub); }
.theme { padding: 28rpx; }
.theme-title { display: block; margin-bottom: 24rpx; }
.theme-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20rpx; }
.theme-item { position: relative; height: 200rpx; border-radius: 24rpx; overflow: hidden; }
.theme-paid {
  position: absolute; top: 12rpx; right: 12rpx; font-size: 18rpx;
  color: #c8a876; background: rgba(0,0,0,.45); padding: 2rpx 10rpx; border-radius: 6rpx;
}
.theme-img { width: 100%; height: 100%; }
.theme-mask {
  position: absolute; left: 20rpx; bottom: 20rpx; color: #fff;
  display: flex; flex-direction: column; gap: 4rpx;
}
.theme-name { font-weight: 700; font-size: 28rpx; }
.theme-desc { font-size: 20rpx; opacity: .85; }
.shop-score { display: block; margin-top: 8rpx; }
.shop-promo { display: flex; justify-content: space-between; margin-top: 12rpx; font-size: 22rpx; }
.rank, .goods { padding: 28rpx; margin-bottom: 24rpx; }
.rank-item { display: flex; align-items: center; gap: 20rpx; margin-bottom: 20rpx; }
.rank-empty { padding: 12rpx 0 8rpx; }
.rank-no.r2 { background: linear-gradient(135deg,#f59e0b,#fbbf24); }
.rank-no.r3 { background: linear-gradient(135deg,#94a3b8,#cbd5e1); }
.rank-no {
  width: 40rpx; height: 40rpx; border-radius: 50%; color: #fff; font-size: 22rpx; font-weight: 700;
  display: flex; align-items: center; justify-content: center; background: var(--shop-gold);
}
.rank-no.r1 { background: linear-gradient(135deg,#d83636,#f25757); }
.rank-img { width: 112rpx; height: 112rpx; border-radius: 16rpx; }
.rank-info { flex: 1; min-width: 0; }
.rank-name { font-size: 26rpx; }
.rank-foot { display: flex; justify-content: space-between; margin-top: 8rpx; }
.goods-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 24rpx; }
.goods-img { width: 100%; height: 280rpx; border-radius: 24rpx; background: #f3f3f3; }
.goods-name { font-size: 26rpx; margin-top: 12rpx; font-weight: 500; min-height: 72rpx; }
.goods-price { font-size: 32rpx; }
.sold { display: block; margin-top: 6rpx; }
.load-more { text-align: center; color: var(--shop-sub); font-size: 24rpx; padding: 24rpx 0 8rpx; }
</style>
