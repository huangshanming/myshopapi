import http from './http'

/** @param {'admin'|'merchant'} scope */
function base(scope) {
  return scope === 'merchant' ? '/api/v1/merchant' : '/api/v1/admin'
}

export function listOrders(scope, params) {
  return http.get(`${base(scope)}/orders`, { params })
}

export function getOrder(scope, id) {
  return http.get(`${base(scope)}/orders/${id}`)
}

export function shipOrder(scope, id, data) {
  return http.put(`${base(scope)}/orders/${id}/ship`, data)
}

export function completeOrder(scope, id) {
  return http.put(`${base(scope)}/orders/${id}/complete`)
}

export function remarkOrder(scope, id, remark) {
  return http.put(`${base(scope)}/orders/${id}/remark`, { remark })
}

export function listAfterSales(scope, params) {
  return http.get(`${base(scope)}/after-sales`, { params })
}

export function handleAfterSale(scope, id, data) {
  return http.put(`${base(scope)}/after-sales/${id}/handle`, data)
}

export function listLogistics(params) {
  return http.get('/api/v1/admin/logistics', { params })
}

export function createLogistics(data) {
  return http.post('/api/v1/admin/logistics', data)
}

export function updateLogistics(id, data) {
  return http.put(`/api/v1/admin/logistics/${id}`, data)
}

export function updateLogisticsStatus(id, status) {
  return http.put(`/api/v1/admin/logistics/${id}/status`, { status })
}

export function deleteLogistics(id) {
  return http.delete(`/api/v1/admin/logistics/${id}`)
}

export function listLogisticsOptions(keyword) {
  return http.get('/api/v1/logistics/options', { params: { keyword: keyword || undefined } })
}

export function listReviews(scope, params) {
  return http.get(`${base(scope)}/reviews`, { params })
}

export function replyReview(id, data) {
  return http.put(`/api/v1/merchant/reviews/${id}/reply`, data)
}

export function deleteReview(scope, id) {
  return http.delete(`${base(scope)}/reviews/${id}`)
}

export const ORDER_STATUS_OPTIONS = [
  { value: 'pending', label: '待确认' },
  { value: 'confirmed', label: '待发货' },
  { value: 'shipped', label: '已发货' },
  { value: 'completed', label: '已完成' },
  { value: 'reviewed', label: '已评价' },
  { value: 'cancelled', label: '已取消' },
  { value: 'failed', label: '失败' },
]

export const AFTER_SALE_STATUS_OPTIONS = [
  { value: 'pending', label: '待处理' },
  { value: 'approved', label: '已同意' },
  { value: 'rejected', label: '已拒绝' },
  { value: 'refunded', label: '已退款' },
  { value: 'closed', label: '已关闭' },
]

const orderStatusMap = Object.fromEntries(ORDER_STATUS_OPTIONS.map((o) => [o.value, o.label]))
const afterSaleStatusMap = Object.fromEntries(AFTER_SALE_STATUS_OPTIONS.map((o) => [o.value, o.label]))

export function orderStatusLabel(s) {
  return orderStatusMap[s] || s
}

export function afterSaleStatusLabel(s) {
  return afterSaleStatusMap[s] || s
}

export function orderStatusType(s) {
  const map = {
    pending: 'info',
    confirmed: 'warning',
    shipped: 'primary',
    completed: 'success',
    reviewed: 'success',
    cancelled: 'info',
    failed: 'danger',
  }
  return map[s] || 'info'
}

/** 解析订单行 sku_snapshot 为可读文案 */
export function formatSkuSnapshot(snap) {
  if (!snap) return ''
  let obj = snap
  if (typeof snap === 'string') {
    try {
      obj = JSON.parse(snap)
    } catch {
      return snap
    }
  }
  if (!obj || typeof obj !== 'object') return ''
  const parts = []
  for (const [k, v] of Object.entries(obj)) {
    if (k === 'sku_id') {
      parts.push(`SKU#${v}`)
      continue
    }
    if (v == null || v === '') continue
    if (typeof v === 'object') continue
    parts.push(`${k}:${v}`)
  }
  return parts.join(' · ')
}

export function getProductDetail(id) {
  return http.get('/api/v1/products/detail', { params: { id } })
}
