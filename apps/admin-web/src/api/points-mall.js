import http from './http'

const base = '/api/v1/admin/points-products'

export function listPointsProducts(params) {
  return http.get(base, { params })
}

export function getPointsProduct(id) {
  return http.get(`${base}/${id}`)
}

export function createPointsProduct(data) {
  return http.post(base, data)
}

export function updatePointsProduct(id, data) {
  return http.put(`${base}/${id}`, data)
}

export function setPointsProductStatus(id, status) {
  return http.put(`${base}/${id}/status`, { status })
}

export function deletePointsProduct(id) {
  return http.delete(`${base}/${id}`)
}

export function uploadPointsProductImage(file) {
  const fd = new FormData()
  fd.append('file', file)
  return http.post(`${base}/upload`, fd, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
