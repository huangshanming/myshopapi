import http from './http'

export const listHomepagePackages = (params) => http.get('/api/v1/admin/homepage-packages', { params })
export const createHomepagePackage = (data) => http.post('/api/v1/admin/homepage-packages', data)
export const updateHomepagePackage = (id, data) => http.put(`/api/v1/admin/homepage-packages/${id}`, data)
export const listHomepageSettings = () => http.get('/api/v1/admin/homepage-settings')
export const updateHomepageSettings = (items) => http.put('/api/v1/admin/homepage-settings', { items })
export const listHomepageOrders = (params) => http.get('/api/v1/admin/homepage-orders', { params })
export const grantHomepageOrder = (data) => http.post('/api/v1/admin/homepage-orders/grant', data)

export const merchantListHomepagePackages = (params) => http.get('/api/v1/merchant/homepage-packages', { params })
export const merchantBuyHomepage = (data) => http.post('/api/v1/merchant/homepage-orders', data)
export const merchantListHomepageOrders = (params) => http.get('/api/v1/merchant/homepage-orders', { params })

export const SLOT_TYPE_OPTIONS = [
  { value: 'brand_shop', label: '品牌商户' },
  { value: 'quality_shop', label: '优质商户' },
  { value: 'article', label: '种草文章' },
]

export const SLOT_STATUS_LABEL = {
  active: '生效中',
  expired: '已过期',
  cancelled: '已取消',
  on: '上架',
  off: '下架',
}

export const SLOT_PAY_SOURCE_LABEL = {
  admin: '超管代开通',
  wallet: '钱包购买',
}

export function slotTypeLabel(t) {
  return SLOT_TYPE_OPTIONS.find((o) => o.value === t)?.label || t || '-'
}

export function slotStatusLabel(s) {
  return SLOT_STATUS_LABEL[s] || s || '-'
}

export function slotPaySourceLabel(s) {
  return SLOT_PAY_SOURCE_LABEL[s] || s || '-'
}

/** 目标展示：店铺展位用店名，文章展位用标题或文章ID */
export function slotTargetLabel(row) {
  if (!row) return '-'
  if (row.slot_type === 'article') {
    return row.target_name || (row.target_id ? `文章 #${row.target_id}` : '-')
  }
  return row.shop_name || (row.shop_id ? `店铺 #${row.shop_id}` : '-')
}
