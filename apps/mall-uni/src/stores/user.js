const TOKEN_KEY = 'mymall_token'
const USER_KEY = 'mymall_user'

export function getToken() {
  return uni.getStorageSync(TOKEN_KEY) || ''
}

export function getUser() {
  try {
    const raw = uni.getStorageSync(USER_KEY)
    return raw ? (typeof raw === 'string' ? JSON.parse(raw) : raw) : null
  } catch {
    return null
  }
}

export function setAuth(token, user) {
  uni.setStorageSync(TOKEN_KEY, token || '')
  uni.setStorageSync(USER_KEY, user ? JSON.stringify(user) : '')
}

export function clearAuth() {
  uni.removeStorageSync(TOKEN_KEY)
  uni.removeStorageSync(USER_KEY)
}

export function isLoggedIn() {
  return !!getToken()
}
