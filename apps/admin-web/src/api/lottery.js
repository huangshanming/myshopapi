import http from './http'

const base = '/api/v1/admin/lottery'

export function listLotteryActivities(params) {
  return http.get(`${base}/activities`, { params })
}

export function getLotteryActivity(id) {
  return http.get(`${base}/activities/${id}`)
}

export function createLotteryActivity(data) {
  return http.post(`${base}/activities`, data)
}

export function updateLotteryActivity(id, data) {
  return http.put(`${base}/activities/${id}`, data)
}

export function saveLotteryPrizes(id, prizes) {
  return http.put(`${base}/activities/${id}/prizes`, { prizes })
}

export function listLotteryRecords(params) {
  return http.get(`${base}/records`, { params })
}

export function listLotteryOrders(params) {
  return http.get(`${base}/orders`, { params })
}

export function shipLotteryOrder(id, data) {
  return http.post(`${base}/orders/${id}/ship`, data)
}

export function uploadLotteryImage(file) {
  const fd = new FormData()
  fd.append('file', file)
  return http.post(`${base}/upload`, fd, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
