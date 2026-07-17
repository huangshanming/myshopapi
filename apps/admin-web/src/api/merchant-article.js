import http from './http'

const base = '/api/v1/merchant'

export function listMyArticles(params) {
  return http.get(`${base}/articles`, { params })
}

export function getMyArticle(id) {
  return http.get(`${base}/articles/${id}`)
}

export function createMyArticle(data) {
  return http.post(`${base}/articles`, data)
}

export function updateMyArticle(id, data) {
  return http.put(`${base}/articles/${id}`, data)
}

export function deleteMyArticle(id) {
  return http.delete(`${base}/articles/${id}`)
}

export function listMyArticleCategories() {
  return http.get(`${base}/article-categories`)
}

export function listMyArticleComments(params) {
  return http.get(`${base}/article-comments`, { params })
}

export function patchMyArticleComment(id, status) {
  return http.patch(`${base}/article-comments/${id}`, { status })
}

export function deleteMyArticleComment(id) {
  return http.delete(`${base}/article-comments/${id}`)
}

export function uploadMyArticleImage(file) {
  const fd = new FormData()
  fd.append('file', file)
  return http.post(`${base}/article-uploads`, fd)
}
