<template>
  <view class="page">
    <!-- 优惠券通栏 -->
    <view class="coupon-bar gold-gradient">
      <text class="coupon-text">平台通用券满100减30，立即领取</text>
      <text class="coupon-btn" @tap="toast('领取功能即将开放')">去领取</text>
    </view>

    <!-- 顶栏 -->
    <view class="header">
      <view class="loc" @tap="toast('定位即将开放')">
        <text class="loc-icon">📍</text>
        <text>长沙市</text>
      </view>
      <view class="search" @tap="toast('搜索即将开放')">
        <text class="search-icon">🔍</text>
        <text class="search-ph">搜索店铺 / 好物 / 种草笔记</text>
      </view>
      <view class="bell" @tap="toast('暂无消息')">🔔</view>
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

    <!-- 分类入口 -->
    <scroll-view scroll-x :show-scrollbar="false" class="cat-scroll scroll-hide section-card">
      <view class="cat-row">
        <view v-for="c in categoryEntries" :key="c.name" class="cat-item" @tap="onCategory(c)">
          <view class="cat-icon" :style="{ background: c.bg, color: c.color }">{{ c.emoji }}</view>
          <text class="cat-name">{{ c.name }}</text>
        </view>
      </view>
    </scroll-view>

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
              <text class="sub shop-score">{{ s.category || '优选商户' }}{{ s.paid ? ' · 推广' : '' }}</text>
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
  </view>
</template>

<script setup>
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { onReachBottom } from '@dcloudio/uni-app'
import { getSeckillCurrent, listArticles, listBanners, listCategories, listHomeSlots, listProducts, listSalesRank, listThemeTiles } from '../../api/index'

const placeholder = 'https://picsum.photos/id/96/400/400'
const fallbackBanners = [
  { id: 'f1', image_url: 'https://picsum.photos/id/1059/750/340', link_type: 'none' },
  { id: 'f2', image_url: 'https://picsum.photos/id/1062/750/340', link_type: 'none' },
  { id: 'f3', image_url: 'https://picsum.photos/id/1068/750/340', link_type: 'none' },
]
const banners = ref([...fallbackBanners])

const fixedCats = [
  { name: '全部商户', emoji: '🏬', bg: 'rgba(230,213,188,.35)', color: '#C8A876' },
  { name: '种草社区', emoji: '✍️', bg: 'rgba(230,213,188,.35)', color: '#C8A876' },
  { name: '限时秒杀', emoji: '⚡', bg: '#FEF2F2', color: '#D83636' },
]
const apiCats = ref([])
const categoryEntries = computed(() => {
  const fromApi = apiCats.value.slice(0, 5).map((c, i) => {
    const colors = [
      { bg: '#DCFCE7', color: '#16A34A' },
      { bg: '#DBEAFE', color: '#2563EB' },
      { bg: '#FCE7F3', color: '#DB2777' },
      { bg: '#FEF3C7', color: '#D97706' },
      { bg: 'rgba(230,213,188,.35)', color: '#C8A876' },
    ]
    const sty = colors[i % colors.length]
    return { name: c.name, emoji: '📦', id: c.id, ...sty }
  })
  return [fixedCats[0], ...fromApi, fixedCats[1], fixedCats[2]]
})

const seckillItems = ref([])
const brandShops = ref([])
const shops = ref([])
const notes = ref([])

const themes = ref([])

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
  if (c.name === '全部商户') {
    goShopList('quality_shop')
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

onMounted(() => {
  tick()
  timer = setInterval(tick, 1000)
  loadBanners()
  loadCats()
  loadHomeSlots()
  loadNotes()
  loadSeckill()
  loadThemes()
  loadSalesRank()
  loadProducts(true)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
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
.bell { font-size: 36rpx; }
.banner-wrap { margin: 24rpx 32rpx 0; border-radius: 32rpx; overflow: hidden; position: relative; box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08); }
.banner { height: 344rpx; }
.banner-img { width: 100%; height: 344rpx; }
.corner-tag { position: absolute; top: 0; right: 0; width: 140rpx; height: 140rpx; overflow: hidden; }
.corner-text {
  position: absolute; top: 24rpx; right: -44rpx; width: 180rpx; text-align: center;
  font-size: 20rpx; color: #fff; padding: 4rpx 0; transform: rotate(45deg);
}
.scroll-hide { width: 100%; white-space: nowrap; }
.cat-scroll { padding: 28rpx 0; white-space: nowrap; }
.cat-row { display: inline-flex; gap: 40rpx; padding: 0 24rpx; }
.cat-item { width: 112rpx; display: inline-flex; flex-direction: column; align-items: center; }
.cat-icon {
  width: 88rpx; height: 88rpx; border-radius: 50%;
  display: flex; align-items: center; justify-content: center; font-size: 36rpx;
}
.cat-name { font-size: 22rpx; margin-top: 12rpx; color: var(--shop-text); }
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
