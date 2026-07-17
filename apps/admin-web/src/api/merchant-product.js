import http from './http'

const base = '/api/v1/merchant'

export function listProducts(params) {
  return http.get(`${base}/products`, { params })
}

export function getProduct(id) {
  return http.get(`${base}/products/${id}`)
}

export function createProduct(data) {
  return http.post(`${base}/products`, data)
}

export function updateProduct(id, data) {
  return http.put(`${base}/products/${id}`, data)
}

export function copyProduct(id) {
  return http.post(`${base}/products/${id}/copy`)
}

export function setProductStatus(id, status) {
  return http.put(`${base}/products/${id}/status`, { status })
}

export function batchProducts(data) {
  return http.post(`${base}/products/batch`, data)
}

export function getBatchJob(id) {
  return http.get(`${base}/products/jobs/${id}`)
}

export function restoreRecycle(product_ids) {
  return http.post(`${base}/products/recycle/restore`, { product_ids })
}

export function deleteRecycle(product_ids) {
  return http.delete(`${base}/products/recycle`, { data: { product_ids } })
}

export function adjustSkuStock(id, data) {
  return http.put(`${base}/skus/${id}/stock`, data)
}

export function stockWarnings(params) {
  return http.get(`${base}/stocks/warnings`, { params })
}

export function uploadImage(file) {
  const fd = new FormData()
  fd.append('file', file)
  return http.post(`${base}/uploads/images`, fd, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function scheduleProduct(id, data) {
  return http.post(`${base}/products/${id}/schedules`, data)
}

export function listOpLogs(params) {
  return http.get(`${base}/products/op-logs`, { params })
}

export function exportProducts() {
  return http.get(`${base}/products/export`, { responseType: 'blob' })
}

export function importProducts(file) {
  const fd = new FormData()
  fd.append('file', file)
  return http.post(`${base}/products/import`, fd, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function listTags() {
  return http.get(`${base}/tags`)
}

export function saveTag(data, id) {
  return id ? http.put(`${base}/tags/${id}`, data) : http.post(`${base}/tags`, data)
}

export function deleteTag(id) {
  return http.delete(`${base}/tags/${id}`)
}

export function listAttrTemplates() {
  return http.get(`${base}/attr-templates`)
}

export function saveAttrTemplate(data, id) {
  return id
    ? http.put(`${base}/attr-templates/${id}`, data)
    : http.post(`${base}/attr-templates`, data)
}

export function fetchMerchantMe() {
  return http.get(`${base}/auth/me`)
}

export function listShopRoles() {
  return http.get(`${base}/shop/roles`)
}

export function listShopMenus() {
  return http.get(`${base}/shop/menus`)
}

export function getRoleMenus(id) {
  return http.get(`${base}/shop/roles/${id}/menus`)
}

export function saveShopRole(data, id) {
  return id ? http.put(`${base}/shop/roles/${id}`, data) : http.post(`${base}/shop/roles`, data)
}

export function listShopStaff() {
  return http.get(`${base}/shop/staff`)
}

export function bindShopStaff(data) {
  return http.post(`${base}/shop/staff`, data)
}
