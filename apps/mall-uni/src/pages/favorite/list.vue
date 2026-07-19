<template>
  <view class="page">
    <view class="tabs">
      <text class="tab" :class="{ active: tab === 'product' }" @tap="switchTab('product')">商品</text>
      <text class="tab" :class="{ active: tab === 'article' }" @tap="switchTab('article')">文章</text>
    </view>

    <template v-if="tab === 'product'">
      <view v-if="list.length" class="toolbar">
        <text @tap="toggleEdit">{{ editing ? '完成' : '管理' }}</text>
        <text v-if="editing" class="danger" @tap="batchCancel">取消收藏</text>
      </view>
      <view v-if="!list.length && !loading" class="empty">暂无收藏商品</view>
      <view
        v-for="item in list"
        :key="item.id"
        class="card"
        :class="{ invalid: item.invalid }"
        @tap="onTapProduct(item)"
      >
        <view v-if="editing" class="check" @tap.stop="toggle(item.product_id)">
          <text>{{ selected.has(item.product_id) ? '✓' : '' }}</text>
        </view>
        <image class="cover" :src="item.main_image || placeholder" mode="aspectFill" />
        <view class="info">
          <text class="name">{{ item.name || '商品已失效' }}</text>
          <text v-if="item.invalid" class="tag">已失效</text>
          <text v-else class="price">¥{{ item.sale_price }}</text>
          <text v-if="!editing" class="cancel" @tap.stop="cancelOne(item.product_id)">取消收藏</text>
        </view>
      </view>
    </template>

    <template v-else>
      <view v-if="!articles.length && !loading" class="empty">暂无收藏文章</view>
      <view
        v-for="a in articles"
        :key="a.id"
        class="card"
        :class="{ invalid: a.invalid }"
        @tap="onTapArticle(a)"
      >
        <image class="cover article" :src="a.cover_url || articlePlaceholder" mode="aspectFill" />
        <view class="info">
          <text class="name line-2">{{ a.title || '文章已失效' }}</text>
          <view class="stats">
            <text v-if="a.invalid" class="tag">已失效</text>
            <template v-else>
              <text>阅 {{ a.read_count || 0 }}</text>
              <text>❤ {{ a.like_count || 0 }}</text>
            </template>
          </view>
          <text class="cancel" @tap.stop="cancelArticle(a)">取消收藏</text>
        </view>
      </view>
    </template>

    <view class="more">
      <text v-if="loading">加载中...</text>
      <text v-else-if="finished && currentLen">没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onReachBottom, onShow } from '@dcloudio/uni-app'
import {
  batchRemoveFavorites,
  listFavorites,
  listMyArticleFavorites,
  removeFavorite,
  unfavoriteArticle,
} from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/200/200'
const articlePlaceholder = 'https://picsum.photos/id/102/400/300'
const tab = ref('product')
const list = ref([])
const articles = ref([])
const page = ref(1)
const loading = ref(false)
const finished = ref(false)
const editing = ref(false)
const selected = ref(new Set())
const currentLen = computed(() => (tab.value === 'product' ? list.value.length : articles.value.length))

function switchTab(t) {
  if (tab.value === t) return
  tab.value = t
  editing.value = false
  selected.value = new Set()
  load(true)
}

function toggleEdit() {
  editing.value = !editing.value
  selected.value = new Set()
}

function toggle(id) {
  const s = new Set(selected.value)
  if (s.has(id)) s.delete(id)
  else s.add(id)
  selected.value = s
}

function onTapProduct(item) {
  if (editing.value) {
    toggle(item.product_id)
    return
  }
  if (item.invalid) {
    uni.showToast({ title: '商品已失效', icon: 'none' })
    return
  }
  uni.navigateTo({ url: `/pages/product/detail?id=${item.product_id}` })
}

function onTapArticle(a) {
  if (a.invalid) {
    uni.showToast({ title: '文章已失效', icon: 'none' })
    return
  }
  uni.navigateTo({ url: `/pages/community/detail?id=${a.id}` })
}

async function cancelOne(productId) {
  try {
    await removeFavorite(productId)
    list.value = list.value.filter((x) => x.product_id !== productId)
    uni.showToast({ title: '已取消', icon: 'none' })
  } catch { /* handled */ }
}

async function batchCancel() {
  const ids = [...selected.value]
  if (!ids.length) {
    uni.showToast({ title: '请先勾选', icon: 'none' })
    return
  }
  try {
    await batchRemoveFavorites(ids)
    list.value = list.value.filter((x) => !selected.value.has(x.product_id))
    selected.value = new Set()
    uni.showToast({ title: '已取消', icon: 'none' })
  } catch { /* handled */ }
}

async function cancelArticle(a) {
  try {
    await unfavoriteArticle(a.id)
    articles.value = articles.value.filter((x) => x.id !== a.id)
    uni.showToast({ title: '已取消', icon: 'none' })
  } catch { /* handled */ }
}

async function load(reset = false) {
  if (!isLoggedIn()) {
    uni.redirectTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/favorite/list') })
    return
  }
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      if (tab.value === 'product') list.value = []
      else articles.value = []
    }
    if (tab.value === 'product') {
      const res = await listFavorites({ page: page.value, page_size: 10 })
      const rows = res?.list || []
      list.value = reset ? rows : list.value.concat(rows)
      const total = res?.total || 0
      if (list.value.length >= total || rows.length < 10) finished.value = true
      else page.value += 1
    } else {
      const res = await listMyArticleFavorites({ page: page.value, page_size: 10 })
      const rows = res?.list || []
      articles.value = reset ? rows : articles.value.concat(rows)
      const total = res?.total || 0
      if (articles.value.length >= total || rows.length < 10) finished.value = true
      else page.value += 1
    }
  } catch {
    if (reset) {
      if (tab.value === 'product') list.value = []
      else articles.value = []
    }
  } finally {
    loading.value = false
  }
}

onShow(() => load(true))
onReachBottom(() => load(false))
</script>

<style scoped>
.page { padding: 24rpx 32rpx; }
.tabs {
  display: flex; gap: 8rpx; margin-bottom: 20rpx; background: #f4f4f5;
  border-radius: 999rpx; padding: 6rpx;
}
.tab {
  flex: 1; text-align: center; padding: 14rpx 0; font-size: 26rpx; color: #71717a; border-radius: 999rpx;
}
.tab.active { background: #fff; color: #c8a876; font-weight: 600; box-shadow: 0 2rpx 8rpx rgba(0,0,0,.04); }
.toolbar { display: flex; justify-content: space-between; margin-bottom: 16rpx; font-size: 26rpx; color: #c8a876; }
.danger { color: #d83636; }
.card {
  display: flex; gap: 20rpx; background: #fff; border-radius: 24rpx; padding: 20rpx; margin-bottom: 16rpx;
  align-items: center;
}
.card.invalid { opacity: 0.55; }
.check {
  width: 40rpx; height: 40rpx; border: 2rpx solid #c8a876; border-radius: 8rpx;
  display: flex; align-items: center; justify-content: center; color: #c8a876; font-size: 24rpx;
}
.cover { width: 160rpx; height: 160rpx; border-radius: 16rpx; background: #f4f4f5; flex-shrink: 0; }
.cover.article { height: 120rpx; }
.info { flex: 1; min-width: 0; }
.name { font-size: 28rpx; display: block; margin-bottom: 8rpx; }
.price { color: #d83636; font-weight: 700; font-size: 30rpx; }
.tag { color: #a1a1aa; font-size: 22rpx; }
.stats { display: flex; gap: 16rpx; font-size: 22rpx; color: #a1a1aa; }
.cancel { display: block; margin-top: 12rpx; color: #71717a; font-size: 22rpx; }
.empty, .more { text-align: center; color: #71717a; padding: 40rpx; font-size: 24rpx; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2;
  overflow: hidden;
}
</style>
