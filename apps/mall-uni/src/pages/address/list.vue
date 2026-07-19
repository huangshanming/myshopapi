<template>
  <view class="page">
    <view v-if="!list.length && !loading" class="empty">暂无收货地址</view>
    <view
      v-for="a in list"
      :key="a.id"
      class="card"
      :class="{ active: fromConfirm && Number(a.id) === pickedId }"
      @tap="onPick(a)"
    >
      <view class="top">
        <text class="name">{{ a.receiver_name }}</text>
        <text class="phone">{{ a.receiver_phone }}</text>
        <text v-if="a.is_default" class="tag">默认</text>
      </view>
      <text class="line">{{ fullAddr(a) }}</text>
      <view class="ops" @tap.stop>
        <text @tap="setDefault(a)">设为默认</text>
        <text @tap="edit(a)">编辑</text>
        <text class="danger" @tap="remove(a)">删除</text>
      </view>
    </view>

    <button class="btn" @tap="add">新增地址</button>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { deleteAddress, listAddresses, setDefaultAddress } from '../../api/index'
import { isLoggedIn } from '../../stores/user'

const list = ref([])
const loading = ref(false)
const fromConfirm = ref(false)
const pickedId = ref(0)

function fullAddr(a) {
  return `${a.province || ''}${a.city || ''}${a.district || ''}${a.detail || ''}`
}

async function load() {
  if (!isLoggedIn()) {
    uni.redirectTo({ url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/address/list') })
    return
  }
  loading.value = true
  try {
    const res = await listAddresses()
    list.value = res || []
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

function onPick(a) {
  if (!fromConfirm.value) return
  uni.setStorageSync('mymall_picked_address_id', String(a.id))
  uni.navigateBack()
}

function add() {
  uni.navigateTo({ url: '/pages/address/edit' })
}

function edit(a) {
  uni.navigateTo({ url: `/pages/address/edit?id=${a.id}` })
}

async function setDefault(a) {
  try {
    await setDefaultAddress(a.id)
    uni.showToast({ title: '已设为默认', icon: 'success' })
    load()
  } catch {
    /* handled */
  }
}

async function remove(a) {
  const ok = await new Promise((resolve) => {
    uni.showModal({
      title: '删除地址',
      content: '确认删除该收货地址？',
      success: (r) => resolve(r.confirm),
    })
  })
  if (!ok) return
  try {
    await deleteAddress(a.id)
    uni.showToast({ title: '已删除', icon: 'success' })
    load()
  } catch {
    /* handled */
  }
}

onLoad((q) => {
  fromConfirm.value = q.from === 'confirm'
})

onShow(load)
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; }
.empty { text-align: center; color: #a1a1aa; padding: 80rpx 0; }
.card {
  background: #fff; border-radius: 24rpx; padding: 28rpx; margin-bottom: 20rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.card.active { outline: 2rpx solid #c8a876; }
.top { display: flex; align-items: center; gap: 16rpx; margin-bottom: 8rpx; }
.name { font-weight: 600; font-size: 30rpx; }
.phone { color: #71717a; font-size: 26rpx; }
.tag {
  font-size: 20rpx; color: #c8a876; border: 1rpx solid #c8a876;
  padding: 0 8rpx; border-radius: 6rpx;
}
.line { display: block; color: #52525b; font-size: 26rpx; line-height: 1.5; }
.ops {
  display: flex; gap: 32rpx; margin-top: 20rpx; padding-top: 16rpx;
  border-top: 1rpx solid #f5f5f5; color: #71717a; font-size: 24rpx;
}
.danger { color: #d83636; }
.btn {
  margin-top: 24rpx; height: 88rpx; line-height: 88rpx; border-radius: 999rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; font-size: 30rpx;
}
</style>
