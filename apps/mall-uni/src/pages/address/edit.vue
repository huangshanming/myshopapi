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
      <picker mode="multiSelector" :range="pickerRange" range-key="name" :value="pickerIndex" @columnchange="onColumnChange" @change="onRegionChange">
        <view class="field">
          <text class="label">所在地区</text>
          <text class="input" :class="{ placeholder: !regionText }">{{ regionText || '请选择省市区' }}</text>
          <text class="arrow">›</text>
        </view>
      </picker>
      <view class="field">
        <text class="label">详细地址</text>
        <input v-model="form.detail" class="input" placeholder="街道、门牌号等" />
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
import { computed, reactive, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { createAddress, listAddresses, listRegions, updateAddress } from '../../api/index'

const form = reactive({
  receiver_name: '',
  receiver_phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  province_code: '',
  city_code: '',
  district_code: '',
  is_default: 0,
})
const saving = ref(false)
let addressId = 0

const provinces = ref([])
const cities = ref([])
const districts = ref([])
const pickerIndex = ref([0, 0, 0])

const pickerRange = computed(() => [provinces.value, cities.value, districts.value])

const regionText = computed(() => {
  if (!form.province) return ''
  return `${form.province} ${form.city} ${form.district}`
})

async function loadChildren(parentCode) {
  const res = await listRegions(parentCode)
  return res?.list || []
}

async function initPicker(pCode = '', cCode = '', dCode = '') {
  provinces.value = await loadChildren('')
  if (!provinces.value.length) return
  let pi = 0
  if (pCode) {
    const i = provinces.value.findIndex((x) => x.code === pCode)
    if (i >= 0) pi = i
  }
  cities.value = await loadChildren(provinces.value[pi].code)
  let ci = 0
  if (cCode && cities.value.length) {
    const i = cities.value.findIndex((x) => x.code === cCode)
    if (i >= 0) ci = i
  }
  if (!cities.value.length) {
    districts.value = []
    pickerIndex.value = [pi, 0, 0]
    return
  }
  districts.value = await loadChildren(cities.value[ci].code)
  let di = 0
  if (dCode && districts.value.length) {
    const i = districts.value.findIndex((x) => x.code === dCode)
    if (i >= 0) di = i
  }
  pickerIndex.value = [pi, ci, di]
}

async function onColumnChange(e) {
  const col = e.detail.column
  const row = e.detail.value
  const idx = [...pickerIndex.value]
  idx[col] = row
  if (col === 0) {
    cities.value = await loadChildren(provinces.value[row]?.code || '')
    idx[1] = 0
    idx[2] = 0
    districts.value = cities.value.length ? await loadChildren(cities.value[0].code) : []
  } else if (col === 1) {
    districts.value = await loadChildren(cities.value[row]?.code || '')
    idx[2] = 0
  }
  pickerIndex.value = idx
}

function onRegionChange(e) {
  const [pi, ci, di] = e.detail.value
  pickerIndex.value = [pi, ci, di]
  const p = provinces.value[pi]
  const c = cities.value[ci]
  const d = districts.value[di]
  form.province = p?.name || ''
  form.city = c?.name || ''
  form.district = d?.name || ''
  form.province_code = p?.code || ''
  form.city_code = c?.code || ''
  form.district_code = d?.code || ''
}

onLoad(async (q) => {
  addressId = Number(q.id || 0)
  await initPicker()
  if (!addressId) return
  try {
    const res = await listAddresses()
    const a = (res?.list || []).find((x) => Number(x.id) === addressId)
    if (!a) return
    form.receiver_name = a.receiver_name || ''
    form.receiver_phone = a.receiver_phone || ''
    form.province = a.province || ''
    form.city = a.city || ''
    form.district = a.district || ''
    form.detail = a.detail || ''
    form.province_code = a.province_code || ''
    form.city_code = a.city_code || ''
    form.district_code = a.district_code || ''
    form.is_default = a.is_default ? 1 : 0
    await initPicker(form.province_code, form.city_code, form.district_code)
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
  if (!form.province_code || !form.city_code || !form.district_code) {
    uni.showToast({ title: '请选择省市区', icon: 'none' })
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
.input { flex: 1; font-size: 28rpx; color: #18181b; }
.input.placeholder { color: #a1a1aa; }
.arrow { color: #c8a876; font-size: 32rpx; }
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
