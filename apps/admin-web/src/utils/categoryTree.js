/** 将扁平分类列表转为 el-tree-select 树 */
export function buildCategoryTree(list) {
  const nodes = (list || []).map((c) => ({
    id: c.id,
    parent_id: c.parent_id ?? c.ParentId ?? 0,
    name: c.name,
    children: [],
  }))
  const byId = new Map(nodes.map((n) => [n.id, n]))
  const roots = []
  for (const n of nodes) {
    const parent = n.parent_id ? byId.get(n.parent_id) : null
    if (parent) parent.children.push(n)
    else roots.push(n)
  }
  const prune = (arr) =>
    arr.map((n) => {
      const children = n.children?.length ? prune(n.children) : undefined
      return children?.length ? { id: n.id, name: n.name, children } : { id: n.id, name: n.name }
    })
  return prune(roots)
}
