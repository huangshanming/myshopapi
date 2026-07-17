/** 从分页接口 data 中取出数组，避免 el-table 收到对象报 rows is not iterable */
export function pickList(data) {
  if (Array.isArray(data)) return data
  if (!data || typeof data !== 'object') return []
  if (Array.isArray(data.list)) return data.list
  if (data.data && Array.isArray(data.data.list)) return data.data.list
  if (Array.isArray(data.data)) return data.data
  return []
}
