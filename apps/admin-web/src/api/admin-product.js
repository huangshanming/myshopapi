import http from './http'

export function listAdminProducts(params) {
  return http.get('/api/v1/admin/products', { params })
}

export function forceOffSaleProduct(id, remark) {
  return http.put(`/api/v1/admin/products/${id}/off_sale`, { remark })
}

export function deleteAdminProduct(id, remark) {
  return http.delete(`/api/v1/admin/products/${id}`, { data: { remark } })
}

/** 前台可见分类（商家选类目等） */
export function listProductCategories(params = {}) {
  return http.get('/api/v1/product_category/list', { params: { page: 1, page_size: 500, ...params } })
}

/** 平台管理：含隐藏分类 */
export function listAdminCategories() {
  return http.get('/api/v1/admin/categories')
}

export function createAdminCategory(data) {
  return http.post('/api/v1/admin/categories', data)
}

export function updateAdminCategory(id, data) {
  return http.put(`/api/v1/admin/categories/${id}`, data)
}

export function deleteAdminCategory(id) {
  return http.delete(`/api/v1/admin/categories/${id}`)
}
