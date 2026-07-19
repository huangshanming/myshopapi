<template>
  <view class="page" v-if="order">
    <view class="card">
      <text class="sec">评价商品</text>
      <radio-group @change="onItemChange">
        <label v-for="it in order.items || []" :key="it.id" class="item">
          <radio :value="String(it.id)" :checked="orderItemId === it.id" color="#c8a876" />
          <view class="item-info">
            <text class="name">{{ it.product_name }}</text>
            <text class="sub">×{{ it.quantity }} · ¥{{ it.price }}</text>
          </view>
        </label>
      </radio-group>
    </view>

    <view class="card">
      <text class="sec">评分</text>
      <view class="stars">
        <text
          v-for="n in 5"
          :key="n"
          class="star"
          :class="{ on: n <= rating }"
          @tap="rating = n"
        >★</text>
      </view>
      <textarea
        v-model="content"
        class="textarea"
        maxlength="1000"
        placeholder="分享你的使用感受（选填）"
      />
      <view class="imgs">
        <view v-for="(url, i) in images" :key="url" class="img-wrap">
          <image :src="url" mode="aspectFill" class="img" />
          <text class="rm" @tap="images.splice(i, 1)">×</text>
        </view>
        <view v-if="images.length < 9" class="add" @tap="pickImage">+</view>
      </view>
      <label class="anon">
        <switch :checked="isAnonymous" color="#c8a876" @change="(e) => (isAnonymous = e.detail.value)" />
        <text>匿名评价</text>
      </label>
    </view>

    <button class="btn" :loading="submitting" @tap="submit">提交评价</button>
  </view>
  <view v-else class="empty">{{ loading ? '加载中...' : '订单不存在' }}</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getOrder, submitOrderReview, uploadReviewImage } from '../../api/index'

const order = ref(null)
const loading = ref(false)
const submitting = ref(false)
const rating = ref(5)
const content = ref('')
const isAnonymous = ref(false)
const images = ref([])
const orderItemId = ref(0)
let orderId = 0

onLoad((q) => {
  orderId = Number(q.id || 0)
  load()
})

function onItemChange(e) {
  orderItemId.value = Number(e.detail.value)
}

async function load() {
  if (!orderId) return
  loading.value = true
  try {
    const res = await getOrder(orderId)
    order.value = res || null
    if (order.value?.status !== 'completed') {
      uni.showToast({ title: '当前订单不可评价', icon: 'none' })
    }
    const items = order.value?.items || []
    if (items.length) orderItemId.value = items[0].id
  } catch {
    order.value = null
  } finally {
    loading.value = false
  }
}

function pickImage() {
  uni.chooseImage({
    count: 9 - images.value.length,
    success: async (res) => {
      for (const path of res.tempFilePaths || []) {
        try {
          const up = await uploadReviewImage(path)
          if (up?.url) images.value.push(up.url)
        } catch {
          /* toast in upload */
        }
      }
    },
  })
}

async function submit() {
  if (!orderItemId.value) {
    uni.showToast({ title: '请选择商品', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    await submitOrderReview(orderId, {
      rating: rating.value,
      content: content.value,
      is_anonymous: isAnonymous.value,
      order_item_id: orderItemId.value,
      images: images.value,
    })
    uni.showToast({ title: '评价成功', icon: 'success' })
    setTimeout(() => {
      uni.redirectTo({ url: `/pages/order/review-view?id=${orderId}` })
    }, 500)
  } catch {
    /* handled */
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; }
.card {
  background: #fff; border-radius: 24rpx; padding: 28rpx; margin-bottom: 20rpx;
}
.sec { font-weight: 600; font-size: 28rpx; display: block; margin-bottom: 16rpx; }
.item { display: flex; align-items: flex-start; gap: 12rpx; padding: 12rpx 0; }
.item-info { flex: 1; }
.name { display: block; font-size: 28rpx; }
.sub { color: #71717a; font-size: 22rpx; }
.stars { display: flex; gap: 8rpx; margin-bottom: 20rpx; }
.star { font-size: 48rpx; color: #e4e4e7; }
.star.on { color: #f5a623; }
.textarea {
  width: 100%; min-height: 180rpx; background: #fafafa; border-radius: 16rpx;
  padding: 20rpx; font-size: 26rpx; box-sizing: border-box;
}
.imgs { display: flex; flex-wrap: wrap; gap: 16rpx; margin-top: 20rpx; }
.img-wrap { position: relative; width: 160rpx; height: 160rpx; }
.img { width: 100%; height: 100%; border-radius: 12rpx; }
.rm {
  position: absolute; top: -8rpx; right: -8rpx; width: 36rpx; height: 36rpx;
  background: #52525b; color: #fff; border-radius: 50%; text-align: center; line-height: 36rpx; font-size: 24rpx;
}
.add {
  width: 160rpx; height: 160rpx; border: 2rpx dashed #d4d4d8; border-radius: 12rpx;
  display: flex; align-items: center; justify-content: center; font-size: 56rpx; color: #a1a1aa;
}
.anon { display: flex; align-items: center; gap: 12rpx; margin-top: 24rpx; font-size: 26rpx; color: #52525b; }
.btn {
  margin-top: 24rpx; background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff;
  border-radius: 999rpx; height: 88rpx; line-height: 88rpx; font-size: 30rpx;
}
.empty { text-align: center; padding: 120rpx; color: #71717a; }
</style>
