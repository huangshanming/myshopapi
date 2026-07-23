import { http, unwrapBody } from '../utils/request'
import { getToken, getUser } from '../stores/user'

export function listProducts(params) {
  return http.get('/api/v1/products/list', params)
}

export function listSalesRank(params = { page: 1, page_size: 20 }) {
  return http.get('/api/v1/products/sales-rank', params)
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

export function listHomeSlots(slotType) {
  return http.get('/api/v1/shops/home-slots', { slot_type: slotType })
}

export function listThemeTiles() {
  return http.get('/api/v1/home/theme-tiles')
}

export function getShop(id) {
  return http.get(`/api/v1/shops/${id}`)
}

export function listArticles(params = { page: 1, page_size: 10 }) {
  return http.get('/api/v1/articles/list', params)
}

export function listBanners() {
  return http.get('/api/v1/banners')
}

export function getArticle(id) {
  return http.get(`/api/v1/articles/${id}`)
}

export function listArticleComments(articleId, params = { page: 1, page_size: 20 }) {
  return http.get(`/api/v1/articles/${articleId}/comments`, params)
}

export function createArticleComment(articleId, data) {
  return http.post(`/api/v1/articles/${articleId}/comments`, data)
}

export function listCommentEmojis() {
  return http.get('/api/v1/comment-emojis')
}

export function likeArticle(id) {
  return http.post(`/api/v1/articles/${id}/like`)
}

export function unlikeArticle(id) {
  return http.delete(`/api/v1/articles/${id}/like`)
}

export function favoriteArticle(id) {
  return http.post(`/api/v1/articles/${id}/favorite`)
}

export function unfavoriteArticle(id) {
  return http.delete(`/api/v1/articles/${id}/favorite`)
}

export function listMyArticleFavorites(params = { page: 1, page_size: 10 }) {
  return http.get('/api/v1/user/article-favorites', params)
}

export function listMyArticleLikes(params = { page: 1, page_size: 10 }) {
  return http.get('/api/v1/user/article-likes', params)
}

export function listMyArticles(params = { page: 1, page_size: 10 }) {
  return http.get('/api/v1/user/articles', params)
}

export function getMyArticle(id) {
  return http.get(`/api/v1/user/articles/${id}`)
}

export function createMyArticle(data) {
  return http.post('/api/v1/user/articles', data)
}

export function updateMyArticle(id, data) {
  return http.put(`/api/v1/user/articles/${id}`, data)
}

export function deleteMyArticle(id) {
  return http.delete(`/api/v1/user/articles/${id}`)
}

export function uploadArticleImage(filePath) {
  return new Promise((resolve, reject) => {
    const token = getToken()
    const user = getUser()
    const header = {}
    if (token) header.Authorization = `Bearer ${token}`
    if (user?.id) header['X-User-Id'] = String(user.id)
    if (user?.role) header['X-User-Role'] = user.role
    else if (token) header['X-User-Role'] = 'user'
    uni.uploadFile({
      url: '/api/v1/user/article-uploads',
      filePath,
      name: 'file',
      header,
      success: (res) => {
        try {
          const body = JSON.parse(res.data)
          if (res.statusCode && res.statusCode >= 400) {
            reject(new Error(body?.msg || body?.message || `HTTP ${res.statusCode}`))
            return
          }
          resolve(unwrapBody(body))
        } catch (e) {
          reject(e)
        }
      },
      fail: reject,
    })
  })
}

export function getUserPoints() {
  return http.get('/api/v1/user/points')
}

export function listUserPointLogs(params) {
  return http.get('/api/v1/user/points/logs', params)
}

export function listPointsMallProducts(params) {
  return http.get('/api/v1/user/points-mall/product', params)
}

export function getPointsMallProduct(id) {
  return http.get(`/api/v1/user/points-mall/product/${id}`)
}

export function exchangePointsMall(data) {
  return http.post('/api/v1/user/points-mall/exchange', data)
}

export function listPointsMallOrders(params) {
  return http.get('/api/v1/user/points-mall/orders', params)
}

export function getPointsMallOrder(id) {
  return http.get(`/api/v1/user/points-mall/orders/${id}`)
}

export function listUserTasks() {
  return http.get('/api/v1/user/tasks')
}

export function checkinTask() {
  return http.post('/api/v1/user/tasks/checkin')
}

export function claimTask(code) {
  return http.post(`/api/v1/user/tasks/${code}/claim`)
}

export function reportTaskEvent(data) {
  return http.post('/api/v1/user/tasks/events', data, { silent: true })
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

export function createOrder(items, addressId, userCouponId = 0) {
  const body = { items, address_id: addressId }
  if (userCouponId) body.user_coupon_id = userCouponId
  return http.post('/api/v1/orders', body)
}

export function couponPreview(items, userCouponId = 0) {
  return http.post('/api/v1/orders/coupon-preview', {
    items,
    user_coupon_id: userCouponId || 0,
  })
}

export function listCouponCenter(params = {}) {
  return http.get('/api/v1/coupons/center', params)
}

export function listCouponPopup() {
  return http.get('/api/v1/coupons/popup')
}

export function claimCoupon(id, source = 'direct') {
  return http.post(`/api/v1/coupons/${id}/claim`, { source })
}

export function listMyCoupons(params = { page: 1, page_size: 20, status: 'unused' }) {
  return http.get('/api/v1/user/coupons', params)
}

export function listNotifications(params = { page: 1, page_size: 20 }) {
  return http.get('/api/v1/user/notifications', params)
}

export function getNotificationUnreadCount() {
  return http.get('/api/v1/user/notifications/unread-count')
}

export function markNotificationRead(id) {
  return http.post(`/api/v1/user/notifications/${id}/read`)
}

export function markAllNotificationsRead() {
  return http.post('/api/v1/user/notifications/read-all')
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

export function getOrderStatusCounts() {
  return http.get('/api/v1/orders/status-counts')
}

export function listMyAfterSales(params) {
  return http.get('/api/v1/orders/after-sales', params)
}

export function getOrder(id) {
  return http.get(`/api/v1/orders/${id}`)
}

export function cancelOrder(id) {
  return http.put(`/api/v1/orders/${id}/cancel`)
}

export function confirmReceive(id) {
  return http.put(`/api/v1/orders/${id}/confirm-receive`)
}

export function getReviewEligible(orderId) {
  return http.get(`/api/v1/orders/${orderId}/review-eligible`)
}

export function submitOrderReview(orderId, data) {
  return http.post(`/api/v1/orders/${orderId}/reviews`, data)
}

export function getOrderReview(orderId) {
  return http.get(`/api/v1/orders/${orderId}/review`)
}

export function listProductReviews(productId, params) {
  return http.get(`/api/v1/products/${productId}/reviews`, params)
}

export function uploadReviewImage(filePath) {
  return new Promise((resolve, reject) => {
    const token = getToken()
    const user = getUser()
    const header = {}
    if (token) header.Authorization = `Bearer ${token}`
    if (user?.id) header['X-User-Id'] = String(user.id)
    if (user?.role) header['X-User-Role'] = user.role
    else if (token) header['X-User-Role'] = 'user'
    uni.uploadFile({
      url: '/api/v1/user/review-uploads',
      filePath,
      name: 'file',
      header,
      success: (res) => {
        try {
          const body = JSON.parse(res.data)
          if (res.statusCode && res.statusCode >= 400) {
            reject(new Error(body?.msg || body?.message || `HTTP ${res.statusCode}`))
            return
          }
          resolve(body)
        } catch (e) {
          reject(e)
        }
      },
      fail: reject,
    })
  })
}

export function addFavorite(productId) {
  return http.post('/api/v1/user/favorites', { product_id: productId })
}

export function removeFavorite(productId) {
  return http.delete(`/api/v1/user/favorites/${productId}`)
}

export function batchRemoveFavorites(productIds) {
  return http.post('/api/v1/user/favorites/batch-remove', { product_ids: productIds })
}

export function listFavorites(params) {
  return http.get('/api/v1/user/favorites', params)
}

export function getFavoriteStatus(productId) {
  return http.get(`/api/v1/products/${productId}/favorite`)
}

export const ORDER_STATUS = {
  pending: '待付款',
  confirmed: '待发货',
  shipped: '待收货',
  completed: '待评价',
  reviewed: '已评价',
  cancelled: '已取消',
  failed: '已关闭',
}

export const AFTER_SALE_STATUS = {
  pending: '处理中',
  approved: '已同意',
  rejected: '已拒绝',
  refunded: '已退款',
  closed: '已关闭',
}
