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

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  // 直连微服务时网关不会注入头，本地代理用 JWT 字段手动补
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
    const body = res.data
    if (body && typeof body.code === 'number' && body.code !== 200) {
      return Promise.reject(new Error(body.msg || '请求失败'))
    }
    return body
  },
  (err) => Promise.reject(new Error(pickErrMsg(err)))
)

export default http
