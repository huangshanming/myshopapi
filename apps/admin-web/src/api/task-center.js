import http from './http'

export const listAdminTasks = () => http.get('/api/v1/admin/tasks')
export const updateAdminTask = (id, data) => http.put(`/api/v1/admin/tasks/${id}`, data)
