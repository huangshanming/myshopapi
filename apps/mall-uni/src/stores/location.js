const STORAGE_KEY = 'mymall_city'

/** 常用市中心（热门城市快捷跳转） */
export const CITY_CENTERS = {
  长沙市: { province: '湖南省', latitude: 28.228209, longitude: 112.938814 },
  北京市: { province: '北京市', latitude: 39.904211, longitude: 116.407395 },
  上海市: { province: '上海市', latitude: 31.230416, longitude: 121.473701 },
  广州市: { province: '广东省', latitude: 23.12911, longitude: 113.264385 },
  深圳市: { province: '广东省', latitude: 22.543099, longitude: 114.057868 },
  杭州市: { province: '浙江省', latitude: 30.274085, longitude: 120.15507 },
  成都市: { province: '四川省', latitude: 30.572961, longitude: 104.066301 },
  武汉市: { province: '湖北省', latitude: 30.592849, longitude: 114.305539 },
}

const DEFAULT = {
  city: '长沙市',
  cityCode: '',
  province: '湖南省',
  provinceCode: '',
  address: '',
  latitude: 28.228209,
  longitude: 112.938814,
  userPicked: false,
}

function readStored() {
  try {
    const raw = uni.getStorageSync(STORAGE_KEY)
    if (!raw) return { ...DEFAULT }
    const data = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (!data || !data.city) return { ...DEFAULT }
    const center = CITY_CENTERS[data.city]
    return {
      city: data.city || DEFAULT.city,
      cityCode: data.cityCode || '',
      province: data.province || '',
      provinceCode: data.provinceCode || '',
      address: data.address || '',
      latitude: Number(data.latitude) || center?.latitude || DEFAULT.latitude,
      longitude: Number(data.longitude) || center?.longitude || DEFAULT.longitude,
      userPicked: !!data.userPicked,
    }
  } catch {
    return { ...DEFAULT }
  }
}

let state = readStored()

function persist() {
  uni.setStorageSync(STORAGE_KEY, JSON.stringify(state))
}

export function getLocation() {
  return { ...state }
}

export function getCity() {
  return state.city || DEFAULT.city
}

export function getCoords() {
  const lat = Number(state.latitude) || 0
  const lng = Number(state.longitude) || 0
  if (!lat && !lng) return null
  return { latitude: lat, longitude: lng }
}

export function hasCoords() {
  return !!getCoords()
}

export function setLocation(loc, { userPicked = true } = {}) {
  if (!loc || !loc.city) return state
  const center = CITY_CENTERS[loc.city]
  state = {
    city: loc.city,
    cityCode: loc.cityCode || '',
    province: loc.province || center?.province || '',
    provinceCode: loc.provinceCode || '',
    address: loc.address || '',
    latitude: Number(loc.latitude) || center?.latitude || state.latitude || DEFAULT.latitude,
    longitude: Number(loc.longitude) || center?.longitude || state.longitude || DEFAULT.longitude,
    userPicked: !!userPicked,
  }
  persist()
  return { ...state }
}

/** 仅在用户未手选时，用收货地址城市初始化（落到市中心坐标） */
export function applyAddressCityIfNeeded(addr) {
  if (!addr || !addr.city) return false
  if (state.userPicked) return false
  const center = CITY_CENTERS[addr.city]
  setLocation(
    {
      city: addr.city,
      cityCode: addr.city_code || addr.cityCode || '',
      province: addr.province || center?.province || '',
      provinceCode: addr.province_code || addr.provinceCode || '',
      address: [addr.province, addr.city, addr.district, addr.detail].filter(Boolean).join(''),
      latitude: center?.latitude,
      longitude: center?.longitude,
    },
    { userPicked: false },
  )
  return true
}
