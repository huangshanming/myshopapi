/** 从分页/树/ID 列表接口中取出数组，避免 el-table 收到对象 */
export function pickList(data) {
  if (Array.isArray(data)) return data
  if (!data || typeof data !== 'object') return []
  if (Array.isArray(data.list)) return data.list
  if (Array.isArray(data.tree)) return data.tree
  if (Array.isArray(data.ids)) return data.ids
  if (data.data && Array.isArray(data.data.list)) return data.data.list
  if (Array.isArray(data.data)) return data.data
  return []
}
