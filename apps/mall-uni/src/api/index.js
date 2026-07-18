import { http } from '../utils/request'

export function listProducts(params) {
  return http.get('/api/v1/products/list', params)
}

export function getProductDetail(id) {
  return http.get('/api/v1/products/detail', { id })
}

export function listCategories(params = { page: 1, page_size: 20 }) {
  return http.get('/api/v1/product_category/list', params)
}

export function login(data) {
  return http.post('/api/v1/user/login', data)
}

export function register(data) {
  return http.post('/api/v1/user/register', data)
}

export function getProfile() {
  return http.get('/api/v1/user/profile')
}

export function createOrder(items) {
  return http.post('/api/v1/orders', { items })
}

export function listOrders(params) {
  return http.get('/api/v1/orders', params)
}

export function getOrder(id) {
  return http.get(`/api/v1/orders/${id}`)
}

export function cancelOrder(id) {
  return http.put(`/api/v1/orders/${id}/cancel`)
}

export const ORDER_STATUS = {
  pending: '待确认',
  confirmed: '待发货',
  shipped: '已发货',
  completed: '已完成',
  cancelled: '已取消',
  failed: '失败',
}
