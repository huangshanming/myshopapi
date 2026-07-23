<template>
  <view class="page">
    <view class="tip">复制推广链接后，到第三方 App / 浏览器打开购买。</view>

    <view
      v-for="a in items"
      :key="a.act_id"
      class="card"
    >
      <image class="cover" :src="a.img || a.poster || a.icon || placeholder" mode="aspectFill" />
      <view class="body">
        <view class="title-row">
          <image v-if="a.icon" class="icon" :src="a.icon" mode="aspectFit" />
          <text class="title line-2">{{ a.title }}</text>
        </view>
        <text v-if="a.desc" class="desc line-2">{{ a.desc }}</text>
        <view class="meta">
          <text v-if="a.commission_rate_des" class="tag">{{ a.commission_rate_des }}</text>
          <text v-if="a.activity_date" class="date">{{ a.activity_date }}</text>
        </view>
        <button
          class="btn"
          :loading="copyingId === a.act_id"
          :disabled="!!copyingId"
          @tap="onCopy(a)"
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
import { ref } from 'vue'
import { onReachBottom, onShow } from '@dcloudio/uni-app'
import { convertCpsAct, listCpsActs } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const placeholder = 'https://picsum.photos/id/96/400/240'
const items = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)
const copyingId = ref(0)
const emptyText = ref('暂无优惠活动')

async function load(reset = false) {
  if (loading.value || (finished.value && !reset)) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      finished.value = false
      items.value = []
    }
    const res = await listCpsActs({ page: page.value, page_size: pageSize })
    const list = res?.list || []
    items.value = reset ? list : items.value.concat(list)
    const total = Number(res?.total) || 0
    if (items.value.length >= total || list.length < pageSize) {
      finished.value = true
    } else {
      page.value += 1
    }
    if (!items.value.length) emptyText.value = '暂无优惠活动'
  } catch (e) {
    if (reset) {
      items.value = []
      emptyText.value = e.message || '加载失败，请检查聚推客配置'
    }
  } finally {
    loading.value = false
  }
}

function loadMore() {
  load(false)
}

async function onCopy(a) {
  if (!isLoggedIn()) {
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/cps/redpacket'),
    })
    return
  }
  copyingId.value = a.act_id
  try {
    const res = await convertCpsAct({ act_id: Number(a.act_id) })
    const data = res?.data || res
    const link = data?.h5 || data?.long_h5 || ''
    if (!link) {
      uni.showToast({ title: '该活动暂无可用链接', icon: 'none' })
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
    copyingId.value = 0
  }
}

onReachBottom(() => loadMore())
onShow(() => {
  if (!items.value.length) load(true)
})
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; min-height: 100vh; background: #f7f3ec; }
.tip {
  font-size: 22rpx; color: #8b7355; line-height: 1.5;
  background: rgba(200,168,118,.12); border-radius: 16rpx;
  padding: 16rpx 20rpx; margin-bottom: 20rpx;
}
.card {
  display: flex; gap: 20rpx; background: #fff; border-radius: 24rpx;
  padding: 20rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.cover {
  width: 200rpx; height: 200rpx; border-radius: 16rpx; background: #f3f3f3; flex-shrink: 0;
}
.body { flex: 1; min-width: 0; display: flex; flex-direction: column; }
.title-row { display: flex; align-items: flex-start; gap: 10rpx; }
.icon { width: 36rpx; height: 36rpx; border-radius: 8rpx; flex-shrink: 0; margin-top: 4rpx; }
.title { flex: 1; font-size: 28rpx; font-weight: 700; color: #18181b; }
.desc { margin-top: 8rpx; font-size: 22rpx; color: #71717a; }
.meta { margin-top: 12rpx; display: flex; flex-wrap: wrap; gap: 10rpx; align-items: center; }
.tag {
  font-size: 20rpx; color: #8b7355; background: #f7f3ec;
  padding: 4rpx 12rpx; border-radius: 8rpx;
}
.date { font-size: 20rpx; color: #a1a1aa; }
.btn {
  margin-top: auto; margin-left: 0; margin-right: 0;
  height: 64rpx; line-height: 64rpx; border-radius: 999rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; font-size: 26rpx;
}
.btn[disabled] { opacity: .65; }
.line-2 {
  display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 2; overflow: hidden;
}
.empty { text-align: center; color: #a1a1aa; padding: 100rpx 24rpx; font-size: 26rpx; }
.foot { text-align: center; padding: 24rpx; color: #a1a1aa; font-size: 24rpx; }
</style>
