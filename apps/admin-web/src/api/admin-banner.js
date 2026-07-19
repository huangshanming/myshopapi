import http from './http'

const base = '/api/v1/admin/banners'

export function listBanners(params) {
  return http.get(base, { params })
}

export function getBanner(id) {
  return http.get(`${base}/${id}`)
}

export function createBanner(data) {
  return http.post(base, data)
}

export function updateBanner(id, data) {
  return http.put(`${base}/${id}`, data)
}

export function deleteBanner(id) {
  return http.delete(`${base}/${id}`)
}

export function uploadBannerImage(file) {
  const fd = new FormData()
  fd.append('file', file)
  return http.post(`${base}/upload`, fd, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export const BANNER_LINK_OPTIONS = [
  { value: 'none', label: '不跳转' },
  { value: 'product', label: '商品' },
  { value: 'article', label: '文章' },
]

export function bannerLinkLabel(t) {
  return BANNER_LINK_OPTIONS.find((o) => o.value === t)?.label || t || '-'
}
