import http from './http'

export const COUPON_TYPES = [
  { value: 'full_reduce', label: '满减券' },
  { value: 'no_threshold', label: '无门槛券' },
  { value: 'category', label: '品类券' },
  { value: 'product', label: '商品专享券' },
  { value: 'discount', label: '折扣券' },
]

export const couponTypeLabel = (v) => COUPON_TYPES.find((x) => x.value === v)?.label || v

export const displayStatusLabel = (v) => ({
  draft: '草稿', not_started: '未开始', active: '发放中', sold_out: '已领完', expired: '已过期', off: '已下架',
}[v] || v)

export const listAdminCoupons = (params) => http.get('/api/v1/admin/coupons', { params })
export const createAdminCoupon = (data) => http.post('/api/v1/admin/coupons', data)
export const updateAdminCoupon = (id, data) => http.put(`/api/v1/admin/coupons/${id}`, data)
export const offAdminCoupon = (id) => http.put(`/api/v1/admin/coupons/${id}/off`)
export const copyAdminCoupon = (id) => http.post(`/api/v1/admin/coupons/${id}/copy`)
export const grantAdminCoupon = (data) => http.post('/api/v1/admin/coupons/grant', data)
export const adminCouponClaims = (id, params) => http.get(`/api/v1/admin/coupons/${id}/claims`, { params })
export const adminCouponRedeems = (id, params) => http.get(`/api/v1/admin/coupons/${id}/redeems`, { params })
export const adminCouponStats = (id) => http.get(`/api/v1/admin/coupons/${id}/stats`)

export const listMerchantCoupons = (params) => http.get('/api/v1/merchant/coupons', { params })
export const createMerchantCoupon = (data) => http.post('/api/v1/merchant/coupons', data)
export const updateMerchantCoupon = (id, data) => http.put(`/api/v1/merchant/coupons/${id}`, data)
export const offMerchantCoupon = (id) => http.put(`/api/v1/merchant/coupons/${id}/off`)
export const copyMerchantCoupon = (id) => http.post(`/api/v1/merchant/coupons/${id}/copy`)
export const grantMerchantCoupon = (data) => http.post('/api/v1/merchant/coupons/grant', data)
export const merchantCouponClaims = (id, params) => http.get(`/api/v1/merchant/coupons/${id}/claims`, { params })
export const merchantCouponRedeems = (id, params) => http.get(`/api/v1/merchant/coupons/${id}/redeems`, { params })
export const merchantCouponStats = (id) => http.get(`/api/v1/merchant/coupons/${id}/stats`)
