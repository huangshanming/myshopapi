import http from './http'

export function listMyShops() {
  return http.get('/api/v1/merchant/shops')
}

export function updateMyShop(id, data) {
  return http.put(`/api/v1/merchant/shops/${id}`, data)
}

export function getMapConfig() {
  return http.get('/api/v1/map/config')
}

export function reverseGeocode(lat, lng) {
  return http.get('/api/v1/map/geocoder', { params: { lat, lng } })
}
