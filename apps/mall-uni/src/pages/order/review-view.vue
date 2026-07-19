<template>
  <view class="page" v-if="review">
    <view class="card">
      <view class="stars">
        <text v-for="n in 5" :key="n" class="star" :class="{ on: n <= review.rating }">★</text>
      </view>
      <text class="content">{{ review.content || '（无文字评价）' }}</text>
      <view v-if="review.images?.length" class="imgs">
        <image
          v-for="img in review.images"
          :key="img.id || img.url"
          :src="img.url"
          mode="aspectFill"
          class="img"
          @tap="preview(img.url)"
        />
      </view>
      <text class="meta">{{ review.is_anonymous ? '匿名用户' : '我' }} · {{ review.created_at }}</text>
    </view>
    <view v-if="review.merchant_reply" class="card reply">
      <text class="sec">商家回复</text>
      <text class="content">{{ review.merchant_reply }}</text>
      <text v-if="review.replied_at" class="meta">{{ review.replied_at }}</text>
    </view>
  </view>
  <view v-else class="empty">{{ loading ? '加载中...' : '暂无评价' }}</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getOrderReview } from '../../api/index'

const review = ref(null)
const loading = ref(false)
let orderId = 0

onLoad((q) => {
  orderId = Number(q.id || 0)
  load()
})

async function load() {
  if (!orderId) return
  loading.value = true
  try {
    const res = await getOrderReview(orderId)
    review.value = res || null
  } catch {
    review.value = null
  } finally {
    loading.value = false
  }
}

function preview(url) {
  const urls = (review.value?.images || []).map((i) => i.url)
  uni.previewImage({ current: url, urls })
}
</script>

<style scoped>
.page { padding: 24rpx 32rpx; }
.card {
  background: #fff; border-radius: 24rpx; padding: 28rpx; margin-bottom: 20rpx;
}
.stars { margin-bottom: 16rpx; }
.star { font-size: 36rpx; color: #e4e4e7; }
.star.on { color: #f5a623; }
.content { font-size: 28rpx; line-height: 1.6; color: #3f3f46; display: block; }
.imgs { display: flex; flex-wrap: wrap; gap: 12rpx; margin-top: 16rpx; }
.img { width: 180rpx; height: 180rpx; border-radius: 12rpx; }
.meta { display: block; margin-top: 16rpx; color: #a1a1aa; font-size: 22rpx; }
.sec { font-weight: 600; display: block; margin-bottom: 12rpx; }
.reply { background: #faf8f4; }
.empty { text-align: center; padding: 120rpx; color: #71717a; }
</style>
