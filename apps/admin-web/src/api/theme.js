import http from './http'

export const listThemeSlots = () => http.get('/api/v1/admin/theme-slots')
export const updateThemeSlot = (id, data) => http.put(`/api/v1/admin/theme-slots/${id}`, data)
export const listThemePackages = (params) => http.get('/api/v1/admin/theme-packages', { params })
export const createThemePackage = (data) => http.post('/api/v1/admin/theme-packages', data)
export const updateThemePackage = (id, data) => http.put(`/api/v1/admin/theme-packages/${id}`, data)
export const listThemeOrders = (params) => http.get('/api/v1/admin/theme-orders', { params })
export const grantThemeOrder = (data) => http.post('/api/v1/admin/theme-orders/grant', data)

export const merchantListThemeSlots = () => http.get('/api/v1/merchant/theme-slots')
export const merchantListThemePackages = (params) => http.get('/api/v1/merchant/theme-packages', { params })
export const merchantBuyTheme = (data) => http.post('/api/v1/merchant/theme-orders', data)
export const merchantListThemeOrders = (params) => http.get('/api/v1/merchant/theme-orders', { params })

export const THEME_LINK_OPTIONS = [
  { value: 'shop', label: '本店店铺' },
  { value: 'category', label: '商品分类' },
  { value: 'product', label: '指定商品' },
]

export const THEME_DEFAULT_LINK_OPTIONS = [
  { value: 'none', label: '不跳转' },
  { value: 'category', label: '商品分类' },
]

export function themeLinkLabel(t) {
  return [...THEME_LINK_OPTIONS, ...THEME_DEFAULT_LINK_OPTIONS].find((o) => o.value === t)?.label || t || '-'
}

export function themeStatusLabel(s) {
  return ({ active: '生效中', expired: '已过期', cancelled: '已取消', on: '上架', off: '下架' })[s] || s || '-'
}
