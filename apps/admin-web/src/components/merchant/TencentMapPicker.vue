<template>
  <div class="map-wrap">
    <div v-if="error" class="map-err">{{ error }}</div>
    <div ref="box" class="map-box" />
    <div class="map-tip">点击地图选点，可拖动标记微调</div>
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { getMapConfig } from '../../api/shop-profile'

const props = defineProps({
  latitude: { type: Number, default: 0 },
  longitude: { type: Number, default: 0 },
})
const emit = defineEmits(['pick'])

const box = ref(null)
const error = ref('')
let map = null
let marker = null
let TMap = null

function loadScript(key) {
  return new Promise((resolve, reject) => {
    if (window.TMap) {
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

async function init() {
  error.value = ''
  try {
    const cfg = await getMapConfig()
    const key = cfg?.key || import.meta.env.VITE_TENCENT_MAP_KEY || ''
    if (!key) {
      error.value = '未配置腾讯地图 Key（merchant-service TencentMap.Key 或 VITE_TENCENT_MAP_KEY）'
      return
    }
    TMap = await loadScript(key)
    const lat = props.latitude || 28.2282
    const lng = props.longitude || 112.9388
    map = new TMap.Map(box.value, {
      center: new TMap.LatLng(lat, lng),
      zoom: 14,
    })
    setMarker(lat, lng)
    map.on('click', (evt) => {
      const { lat: la, lng: ln } = evt.latLng
      setMarker(la, ln)
      emit('pick', { latitude: la, longitude: ln })
    })
  } catch (e) {
    error.value = e.message || '地图初始化失败'
  }
}

watch(() => [props.latitude, props.longitude], ([la, ln]) => {
  if (la && ln) setMarker(la, ln)
})

onMounted(init)
onBeforeUnmount(() => {
  if (map) {
    map.destroy()
    map = null
  }
})
</script>

<style scoped>
.map-wrap { width: 100%; }
.map-box { width: 100%; height: 360px; border-radius: 8px; overflow: hidden; background: #f1f5f9; }
.map-tip { margin-top: 8px; font-size: 12px; color: #94a3b8; }
.map-err { color: #dc2626; margin-bottom: 8px; font-size: 13px; }
</style>
