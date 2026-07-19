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
        <text>评论 {{ article.comment_count || commentTotal }}</text>
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

    <view class="panel comments">
      <text class="sec-title">评论 {{ commentTotal }}</text>
      <view v-if="!comments.length && !cmtLoading" class="cmt-empty">暂无评论，来说两句吧</view>
      <view
        v-for="c in comments"
        :id="'cmt-' + c.id"
        :key="c.id"
        class="cmt"
        :class="{ flash: highlightId === c.id }"
      >
        <view class="cmt-hd">
          <text class="cmt-name">{{ c.user_nickname || '用户' }}</text>
          <text class="cmt-time">{{ c.created_at }}</text>
        </view>
        <view class="cmt-body">
          <template v-for="(seg, si) in parseContent(c.content)" :key="si">
            <image v-if="seg.type === 'em'" class="em-inline" :src="seg.url" mode="aspectFit" />
            <text v-else class="cmt-txt">{{ seg.text }}</text>
          </template>
        </view>
        <text class="cmt-reply-btn" @tap="startReply(c)">回复</text>
        <view v-if="c.children?.length" class="cmt-children">
          <view
            v-for="ch in c.children"
            :id="'cmt-' + ch.id"
            :key="ch.id"
            class="cmt child"
            :class="{ flash: highlightId === ch.id }"
          >
            <view class="cmt-hd">
              <text class="cmt-name">{{ ch.user_nickname || '用户' }}</text>
              <text v-if="ch.reply_to_nickname" class="cmt-to">回复 {{ ch.reply_to_nickname }}</text>
              <text class="cmt-time">{{ ch.created_at }}</text>
            </view>
            <view class="cmt-body">
              <template v-for="(seg, si) in parseContent(ch.content)" :key="si">
                <image v-if="seg.type === 'em'" class="em-inline" :src="seg.url" mode="aspectFit" />
                <text v-else class="cmt-txt">{{ seg.text }}</text>
              </template>
            </view>
            <text class="cmt-reply-btn" @tap="startReply(ch, c)">回复</text>
          </view>
        </view>
      </view>
      <view v-if="cmtHasMore" class="more" @tap="loadComments(false)">加载更多</view>
    </view>

    <view class="composer" :class="{ open: composerOpen || replyTarget }">
      <view v-if="replyTarget" class="reply-tip">
        <text>回复 {{ replyTarget.user_nickname || '用户' }}</text>
        <text class="cancel" @tap="cancelReply">取消</text>
      </view>
      <view class="composer-row">
        <input
          class="input"
          v-model="draft"
          :placeholder="replyTarget ? '写回复…' : '写评论…'"
          confirm-type="send"
          @confirm="submitComment"
          @focus="composerOpen = true"
        />
        <text class="em-btn" @tap="toggleEmoji">表情</text>
        <text class="send" @tap="submitComment">发送</text>
      </view>
      <scroll-view v-if="emojiOpen" class="emoji-panel" scroll-y>
        <view class="emoji-grid">
          <view v-for="e in emojis" :key="e.id" class="emoji-item" @tap="insertEmoji(e)">
            <image class="emoji-img" :src="e.image_url" mode="aspectFit" />
            <text class="emoji-name">{{ e.name }}</text>
          </view>
          <view v-if="!emojis.length" class="cmt-empty">暂无表情包</view>
        </view>
      </scroll-view>
    </view>

    <view class="bar">
      <view class="action" @tap="toggleLike">
        <text>{{ liked ? '❤' : '♡' }}</text>
        <text>{{ article.like_count || 0 }}</text>
      </view>
      <view class="action" @tap="focusComposer">
        <text>💬</text>
        <text>{{ article.comment_count || commentTotal }}</text>
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
import { nextTick, reactive, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import {
  createArticleComment,
  favoriteArticle,
  getArticle,
  likeArticle,
  listArticleComments,
  listCommentEmojis,
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
const focusCommentId = ref(0)
const highlightId = ref(0)

const comments = ref([])
const commentTotal = ref(0)
const cmtPage = ref(1)
const cmtLoading = ref(false)
const cmtHasMore = ref(false)

const draft = ref('')
const replyTarget = ref(null)
const replyRoot = ref(null)
const composerOpen = ref(false)
const emojiOpen = ref(false)
const emojis = ref([])
const emojiMap = ref({})

function ensureLogin() {
  if (isLoggedIn()) return true
  uni.navigateTo({
    url: '/pages/login/login?redirect=' + encodeURIComponent(`/pages/community/detail?id=${id.value}`),
  })
  return false
}

function preview(i) {
  const urls = images.value.map((x) => x.url || x).filter(Boolean)
  if (!urls.length) return
  uni.previewImage({ urls, current: urls[i] })
}

function parseContent(content) {
  const text = String(content || '')
  const re = /\[em:(\d+)\]/g
  const segs = []
  let last = 0
  let m
  while ((m = re.exec(text))) {
    if (m.index > last) segs.push({ type: 'text', text: text.slice(last, m.index) })
    const eid = Number(m[1])
    const url = emojiMap.value[eid]
    if (url) segs.push({ type: 'em', url })
    else segs.push({ type: 'text', text: m[0] })
    last = m.index + m[0].length
  }
  if (last < text.length) segs.push({ type: 'text', text: text.slice(last) })
  return segs.length ? segs : [{ type: 'text', text }]
}

async function loadEmojis() {
  try {
    const res = await listCommentEmojis()
    emojis.value = res?.list || []
    const map = {}
    emojis.value.forEach((e) => { map[e.id] = e.image_url })
    emojiMap.value = map
  } catch {
    emojis.value = []
  }
}

async function loadComments(reset = true) {
  if (!id.value || cmtLoading.value) return
  cmtLoading.value = true
  try {
    if (reset) {
      cmtPage.value = 1
      comments.value = []
    }
    const res = await listArticleComments(id.value, { page: cmtPage.value, page_size: 10 })
    const rows = res?.list || []
    comments.value = reset ? rows : comments.value.concat(rows)
    commentTotal.value = Number(res?.total) || 0
    cmtHasMore.value = comments.value.length < commentTotal.value
    if (rows.length) cmtPage.value += 1
  } catch {
    if (reset) comments.value = []
  } finally {
    cmtLoading.value = false
  }
}

async function scrollToComment(cid) {
  if (!cid) return
  await nextTick()
  setTimeout(() => {
    uni.pageScrollTo({
      selector: `#cmt-${cid}`,
      duration: 280,
    })
    highlightId.value = cid
    setTimeout(() => { highlightId.value = 0 }, 2200)
  }, 120)
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
    await loadEmojis()
    await loadComments(true)
    if (focusCommentId.value) scrollToComment(focusCommentId.value)
  } catch {
    Object.keys(article).forEach((k) => delete article[k])
  } finally {
    loading.value = false
  }
}

function startReply(c, root) {
  if (!ensureLogin()) return
  replyTarget.value = c
  replyRoot.value = root || (c.parent_id ? null : c)
  composerOpen.value = true
  emojiOpen.value = false
}

function cancelReply() {
  replyTarget.value = null
  replyRoot.value = null
}

function focusComposer() {
  composerOpen.value = true
}

function toggleEmoji() {
  if (!ensureLogin()) return
  emojiOpen.value = !emojiOpen.value
  composerOpen.value = true
}

function insertEmoji(e) {
  draft.value += `[em:${e.id}]`
}

async function submitComment() {
  if (!ensureLogin()) return
  const content = draft.value.trim()
  if (!content) {
    uni.showToast({ title: '请输入内容', icon: 'none' })
    return
  }
  try {
    const parentId = replyTarget.value?.id || 0
    await createArticleComment(id.value, { content, parent_id: parentId })
    draft.value = ''
    cancelReply()
    emojiOpen.value = false
    article.comment_count = (article.comment_count || 0) + 1
    await loadComments(true)
    uni.showToast({ title: '已发布', icon: 'none' })
  } catch (e) {
    uni.showToast({ title: e.message || '发布失败', icon: 'none' })
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
  focusCommentId.value = Number(q?.comment_id || 0)
  if (!id.value) {
    loading.value = false
    return
  }
  load()
})
</script>

<style scoped>
.page { padding-bottom: 280rpx; background: #fafafa; min-height: 100vh; }
.cover { width: 100%; height: 420rpx; background: #e4e4e7; }
.panel {
  margin: 0 24rpx 24rpx; background: #fff; border-radius: 16rpx;
  padding: 28rpx; margin-top: -32rpx; position: relative; z-index: 1;
}
.comments { margin-top: 0; }
.row { display: flex; gap: 12rpx; align-items: flex-start; }
.title { flex: 1; font-size: 36rpx; font-weight: 700; color: #18181b; line-height: 1.4; }
.badge {
  font-size: 20rpx; color: #c8a876; background: rgba(200,168,118,.15);
  padding: 4rpx 12rpx; border-radius: 6rpx;
}
.stats {
  display: flex; flex-wrap: wrap; gap: 24rpx; margin: 20rpx 0 28rpx;
  font-size: 22rpx; color: #a1a1aa;
}
.content { font-size: 28rpx; color: #3f3f46; line-height: 1.7; }
.content.plain { display: block; white-space: pre-wrap; }
.imgs { margin-top: 24rpx; display: flex; flex-direction: column; gap: 16rpx; }
.img { width: 100%; border-radius: 12rpx; background: #f4f4f5; }
.sec-title { display: block; font-size: 30rpx; font-weight: 700; margin-bottom: 20rpx; color: #18181b; }
.cmt { padding: 20rpx 0; border-bottom: 1rpx solid #f4f4f5; }
.cmt.flash { background: rgba(200,168,118,.12); border-radius: 12rpx; padding: 20rpx 12rpx; }
.cmt-hd { display: flex; align-items: center; flex-wrap: wrap; gap: 12rpx; margin-bottom: 8rpx; }
.cmt-name { font-size: 26rpx; font-weight: 600; color: #3f3f46; }
.cmt-to { font-size: 22rpx; color: #c4894a; }
.cmt-time { font-size: 20rpx; color: #a1a1aa; margin-left: auto; }
.cmt-body { display: flex; flex-wrap: wrap; align-items: center; gap: 4rpx; }
.cmt-txt { font-size: 28rpx; color: #27272a; line-height: 1.5; }
.em-inline { width: 48rpx; height: 48rpx; }
.cmt-reply-btn { display: inline-block; margin-top: 10rpx; font-size: 22rpx; color: #c4894a; }
.cmt-children {
  margin-top: 12rpx; padding: 8rpx 16rpx; background: #fafafa; border-radius: 12rpx;
}
.cmt.child { border-bottom: 1rpx solid #f0f0f0; }
.cmt.child:last-child { border-bottom: none; }
.cmt-empty, .more { text-align: center; color: #a1a1aa; padding: 24rpx 0; font-size: 24rpx; }
.composer {
  position: fixed; left: 0; right: 0; bottom: 100rpx;
  background: #fff; border-top: 1rpx solid #f0f0f0;
  padding: 12rpx 20rpx; padding-bottom: calc(12rpx + env(safe-area-inset-bottom));
  z-index: 20;
}
.reply-tip {
  display: flex; justify-content: space-between; font-size: 22rpx; color: #71717a;
  margin-bottom: 8rpx;
}
.cancel { color: #c4894a; }
.composer-row { display: flex; align-items: center; gap: 12rpx; }
.input {
  flex: 1; height: 68rpx; background: #f4f4f5; border-radius: 34rpx;
  padding: 0 24rpx; font-size: 26rpx;
}
.em-btn, .send { font-size: 26rpx; color: #c4894a; padding: 0 8rpx; white-space: nowrap; }
.emoji-panel { max-height: 280rpx; margin-top: 12rpx; }
.emoji-grid { display: flex; flex-wrap: wrap; gap: 16rpx; }
.emoji-item {
  width: 112rpx; display: flex; flex-direction: column; align-items: center; gap: 6rpx;
}
.emoji-img { width: 72rpx; height: 72rpx; }
.emoji-name { font-size: 20rpx; color: #71717a; }
.bar {
  position: fixed; left: 0; right: 0; bottom: 0;
  display: flex; background: #fff; border-top: 1rpx solid #f4f4f5;
  padding: 20rpx 48rpx calc(20rpx + env(safe-area-inset-bottom));
  z-index: 21;
}
.action {
  flex: 1; display: flex; align-items: center; justify-content: center;
  gap: 12rpx; font-size: 28rpx; color: #3f3f46;
}
.empty { text-align: center; color: #a1a1aa; padding: 120rpx 0; font-size: 28rpx; }
</style>
