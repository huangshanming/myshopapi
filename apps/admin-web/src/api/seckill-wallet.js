import http from './http'

export const getSeckillRule = () => http.get('/api/v1/admin/seckill/rule')
export const updateSeckillRule = (data) => http.put('/api/v1/admin/seckill/rule', data)
export const fetchSeckillSessions = (params) => http.get('/api/v1/admin/seckill/sessions', { params })
export const fetchSeckillEntries = (params) => http.get('/api/v1/admin/seckill/entries', { params })

export const getShopWallet = (shopId) => http.get(`/api/v1/admin/shops/${shopId}/wallet`)
export const adjustShopWallet = (shopId, data) => http.post(`/api/v1/admin/shops/${shopId}/wallet/adjust`, data)
export const fetchShopWalletLogs = (shopId, params) =>
  http.get(`/api/v1/admin/shops/${shopId}/wallet/logs`, { params })

export const getMerchantWallet = () => http.get('/api/v1/merchant/wallet')
export const fetchMerchantWalletLogs = (params) => http.get('/api/v1/merchant/wallet/logs', { params })
export const fetchMerchantSeckillSessions = () => http.get('/api/v1/merchant/seckill/sessions')
export const applySeckill = (data) => http.post('/api/v1/merchant/seckill/entries', data)
export const fetchMerchantSeckillEntries = (params) => http.get('/api/v1/merchant/seckill/entries', { params })
