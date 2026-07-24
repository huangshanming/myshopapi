<template>
  <view class="page">
    <view class="current">
      <view>
        <text class="label">当前位置</text>
        <text class="value">{{ current.city || '未选择' }}</text>
      </view>
      <text class="addr">{{ current.address || '在地图上点选精确位置' }}</text>
    </view>

    <view class="map-card">
      <view v-if="mapError" class="map-err">{{ mapError }}</view>
      <!-- #ifdef H5 -->
      <view id="loc-map" class="map-box" />
      <!-- #endif -->
      <!-- #ifndef H5 -->
      <view class="map-fallback">
        <text>当前端请用热门城市选择；H5 支持地图选点</text>
      </view>
      <!-- #endif -->
      <text class="map-tip">点击地图选点，可再点热门城市跳转市中心</text>
    </view>

    <view v-if="loggedIn" class="addr-card" @tap="useAddressCity">
      <view class="addr-main">
        <text class="addr-title">使用收货地址城市</text>
        <text class="addr-sub">{{ addressHint }}</text>
      </view>
      <text class="arrow">›</text>
    </view>

    <view class="section">
      <text class="sec-title">热门城市</text>
      <view class="hot-grid">
        <text
          v-for="c in hotCities"
          :key="c.city"
          class="hot-item"
          :class="{ on: current.city === c.city }"
          @tap="jumpHotCity(c)"
        >{{ c.city }}</text>
      </view>
    </view>

    <view class="section">
      <text class="sec-title">选择省市</text>
      <picker mode="multiSelector" :range="pickerRange" range-key="name" :value="pickerIndex" @columnchange="onColumnChange" @change="onRegionChange">
        <view class="picker-row">
          <text class="picker-text">{{ pickerText || '请选择省 / 市' }}</text>
          <text class="arrow">›</text>
        </view>
      </picker>
    </view>

    <button class="btn" type="primary" @tap.stop="confirm">确认定位</button>
  </view>
</template>

<script setup>
import { computed, getCurrentInstance, nextTick, onBeforeUnmount, ref } from 'vue'
import { onReady, onShow } from '@dcloudio/uni-app'
import { getMapConfig, listAddresses, listRegions, reverseGeocode } from '../../api/index'
import { CITY_CENTERS, getLocation, setLocation } from '../../stores/location'
import { isLoggedIn } from '../../stores/user'

const current = ref(getLocation())
const draft = ref({ ...getLocation() })
const loggedIn = ref(isLoggedIn())
const defaultAddr = ref(null)
const mapError = ref('')

const hotCities = Object.keys(CITY_CENTERS).map((city) => ({
  city,
  province: CITY_CENTERS[city].province,
  latitude: CITY_CENTERS[city].latitude,
  longitude: CITY_CENTERS[city].longitude,
}))

const provinces = ref([])
const cities = ref([])
const pickerIndex = ref([0, 0])
const pickerRange = computed(() => [provinces.value, cities.value])
const pickerText = computed(() => {
  const p = provinces.value[pickerIndex.value[0]]
  const c = cities.value[pickerIndex.value[1]]
  if (!p || !c) return ''
  return `${p.name} / ${c.name}`
})
const addressHint = computed(() => {
  const a = defaultAddr.value
  if (!a) return loggedIn.value ? '暂无默认收货地址' : '登录后可用'
  return `${a.province || ''} ${a.city || ''}`.trim() || '暂无默认收货地址'
})

let map = null
let marker = null
let TMap = null

async function loadChildren(parentCode) {
  const res = await listRegions(parentCode)
  return res?.list || []
}

async function initPicker() {
  provinces.value = await loadChildren('')
  if (!provinces.value.length) return
  let pi = 0
  if (draft.value.province) {
    const i = provinces.value.findIndex((x) => x.name === draft.value.province)
    if (i >= 0) pi = i
  }
  cities.value = await loadChildren(provinces.value[pi].code)
  let ci = 0
  if (draft.value.city && cities.value.length) {
    const i = cities.value.findIndex((x) => x.name === draft.value.city)
    if (i >= 0) ci = i
  }
  pickerIndex.value = [pi, ci]
}

async function onColumnChange(e) {
  const col = e.detail.column
  const row = e.detail.value
  const idx = [...pickerIndex.value]
  idx[col] = row
  if (col === 0) {
    cities.value = await loadChildren(provinces.value[row]?.code || '')
    idx[1] = 0
  }
  pickerIndex.value = idx
}

function onRegionChange(e) {
  const [pi, ci] = e.detail.value
  pickerIndex.value = [pi, ci]
  const p = provinces.value[pi]
  const c = cities.value[ci]
  if (!c?.name) return
  const center = CITY_CENTERS[c.name]
  applyDraft({
    city: c.name,
    cityCode: c.code || '',
    province: p?.name || '',
    provinceCode: p?.code || '',
    latitude: center?.latitude,
    longitude: center?.longitude,
  })
  moveMap(draft.value.latitude, draft.value.longitude)
}

function applyDraft(partial) {
  draft.value = { ...draft.value, ...partial }
  current.value = { ...draft.value }
}

async function applyGeo(lat, lng) {
  applyDraft({ latitude: lat, longitude: lng })
  try {
    const geo = await reverseGeocode(lat, lng)
    if (geo) {
      applyDraft({
        province: geo.province || draft.value.province,
        city: geo.city || draft.value.city,
        address: geo.address || draft.value.address,
        latitude: lat,
        longitude: lng,
      })
    }
  } catch (_) { /* ignore */ }
}

function jumpHotCity(c) {
  applyDraft({
    city: c.city,
    province: c.province,
    latitude: c.latitude,
    longitude: c.longitude,
    address: '',
  })
  moveMap(c.latitude, c.longitude)
  applyGeo(c.latitude, c.longitude)
}

function useAddressCity() {
  const a = defaultAddr.value
  if (!a?.city) {
    uni.showToast({ title: '暂无可用收货地址', icon: 'none' })
    return
  }
  const center = CITY_CENTERS[a.city]
  applyDraft({
    city: a.city,
    cityCode: a.city_code || '',
    province: a.province || center?.province || '',
    provinceCode: a.province_code || '',
    address: [a.province, a.city, a.district, a.detail].filter(Boolean).join(''),
    latitude: center?.latitude || draft.value.latitude,
    longitude: center?.longitude || draft.value.longitude,
  })
  moveMap(draft.value.latitude, draft.value.longitude)
}

function confirm() {
  if (!draft.value.city) {
    uni.showToast({ title: '请先选择位置', icon: 'none' })
    return
  }
  setLocation(draft.value, { userPicked: true })
  uni.showToast({ title: `已定位到${draft.value.city}`, icon: 'none' })
  setTimeout(() => {
    uni.navigateBack({
      fail: () => {
        uni.switchTab({ url: '/pages/index/index' })
      },
    })
  }, 200)
}

function loadScript(key) {
  return new Promise((resolve, reject) => {
    // #ifdef H5
    if (typeof window !== 'undefined' && window.TMap) {
      resolve(window.TMap)
      return
    }
    const id = 'tencent-map-gljs'
    if (document.getElementById(id)) {
      const t = setInterval(() => {
        if (window.TMap) {
          clearInterval(t)
          resolve(window.TMap)
        }
      }, 50)
      return
    }
    const s = document.createElement('script')
    s.id = id
    s.src = `https://map.qq.com/api/gljs?v=1.exp&key=${encodeURIComponent(key)}`
    s.onload = () => resolve(window.TMap)
    s.onerror = () => reject(new Error('腾讯地图脚本加载失败'))
    document.head.appendChild(s)
    // #endif
    // #ifndef H5
    reject(new Error('当前端不支持地图选点'))
    // #endif
  })
}

function setMarker(lat, lng) {
  if (!map || !TMap) return
  const pos = new TMap.LatLng(lat, lng)
  if (!marker) {
    marker = new TMap.MultiMarker({
      map,
      geometries: [{ id: 'p', position: pos }],
    })
  } else {
    marker.updateGeometries([{ id: 'p', position: pos }])
  }
  map.setCenter(pos)
}

function moveMap(lat, lng) {
  if (lat && lng) setMarker(lat, lng)
}

async function initMap() {
  // #ifdef H5
  mapError.value = ''
  try {
    const cfg = await getMapConfig()
    const key = cfg?.key || ''
    if (!key) {
      mapError.value = '未配置腾讯地图 Key'
      return
    }
    TMap = await loadScript(key)
    await nextTick()
    const el = document.getElementById('loc-map')
    if (!el) return
    const lat = draft.value.latitude || 28.228209
    const lng = draft.value.longitude || 112.938814
    map = new TMap.Map(el, {
      center: new TMap.LatLng(lat, lng),
      zoom: 14,
    })
    setMarker(lat, lng)
    map.on('click', (evt) => {
      const { lat: la, lng: ln } = evt.latLng
      setMarker(la, ln)
      applyGeo(la, ln)
    })
  } catch (e) {
    mapError.value = e.message || '地图初始化失败'
  }
  // #endif
}

async function loadDefaultAddress() {
  loggedIn.value = isLoggedIn()
  defaultAddr.value = null
  if (!loggedIn.value) return
  try {
    const res = await listAddresses({ silent: true })
    const list = res?.list || res || []
    const rows = Array.isArray(list) ? list : []
    defaultAddr.value = rows.find((a) => a.is_default) || rows[0] || null
  } catch {
    defaultAddr.value = null
  }
}

onShow(() => {
  current.value = getLocation()
  draft.value = { ...getLocation() }
  loadDefaultAddress()
  initPicker()
})

onReady(() => {
  initMap()
})

onBeforeUnmount(() => {
  if (map) {
    map.destroy()
    map = null
  }
})

// silence unused in non-H5
void getCurrentInstance
</script>

<style scoped>
.page { padding: 24rpx 32rpx 48rpx; }
.current {
  background: #fff; border-radius: 20rpx; padding: 28rpx 32rpx; margin-bottom: 24rpx;
}
.label { display: block; font-size: 24rpx; color: #71717a; }
.value { display: block; margin-top: 8rpx; font-size: 34rpx; font-weight: 600; color: #18181b; }
.addr { display: block; margin-top: 10rpx; font-size: 24rpx; color: #a1a1aa; }
.map-card {
  background: #fff; border-radius: 20rpx; padding: 16rpx; margin-bottom: 24rpx;
}
.map-box { width: 100%; height: 420rpx; border-radius: 12rpx; overflow: hidden; background: #f4f4f5; }
.map-fallback {
  height: 200rpx; display: flex; align-items: center; justify-content: center;
  color: #a1a1aa; font-size: 24rpx;
}
.map-tip { display: block; margin-top: 12rpx; font-size: 22rpx; color: #a1a1aa; padding: 0 8rpx; }
.map-err { color: #dc2626; font-size: 24rpx; margin-bottom: 12rpx; }
.addr-card {
  background: #fff; border-radius: 20rpx; padding: 28rpx 32rpx;
  display: flex; align-items: center; gap: 16rpx; margin-bottom: 24rpx;
}
.addr-main { flex: 1; min-width: 0; }
.addr-title { display: block; font-size: 28rpx; font-weight: 600; color: #18181b; }
.addr-sub { display: block; margin-top: 8rpx; font-size: 24rpx; color: #a1a1aa; }
.arrow { color: #a1a1aa; font-size: 36rpx; }
.section { margin-top: 8rpx; }
.sec-title {
  display: block; font-size: 26rpx; color: #71717a; margin: 16rpx 8rpx 20rpx;
}
.hot-grid { display: flex; flex-wrap: wrap; gap: 16rpx; }
.hot-item {
  width: calc(25% - 12rpx); box-sizing: border-box;
  text-align: center; font-size: 26rpx; color: #3f3f46;
  background: #fff; border-radius: 12rpx; padding: 20rpx 8rpx;
}
.hot-item.on {
  color: #c8a876; background: rgba(200,168,118,.12); font-weight: 600;
}
.picker-row {
  background: #fff; border-radius: 20rpx; padding: 28rpx 32rpx;
  display: flex; align-items: center; justify-content: space-between;
}
.picker-text { font-size: 28rpx; color: #18181b; }
.btn {
  margin-top: 40rpx; background: #c8a876; color: #fff; border: none;
  border-radius: 999rpx; font-size: 30rpx;
}
</style>
