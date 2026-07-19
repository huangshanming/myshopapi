import { getToken, getUser } from '../stores/user'

const BASE = ''

function pickErr(data, fallback) {
  if (data && typeof data === 'object') {
    if (data.msg) return data.msg
    if (data.message) return data.message
  }
  return fallback || '请求失败'
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

export function request(options = {}) {
  const { url, method = 'GET', data, header = {}, silent } = options
  const token = getToken()
  const user = getUser()
  const headers = {
    'Content-Type': 'application/json',
    ...header,
  }
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }
  if (user?.id) {
    headers['X-User-Id'] = String(user.id)
  }
  if (user?.role) {
    headers['X-User-Role'] = user.role
  } else if (token) {
    headers['X-User-Role'] = 'user'
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE + url,
      method,
      data,
      header: headers,
      success(res) {
        const body = res.data
        if (res.statusCode >= 400) {
          const msg = pickErr(body, `HTTP ${res.statusCode}`)
          if (!silent) uni.showToast({ title: msg, icon: 'none' })
          reject(new Error(msg))
          return
        }
        try {
          resolve(unwrapBody(body))
        } catch (e) {
          if (!silent) uni.showToast({ title: e.message || '请求失败', icon: 'none' })
          reject(e)
        }
      },
      fail(err) {
        const msg = err?.errMsg || '网络错误'
        if (!silent) uni.showToast({ title: msg, icon: 'none' })
        reject(new Error(msg))
      },
    })
  })
}

export const http = {
  get: (url, data, opts) => request({ url, method: 'GET', data, ...opts }),
  post: (url, data, opts) => request({ url, method: 'POST', data, ...opts }),
  put: (url, data, opts) => request({ url, method: 'PUT', data, ...opts }),
  delete: (url, data, opts) => request({ url, method: 'DELETE', data, ...opts }),
}
