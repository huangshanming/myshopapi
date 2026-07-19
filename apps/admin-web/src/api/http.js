import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '',
  timeout: 15000,
})

function pickErrMsg(err) {
  const data = err?.response?.data
  if (data && typeof data === 'object') {
    if (data.msg) return data.msg
    if (data.message) return data.message
  }
  if (err?.message) return err.message
  return '请求失败'
}

/** 兼容旧信封 {code,msg,data} 与新契约（body 即 DTO） */
export function unwrapBody(body) {
  if (
    body &&
    typeof body === 'object' &&
    typeof body.code === 'number' &&
    Object.prototype.hasOwnProperty.call(body, 'data')
  ) {
    if (body.code !== 200) {
      throw new Error(body.msg || '请求失败')
    }
    return body.data
  }
  return body
}

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  if (auth.userId) {
    config.headers['X-User-Id'] = String(auth.userId)
  }
  if (auth.role) {
    config.headers['X-User-Role'] = auth.role
  }
  if (auth.shopId) {
    config.headers['X-Shop-Id'] = String(auth.shopId)
  }
  return config
})

http.interceptors.response.use(
  (res) => {
    try {
      return unwrapBody(res.data)
    } catch (e) {
      return Promise.reject(e)
    }
  },
  (err) => Promise.reject(new Error(pickErrMsg(err)))
)

export default http
