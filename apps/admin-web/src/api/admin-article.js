import http from './http'

const base = '/api/v1/admin'

export function listArticles(params) {
  return http.get(`${base}/articles`, { params })
}

export function getArticle(id) {
  return http.get(`${base}/articles/${id}`)
}

export function createArticle(data) {
  return http.post(`${base}/articles`, data)
}

export function updateArticle(id, data) {
  return http.put(`${base}/articles/${id}`, data)
}

export function deleteArticle(id, remark) {
  return http.delete(`${base}/articles/${id}`, { data: { remark: remark || '' } })
}

export function auditArticle(id, data) {
  return http.post(`${base}/articles/${id}/audit`, data)
}

export function batchAuditArticles(data) {
  return http.post(`${base}/articles/batch-audit`, data)
}

export function topArticle(id, is_top) {
  return http.post(`${base}/articles/${id}/top`, { is_top })
}

export function offlineArticle(id, remark) {
  return http.post(`${base}/articles/${id}/offline`, { remark: remark || '' })
}

export function listArticleRecycle(params) {
  return http.get(`${base}/articles/recycle`, { params: { ...params, recycle: 1 } })
}

export function restoreArticle(id) {
  return http.post(`${base}/articles/recycle/restore`, { id })
}

export function permanentDeleteArticle(id) {
  return http.delete(`${base}/articles/recycle`, { data: { id } })
}

export function articleStats() {
  return http.get(`${base}/articles/stats`)
}

export function listArticleCategories() {
  return http.get(`${base}/article-categories`)
}

export function createArticleCategory(data) {
  return http.post(`${base}/article-categories`, data)
}

export function updateArticleCategory(id, data) {
  return http.put(`${base}/article-categories/${id}`, data)
}

export function deleteArticleCategory(id) {
  return http.delete(`${base}/article-categories/${id}`)
}

export function listArticleComments(params) {
  return http.get(`${base}/article-comments`, { params })
}

export function patchArticleComment(id, status) {
  return http.patch(`${base}/article-comments/${id}`, { status })
}

export function deleteArticleComment(id) {
  return http.delete(`${base}/article-comments/${id}`)
}

export function uploadArticleImage(file, shopId = 0) {
  const fd = new FormData()
  fd.append('file', file)
  const q = shopId ? `?shop_id=${shopId}` : ''
  return http.post(`${base}/article-uploads${q}`, fd)
}
