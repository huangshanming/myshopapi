<template>
  <view class="page">
    <scroll-view scroll-y class="scroll" :show-scrollbar="false">
      <!-- 图片区：小红书式九宫格 -->
      <view class="photo-section">
        <view class="photo-grid">
          <view
            v-for="(url, i) in photos"
            :key="url + i"
            class="photo-item"
            @tap="previewPhoto(i)"
            @longpress="setCover(i)"
          >
            <image class="photo-img" :src="url" mode="aspectFill" />
            <view v-if="coverIndex === i" class="cover-badge">封面</view>
            <view class="photo-del" @tap.stop="removePhoto(i)">×</view>
          </view>
          <view v-if="photos.length < 9" class="photo-add" @tap="pickPhotos">
            <text class="add-plus">+</text>
            <text class="add-text">添加</text>
          </view>
        </view>
        <text class="photo-tip">最多 9 张 · 长按图片设为封面</text>
      </view>

      <view class="divider" />

      <!-- 标题 -->
      <view class="title-row">
        <input
          v-model="form.title"
          class="title-input"
          maxlength="80"
          placeholder="填写标题，会有更多赞哦~"
          placeholder-class="ph"
        />
        <text class="counter">{{ form.title.length }}/80</text>
      </view>

      <view class="divider thin" />

      <!-- 正文 -->
      <view class="body-row">
        <textarea
          v-model="bodyText"
          class="body-input"
          maxlength="5000"
          :auto-height="true"
          :show-confirm-bar="false"
          placeholder="分享你的真实体验，添加 #话题# 更容易被看到"
          placeholder-class="ph"
        />
      </view>

      <view v-if="topicTags.length" class="topic-row">
        <text v-for="(t, i) in topicTags" :key="i" class="topic-tag">#{{ t }}</text>
      </view>

      <view class="hint-block">
        <text class="hint">提交后需平台审核通过才会公开展示；审核通过可在任务中心领取奖励。</text>
      </view>
    </scroll-view>

    <!-- 底部发布栏 -->
    <view class="bottom-bar">
      <view class="bar-left">
        <text class="bar-icon">📝</text>
        <text class="bar-text">{{ photos.length ? `已选 ${photos.length} 张` : '记得添加图片' }}</text>
      </view>
      <button class="publish-btn" :loading="saving" :disabled="saving" @tap="submit">发布笔记</button>
    </view>
  </view>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { createMyArticle, getMyArticle, updateMyArticle, uploadArticleImage } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const editId = ref(0)
const saving = ref(false)
const photos = ref([])
const coverIndex = ref(0)
const bodyText = ref('')
const form = reactive({ title: '' })

const topicTags = computed(() => {
  const m = bodyText.value.match(/#([^#\s]{1,20})#?/g) || []
  const seen = new Set()
  const out = []
  for (const raw of m) {
    const t = raw.replace(/^#|#?$/g, '')
    if (t && !seen.has(t)) {
      seen.add(t)
      out.push(t)
    }
  }
  return out.slice(0, 8)
})

async function ensureLogin() {
  if (isLoggedIn()) return true
  uni.navigateTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/community/publish') })
  return false
}

async function uploadOne(path) {
  const res = await uploadArticleImage(path)
  return res?.url || ''
}

function escapeHtml(s) {
  return String(s || '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function formatLine(line) {
  return escapeHtml(line).replace(
    /#([^#\s]{1,20})#?/g,
    '<span style="color:#ff2442;font-weight:600">#$1</span>',
  )
}

function plainToHtml(text) {
  const t = String(text || '').trim()
  if (!t) return ''
  return t
    .split(/\n{2,}/)
    .map((p) => {
      const inner = p.split('\n').map((line) => formatLine(line)).join('<br/>')
      return `<p style="margin:0 0 16px;line-height:1.75;color:#333;font-size:15px">${inner}</p>`
    })
    .join('')
}

function htmlToPlain(html) {
  return String(html || '')
    .replace(/<br\s*\/?>/gi, '\n')
    .replace(/<\/p>\s*<p[^>]*>/gi, '\n\n')
    .replace(/<\/p>/gi, '\n')
    .replace(/<p[^>]*>/gi, '')
    .replace(/<[^>]+>/g, '')
    .replace(/&nbsp;/g, ' ')
    .replace(/&amp;/g, '&')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&quot;/g, '"')
    .replace(/\n{3,}/g, '\n\n')
    .trim()
}

function syncCoverIndex() {
  if (!photos.value.length) {
    coverIndex.value = 0
    return
  }
  if (coverIndex.value >= photos.value.length) {
    coverIndex.value = 0
  }
}

function setCover(i) {
  if (i >= 0 && i < photos.value.length) {
    coverIndex.value = i
    uni.showToast({ title: '已设为封面', icon: 'none' })
  }
}

function removePhoto(i) {
  photos.value.splice(i, 1)
  syncCoverIndex()
}

function previewPhoto(i) {
  if (!photos.value.length) return
  uni.previewImage({ urls: photos.value, current: photos.value[i] })
}

function pickPhotos() {
  const left = 9 - photos.value.length
  if (left <= 0) return
  uni.chooseImage({
    count: left,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (r) => {
      uni.showLoading({ title: '上传中' })
      try {
        for (const p of r.tempFilePaths) {
          const url = await uploadOne(p)
          if (url) photos.value.push(url)
        }
        syncCoverIndex()
      } catch (e) {
        uni.showToast({ title: e.message || '上传失败', icon: 'none' })
      } finally {
        uni.hideLoading()
      }
    },
  })
}

async function submit() {
  if (!(await ensureLogin())) return
  if (!photos.value.length) {
    uni.showToast({ title: '请至少添加一张图片', icon: 'none' })
    return
  }
  if (!form.title.trim()) {
    uni.showToast({ title: '请填写标题', icon: 'none' })
    return
  }
  if (!bodyText.value.trim()) {
    uni.showToast({ title: '请填写正文', icon: 'none' })
    return
  }

  syncCoverIndex()
  const cover = photos.value[coverIndex.value]
  const others = photos.value.filter((_, i) => i !== coverIndex.value)

  saving.value = true
  try {
    const payload = {
      title: form.title.trim(),
      cover_url: cover,
      content: plainToHtml(bodyText.value),
      image_urls: others,
    }
    if (editId.value) await updateMyArticle(editId.value, payload)
    else await createMyArticle(payload)
    uni.showToast({ title: '已提交审核', icon: 'none' })
    setTimeout(() => {
      uni.redirectTo({ url: '/pages/community/mine' })
    }, 500)
  } catch (e) {
    uni.showToast({ title: e.message || '提交失败', icon: 'none' })
  } finally {
    saving.value = false
  }
}

onLoad(async (q) => {
  if (!(await ensureLogin())) return
  editId.value = Number(q?.id || 0)
  if (!editId.value) return
  try {
    const res = await getMyArticle(editId.value)
    const a = res?.article || {}
    form.title = a.title || ''
    bodyText.value = htmlToPlain(a.content || '')

    const urls = []
    if (a.cover_url) urls.push(a.cover_url)
    for (const im of res?.images || []) {
      const u = im?.url || im
      if (u && !urls.includes(u)) urls.push(u)
    }
    photos.value = urls
    coverIndex.value = a.cover_url ? Math.max(0, urls.indexOf(a.cover_url)) : 0
  } catch {
    uni.showToast({ title: '加载失败', icon: 'none' })
  }
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #fff;
  display: flex;
  flex-direction: column;
}

.scroll {
  flex: 1;
  height: 0;
  padding-bottom: calc(140rpx + env(safe-area-inset-bottom));
}

.photo-section {
  padding: 28rpx 32rpx 8rpx;
}

.photo-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}

.photo-item,
.photo-add {
  position: relative;
  width: calc((100% - 24rpx) / 3);
  aspect-ratio: 1;
  border-radius: 16rpx;
  overflow: hidden;
}

.photo-img {
  width: 100%;
  height: 100%;
  background: #f5f5f5;
}

.cover-badge {
  position: absolute;
  left: 8rpx;
  bottom: 8rpx;
  padding: 2rpx 12rpx;
  border-radius: 8rpx;
  font-size: 20rpx;
  color: #fff;
  background: rgba(0, 0, 0, 0.45);
}

.photo-del {
  position: absolute;
  top: 8rpx;
  right: 8rpx;
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.55);
  color: #fff;
  font-size: 28rpx;
  line-height: 40rpx;
  text-align: center;
}

.photo-add {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #fafafa;
  border: 2rpx dashed #e5e5e5;
  box-sizing: border-box;
}

.add-plus {
  font-size: 52rpx;
  color: #ccc;
  line-height: 1;
  font-weight: 300;
}

.add-text {
  margin-top: 6rpx;
  font-size: 22rpx;
  color: #bbb;
}

.photo-tip {
  display: block;
  margin-top: 16rpx;
  font-size: 22rpx;
  color: #bbb;
}

.divider {
  height: 16rpx;
  background: #fafafa;
}

.divider.thin {
  height: 1rpx;
  background: #f0f0f0;
}

.title-row {
  padding: 28rpx 32rpx 16rpx;
}

.title-input {
  width: 100%;
  font-size: 34rpx;
  font-weight: 600;
  color: #222;
  line-height: 1.4;
}

.counter {
  display: block;
  margin-top: 8rpx;
  text-align: right;
  font-size: 22rpx;
  color: #ccc;
}

.body-row {
  padding: 24rpx 32rpx 32rpx;
  min-height: 360rpx;
}

.body-input {
  width: 100%;
  min-height: 320rpx;
  font-size: 30rpx;
  line-height: 1.75;
  color: #333;
}

.ph {
  color: #ccc;
}

.topic-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  padding: 0 32rpx 24rpx;
}

.topic-tag {
  padding: 8rpx 20rpx;
  border-radius: 999rpx;
  font-size: 24rpx;
  color: #ff2442;
  background: rgba(255, 36, 66, 0.08);
}

.hint-block {
  padding: 0 32rpx 40rpx;
}

.hint {
  font-size: 22rpx;
  color: #bbb;
  line-height: 1.6;
}

.bottom-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 28rpx calc(16rpx + env(safe-area-inset-bottom));
  background: #fff;
  border-top: 1rpx solid #f0f0f0;
  box-shadow: 0 -8rpx 24rpx rgba(0, 0, 0, 0.04);
}

.bar-left {
  display: flex;
  align-items: center;
  gap: 10rpx;
  flex: 1;
  min-width: 0;
}

.bar-icon {
  font-size: 32rpx;
}

.bar-text {
  font-size: 24rpx;
  color: #999;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.publish-btn {
  margin: 0;
  padding: 0 48rpx;
  height: 72rpx;
  line-height: 72rpx;
  border-radius: 999rpx;
  font-size: 28rpx;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #ff6b81, #ff2442);
  border: none;
}

.publish-btn[disabled] {
  opacity: 0.65;
}

.publish-btn::after {
  border: none;
}
</style>
