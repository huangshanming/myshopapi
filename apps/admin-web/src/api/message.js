import http from './http'

export const sendAdminNotification = (data) => http.post('/api/v1/admin/notifications/send', data)
export const listAdminNotificationSends = (params) => http.get('/api/v1/admin/notifications/sends', { params })
export const listAdminNotificationRecipients = (id, params) =>
  http.get(`/api/v1/admin/notifications/sends/${id}/recipients`, { params })
