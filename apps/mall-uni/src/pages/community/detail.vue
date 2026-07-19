<template>
  <view class="page" v-if="article.id">
    <image
      v-if="article.cover_url"
      class="cover"
      :src="article.cover_url"
      mode="aspectFill"
    />
    <view class="panel">
      <view class="row">
        <text class="title">{{ article.title }}</text>
        <text v-if="paid" class="badge">推广</text>
      </view>
      <view class="stats">
        <text>观众 {{ article.audience_count || 0 }}</text>
        <text>阅读 {{ article.read_count || 0 }}</text>
        <text>收藏 {{ article.collect_count || 0 }}</text>
      </view>
      <rich-text v-if="article.content" class="content" :nodes="article.content" />
      <text v-else class="content plain">暂无正文</text>
      <view v-if="images.length" class="imgs">
        <image
          v-for="(im, i) in images"
          :key="i"
          class="img"
          :src="im.url || im"
          mode="widthFix"
          @tap="preview(i)"
        />
      </view>
    </view>

    <view class="bar">
      <view class="action" @tap="toggleLike">
        <text>{{ liked ? '❤' : '♡' }}</text>
        <text>{{ article.like_count || 0 }}</text>
      </view>
      <view class="action" @tap="toggleFavorite">
        <text>{{ favorited ? '★' : '☆' }}</text>
        <text>{{ favorited ? '已收藏' : '收藏' }}</text>
      </view>
    </view>
  </view>
  <view v-else class="empty">{{ loading ? '加载中...' : '文章不存在' }}</view>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import {
  favoriteArticle,
  getArticle,
  likeArticle,
  unfavoriteArticle,
  unlikeArticle,
} from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const loading = ref(true)
const article = reactive({})
const images = ref([])
const liked = ref(false)
const favorited = ref(false)
const paid = ref(false)
const id = ref(0)

function ensureLogin() {
  if (isLoggedIn()) return true
  uni.navigateTo({ url: '/pages/login/login' })
  return false
}

function preview(i) {
  const urls = images.value.map((x) => x.url || x).filter(Boolean)
  if (!urls.length) return
  uni.previewImage({ urls, current: urls[i] })
}

async function load() {
  loading.value = true
  try {
    const res = await getArticle(id.value)
    const data = res || {}
    Object.assign(article, data.article || {})
    images.value = data.images || []
    liked.value = !!data.liked
    favorited.value = !!data.favorited
    paid.value = !!data.paid
    if (article.title) uni.setNavigationBarTitle({ title: article.title })
  } catch {
    Object.keys(article).forEach((k) => delete article[k])
  } finally {
    loading.value = false
  }
}

async function toggleLike() {
  if (!ensureLogin()) return
  try {
    if (liked.value) {
      await unlikeArticle(id.value)
      liked.value = false
      article.like_count = Math.max(0, (article.like_count || 1) - 1)
    } else {
      await likeArticle(id.value)
      liked.value = true
      article.like_count = (article.like_count || 0) + 1
    }
  } catch { /* handled */ }
}

async function toggleFavorite() {
  if (!ensureLogin()) return
  try {
    if (favorited.value) {
      await unfavoriteArticle(id.value)
      favorited.value = false
      article.collect_count = Math.max(0, (article.collect_count || 1) - 1)
    } else {
      await favoriteArticle(id.value)
      favorited.value = true
      article.collect_count = (article.collect_count || 0) + 1
    }
  } catch { /* handled */ }
}

onLoad((q) => {
  id.value = Number(q?.id || 0)
  if (!id.value) {
    loading.value = false
    return
  }
  load()
})
</script>

<style scoped>
.page { padding-bottom: 140rpx; background: #fafafa; min-height: 100vh; }
.cover { width: 100%; height: 420rpx; background: #e4e4e7; }
.panel {
  margin: 0 24rpx 24rpx; background: #fff; border-radius: 16rpx;
  padding: 28rpx; margin-top: -32rpx; position: relative; z-index: 1;
}
.row { display: flex; gap: 12rpx; align-items: flex-start; }
.title { flex: 1; font-size: 36rpx; font-weight: 700; color: #18181b; line-height: 1.4; }
.badge {
  font-size: 20rpx; color: #c8a876; background: rgba(200,168,118,.15);
  padding: 4rpx 12rpx; border-radius: 6rpx;
}
.stats {
  display: flex; gap: 24rpx; margin: 20rpx 0 28rpx;
  font-size: 22rpx; color: #a1a1aa;
}
.content { font-size: 28rpx; color: #3f3f46; line-height: 1.7; }
.content.plain { display: block; white-space: pre-wrap; }
.imgs { margin-top: 24rpx; display: flex; flex-direction: column; gap: 16rpx; }
.img { width: 100%; border-radius: 12rpx; background: #f4f4f5; }
.bar {
  position: fixed; left: 0; right: 0; bottom: 0;
  display: flex; background: #fff; border-top: 1rpx solid #f4f4f5;
  padding: 20rpx 48rpx calc(20rpx + env(safe-area-inset-bottom));
}
.action {
  flex: 1; display: flex; align-items: center; justify-content: center;
  gap: 12rpx; font-size: 28rpx; color: #3f3f46;
}
.empty { text-align: center; color: #a1a1aa; padding: 120rpx 0; font-size: 28rpx; }
</style>
