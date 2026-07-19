import http from './http'

export const fetchShops = (params) => http.get('/api/v1/admin/shops', { params })
export const fetchShop = (id) => http.get(`/api/v1/admin/shops/${id}`)
export const createShop = (data) => http.post('/api/v1/admin/shops', data)
export const updateShop = (id, data) => http.put(`/api/v1/admin/shops/${id}`, data)
export const disableShop = (id, reason) => http.put(`/api/v1/admin/shops/${id}/disable`, { reason })
export const enableShop = (id) => http.put(`/api/v1/admin/shops/${id}/enable`)
export const resetShopOwnerPassword = (id, password) =>
  http.put(`/api/v1/admin/shops/${id}/owner-password`, { password })

export function uploadShopImage(file, shopId = 0) {
  const fd = new FormData()
  fd.append('file', file)
  const q = shopId ? `?shop_id=${shopId}` : ''
  return http.post(`/api/v1/admin/shop-uploads${q}`, fd, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
