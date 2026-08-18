<template>
  <view class="page">
    <scroll-view class="tabs" scroll-x :show-scrollbar="false">
      <view
        v-for="t in platforms"
        :key="t.v"
        class="tab"
        :class="{ on: platform === t.v }"
        @tap="onTab(t.v)"
      >{{ t.l }}</view>
    </scroll-view>

    <view class="search-row">
      <input
        class="search"
        v-model="keyword"
        confirm-type="search"
        :placeholder="`搜索${platformLabel}好物`"
        @confirm="onSearch"
      />
      <text class="search-btn" @tap="onSearch">搜索</text>
    </view>

    <view class="tip">输入关键词搜索后，复制推广链接到对应 App / 浏览器购买。</view>

    <view
      v-for="g in items"
      :key="g.item_id + g.title"
      class="card"
    >
      <image class="cover" :src="g.cover || placeholder" mode="aspectFill" />
      <view class="body">
        <text class="title line-2">{{ g.title }}</text>
        <view class="price-row">
          <text v-if="g.coupon_price" class="price">¥{{ g.coupon_price }}</text>
          <text v-else-if="g.price" class="price">¥{{ g.price }}</text>
          <text v-if="g.coupon_price && g.price" class="old">¥{{ g.price }}</text>
          <text v-if="g.commission_rate" class="rate">佣 {{ g.commission_rate }}</text>
        </view>
        <button
          class="btn"
          :loading="copyingId === g.item_id"
          :disabled="!!copyingId"
          @tap="onCopy(g)"
        >复制推广链接</button>
      </view>
    </view>

    <view v-if="!loading && !items.length" class="empty">{{ emptyText }}</view>
    <view class="foot">
      <text v-if="loading">加载中...</text>
      <text v-else-if="finished && items.length">没有更多了</text>
      <text v-else-if="items.length" @tap="loadMore">上拉加载更多</text>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onReachBottom } from '@dcloudio/uni-app'
import { convertCpsGoods, listCpsGoods } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/400/400'
const platforms = [
  { v: 'taobao', l: '淘宝' },
  { v: 'jd', l: '京东' },
  { v: 'pdd', l: '拼多多' },
  { v: 'vip', l: '唯品会' },
]

const platform = ref('taobao')
const keyword = ref('')
const appliedKeyword = ref('')
const items = ref([])
const nextMinId = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)
const copyingId = ref('')
const emptyText = ref('请输入关键词搜索')

const platformLabel = computed(() => {
  return platforms.find((p) => p.v === platform.value)?.l || ''
})

function onTab(v) {
  if (platform.value === v) return
  platform.value = v
  items.value = []
  finished.value = false
  nextMinId.value = 1
  emptyText.value = appliedKeyword.value ? '暂无商品' : '请输入关键词搜索'
  if (appliedKeyword.value) {
    load(true)
  }
}

function onSearch() {
  const kw = String(keyword.value || '').trim()
  if (!kw) {
    uni.showToast({ title: '请输入关键词搜索', icon: 'none' })
    return
  }
  appliedKeyword.value = kw
  load(true)
}

async function load(reset = false) {
  if (!appliedKeyword.value) {
    emptyText.value = '请输入关键词搜索'
    return
  }
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      nextMinId.value = 1
      finished.value = false
      items.value = []
    }
    const res = await listCpsGoods({
      platform: platform.value,
      keyword: appliedKeyword.value,
      min_id: nextMinId.value || 1,
      page_size: pageSize,
    })
    const list = res?.list || []
    items.value = reset ? list : items.value.concat(list)
    const next = Number(res?.next_min_id) || 0
    if (!list.length || !next || next === nextMinId.value) {
      finished.value = true
      nextMinId.value = 0
    } else {
      nextMinId.value = next
    }
    if (!items.value.length) emptyText.value = '暂无商品'
  } catch (e) {
    if (reset) {
      items.value = []
      emptyText.value = e.message || '加载失败'
    }
  } finally {
    loading.value = false
  }
}

function loadMore() {
  load(false)
}

async function onCopy(g) {
  if (!isLoggedIn()) {
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/cps/products'),
    })
    return
  }
  copyingId.value = g.item_id || g.title
  try {
    const res = await convertCpsGoods({
      platform: platform.value,
      item_id: g.item_id,
      raw_url: g.raw_url || undefined,
      title: g.title || undefined,
    })
    const data = res?.data || res
    const link = data?.h5 || data?.long_h5 || ''
    if (!link) {
      uni.showToast({ title: '该商品暂无可用链接', icon: 'none' })
      return
    }
    await new Promise((resolve, reject) => {
      uni.setClipboardData({
        data: link,
        success: resolve,
        fail: reject,
      })
    })
    uni.showToast({ title: '已复制，去第三方打开购买', icon: 'none' })
  } catch (e) {
    uni.showToast({ title: e.message || '转链失败', icon: 'none' })
  } finally {
    copyingId.value = ''
  }
}

onReachBottom(() => loadMore())
</script>

<style scoped>
.page { padding: 0 0 48rpx; min-height: 100vh; background: #f7f3ec; }
.tabs {
  white-space: nowrap; padding: 20rpx 32rpx 0;
  background: linear-gradient(180deg, #efe4d2 0%, #f7f3ec 100%);
}
.tab {
  display: inline-block; margin-right: 16rpx; padding: 12rpx 28rpx;
  border-radius: 999rpx; font-size: 26rpx; color: #8b7355; background: rgba(255,255,255,.75);
}
.tab.on { background: #c8a876; color: #fff; font-weight: 600; }
.search-row {
  display: flex; gap: 16rpx; align-items: center;
  padding: 20rpx 32rpx 8rpx;
}
.search {
  flex: 1; height: 72rpx; background: #fff; border-radius: 999rpx;
  padding: 0 28rpx; font-size: 26rpx;
}
.search-btn {
  font-size: 26rpx; color: #fff; background: #c8a876;
  padding: 16rpx 28rpx; border-radius: 999rpx;
}
.tip {
  margin: 8rpx 32rpx 16rpx; font-size: 22rpx; color: #8b7355; line-height: 1.5;
  background: rgba(200,168,118,.12); border-radius: 16rpx; padding: 14rpx 18rpx;
}
.card {
  display: flex; gap: 20rpx; margin: 0 32rpx 20rpx; background: #fff;
  border-radius: 24rpx; padding: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.cover {
  width: 200rpx; height: 200rpx; border-radius: 16rpx; background: #f3f3f3; flex-shrink: 0;
}
.body { flex: 1; min-width: 0; display: flex; flex-direction: column; }
.title { font-size: 28rpx; font-weight: 600; color: #18181b; }
.price-row {
  margin-top: 12rpx; display: flex; align-items: baseline; flex-wrap: wrap; gap: 10rpx;
}
.price { font-size: 32rpx; font-weight: 700; color: #d83636; }
.old { font-size: 22rpx; color: #a1a1aa; text-decoration: line-through; }
.rate {
  font-size: 20rpx; color: #8b7355; background: #f7f3ec;
  padding: 2rpx 10rpx; border-radius: 8rpx;
}
.btn {
  margin-top: auto; margin-left: 0; margin-right: 0;
  height: 64rpx; line-height: 64rpx; border-radius: 999rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; font-size: 26rpx;
}
.btn[disabled] { opacity: .65; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2; overflow: hidden;
}
.empty { text-align: center; color: #a1a1aa; padding: 100rpx 32rpx; font-size: 26rpx; }
.foot { text-align: center; padding: 24rpx; color: #a1a1aa; font-size: 24rpx; }
</style>
