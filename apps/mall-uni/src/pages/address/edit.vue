<template>
  <view class="page">
    <view class="card">
      <view class="field">
        <text class="label">收货人</text>
        <input v-model="form.receiver_name" class="input" placeholder="姓名" />
      </view>
      <view class="field">
        <text class="label">手机号</text>
        <input v-model="form.receiver_phone" class="input" type="number" maxlength="11" placeholder="手机号" />
      </view>
      <view class="field">
        <text class="label">省</text>
        <input v-model="form.province" class="input" placeholder="省" />
      </view>
      <view class="field">
        <text class="label">市</text>
        <input v-model="form.city" class="input" placeholder="市" />
      </view>
      <view class="field">
        <text class="label">区</text>
        <input v-model="form.district" class="input" placeholder="区/县" />
      </view>
      <view class="field">
        <text class="label">详细地址</text>
        <input v-model="form.detail" class="input" placeholder="街道门牌等" />
      </view>
      <view class="field switch-row" @tap="form.is_default = form.is_default ? 0 : 1">
        <text class="label">设为默认</text>
        <text class="box" :class="{ on: form.is_default }">{{ form.is_default ? '✓' : '' }}</text>
      </view>
    </view>
    <button class="btn" :loading="saving" @tap="save">保存</button>
  </view>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { createAddress, listAddresses, updateAddress } from '../../api/index'

const form = reactive({
  receiver_name: '',
  receiver_phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  is_default: 0,
})
const saving = ref(false)
let addressId = 0

onLoad(async (q) => {
  addressId = Number(q.id || 0)
  if (!addressId) return
  try {
    const res = await listAddresses()
    const a = (res.data || []).find((x) => Number(x.id) === addressId)
    if (a) {
      form.receiver_name = a.receiver_name || ''
      form.receiver_phone = a.receiver_phone || ''
      form.province = a.province || ''
      form.city = a.city || ''
      form.district = a.district || ''
      form.detail = a.detail || ''
      form.is_default = a.is_default ? 1 : 0
    }
  } catch {
    /* ignore */
  }
})

async function save() {
  if (!form.receiver_name.trim()) {
    uni.showToast({ title: '请填写收货人', icon: 'none' })
    return
  }
  if (!form.receiver_phone.trim()) {
    uni.showToast({ title: '请填写手机号', icon: 'none' })
    return
  }
  if (!form.detail.trim()) {
    uni.showToast({ title: '请填写详细地址', icon: 'none' })
    return
  }
  saving.value = true
  try {
    const body = { ...form, is_default: form.is_default ? 1 : 0 }
    if (addressId) await updateAddress(addressId, body)
    else await createAddress(body)
    uni.showToast({ title: '已保存', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 400)
  } catch {
    /* handled */
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; }
.card {
  background: #fff; border-radius: 24rpx; padding: 8rpx 28rpx;
  box-shadow: 0 4rpx 24rpx rgba(200,168,118,.08);
}
.field {
  display: flex; align-items: center; gap: 16rpx;
  padding: 24rpx 0; border-bottom: 1rpx solid #f5f5f5;
}
.field:last-child { border-bottom: none; }
.label { width: 140rpx; color: #71717a; font-size: 26rpx; flex-shrink: 0; }
.input { flex: 1; font-size: 28rpx; }
.switch-row { justify-content: space-between; }
.box {
  width: 40rpx; height: 40rpx; border-radius: 50%; border: 2rpx solid #d4d4d8;
  display: flex; align-items: center; justify-content: center; font-size: 22rpx; color: #fff;
}
.box.on { background: #c8a876; border-color: #c8a876; }
.btn {
  margin-top: 40rpx; height: 88rpx; line-height: 88rpx; border-radius: 999rpx;
  background: linear-gradient(135deg, #bfa472, #d4b890); color: #fff; font-size: 30rpx;
}
</style>
