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

export function listShops(params = { page: 1, page_size: 20 }) {
  return http.get('/api/v1/shops/list', params)
}

export function getShop(id) {
  return http.get(`/api/v1/shops/${id}`)
}

export function listRegions(parentCode = '') {
  return http.get('/api/v1/regions', { parent_code: parentCode || '' })
}

export function getRegionTree() {
  return http.get('/api/v1/regions/tree')
}

export function getSeckillCurrent() {
  return http.get('/api/v1/seckill/current')
}

export function listSeckill(params = { page: 1, page_size: 10 }) {
  return http.get('/api/v1/seckill/list', params)
}

export function getSeckillEntry(id) {
  return http.get(`/api/v1/seckill/entries/${id}`)
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

export function getUserWallet() {
  return http.get('/api/v1/user/wallet')
}

export function listUserWalletLogs(params) {
  return http.get('/api/v1/user/wallet/logs', params)
}

export function createOrder(items, addressId) {
  return http.post('/api/v1/orders', { items, address_id: addressId })
}

export function listAddresses() {
  return http.get('/api/v1/user/addresses')
}

export function createAddress(data) {
  return http.post('/api/v1/user/addresses', data)
}

export function updateAddress(id, data) {
  return http.put(`/api/v1/user/addresses/${id}`, data)
}

export function deleteAddress(id) {
  return http.delete(`/api/v1/user/addresses/${id}`)
}

export function setDefaultAddress(id) {
  return http.put(`/api/v1/user/addresses/${id}/default`)
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
