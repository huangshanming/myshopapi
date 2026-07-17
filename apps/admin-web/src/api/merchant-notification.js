import http from './http'

const base = '/api/v1/merchant/notifications'

export function listNotifications(params) {
  return http.get(base, { params })
}

export function unreadNotificationCount() {
  return http.get(`${base}/unread-count`)
}

export function markNotificationRead(id) {
  return http.post(`${base}/${id}/read`)
}

export function markAllNotificationsRead() {
  return http.post(`${base}/read-all`)
}
