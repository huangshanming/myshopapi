import http from './http'

export const fetchAuthMe = () => http.get('/api/v1/admin/auth/me')

export const fetchMenus = () => http.get('/api/v1/admin/menus')
export const createMenu = (data) => http.post('/api/v1/admin/menus', data)
export const updateMenu = (id, data) => http.put(`/api/v1/admin/menus/${id}`, data)
export const deleteMenu = (id) => http.delete(`/api/v1/admin/menus/${id}`)

export const fetchRoles = () => http.get('/api/v1/admin/roles')
export const createRole = (data) => http.post('/api/v1/admin/roles', data)
export const updateRole = (id, data) => http.put(`/api/v1/admin/roles/${id}`, data)
export const deleteRole = (id) => http.delete(`/api/v1/admin/roles/${id}`)
export const fetchRoleMenus = (id) => http.get(`/api/v1/admin/roles/${id}/menus`)
export const assignRoleMenus = (id, menu_ids) => http.put(`/api/v1/admin/roles/${id}/menus`, { menu_ids })

export const fetchUsers = (params) => http.get('/api/v1/admin/users', { params })
export const fetchUser = (id) => http.get(`/api/v1/admin/users/${id}`)
export const updateUser = (id, data) => http.put(`/api/v1/admin/users/${id}`, data)
export const setUserStatus = (id, status) => http.put(`/api/v1/admin/users/${id}/status`, { status })
export const resetUserPassword = (id, password) => http.put(`/api/v1/admin/users/${id}/password`, { password })
export const generateUserToken = (id) => http.post(`/api/v1/admin/users/${id}/token`)

export const fetchAdmins = (params) => http.get('/api/v1/admin/admins', { params })
export const createAdmin = (data) => http.post('/api/v1/admin/admins', data)
export const fetchAdminRoles = (id) => http.get(`/api/v1/admin/admins/${id}/roles`)
export const assignAdminRoles = (id, role_ids) => http.put(`/api/v1/admin/admins/${id}/roles`, { role_ids })
export const resetAdminPassword = (id, password) => http.put(`/api/v1/admin/admins/${id}/password`, { password })

export const fetchConfigs = () => http.get('/api/v1/admin/configs')
export const saveConfigs = (items) => http.put('/api/v1/admin/configs', { items })
