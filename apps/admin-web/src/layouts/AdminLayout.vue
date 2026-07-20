<template>
  <el-container class="layout">
    <el-aside width="220px" class="aside">
      <div class="brand">mymall 平台</div>
      <el-menu :default-active="$route.path" router>
        <template v-for="item in menuNodes" :key="item.id">
          <el-sub-menu v-if="item.children?.length" :index="'dir-' + item.id">
            <template #title>{{ item.name }}</template>
            <el-menu-item
              v-for="child in item.children"
              :key="child.id"
              :index="child.path || ('/admin/' + child.id)"
            >{{ child.name }}</el-menu-item>
          </el-sub-menu>
          <el-menu-item v-else-if="item.path" :index="item.path">{{ item.name }}</el-menu-item>
        </template>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <span>{{ auth.user?.nickname || auth.user?.mobile }}</span>
        <el-button link type="danger" @click="onLogout">退出</el-button>
      </el-header>
      <el-main><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

/** 侧栏目录定义：按业务域拆分，避免塞进少数大类 */
const MENU_GROUPS = [
  {
    id: 140, name: '商家管理', sort: 10,
    items: [
      { id: 2, name: '入驻审核', path: '/admin/applications' },
      { id: 3, name: '店铺管理', path: '/admin/shops' },
    ],
  },
  {
    id: 141, name: '商品中心', sort: 15,
    items: [
      { id: 5, name: '全站商品', path: '/admin/products' },
      { id: 6, name: '商品分类', path: '/admin/products/categories' },
    ],
  },
  {
    id: 142, name: '交易中心', sort: 20,
    items: [
      { id: 4, name: '全站订单', path: '/admin/orders' },
      { id: 7, name: '售后管理', path: '/admin/after-sales' },
      { id: 8, name: '物流管理', path: '/admin/logistics' },
      { id: 110, name: '评价管理', path: '/admin/reviews' },
    ],
  },
  {
    id: 90, name: '内容社区', sort: 25,
    items: [
      { id: 91, name: '文章列表', path: '/admin/articles' },
      { id: 92, name: '分类管理', path: '/admin/articles/categories' },
      { id: 93, name: '评论管理', path: '/admin/articles/comments' },
      { id: 100, name: '评论表情', path: '/admin/articles/emojis' },
      { id: 94, name: '文章回收站', path: '/admin/articles/recycle' },
      { id: 95, name: '文章统计', path: '/admin/articles/stats' },
    ],
  },
  {
    id: 143, name: '首页运营', sort: 35,
    items: [
      { id: 115, name: '首页 Banner', path: '/admin/banners' },
      { id: 118, name: '主题集市', path: '/admin/themes' },
      { id: 112, name: '首页展位', path: '/admin/homepage' },
    ],
  },
  {
    id: 16, name: '营销玩法', sort: 40,
    items: [
      { id: 121, name: '优惠券', path: '/admin/coupons' },
      { id: 17, name: '秒杀规则', path: '/admin/seckill/rule' },
      { id: 18, name: '秒杀场次', path: '/admin/seckill/sessions' },
      { id: 125, name: '任务中心', path: '/admin/tasks' },
      { id: 127, name: '积分商品', path: '/admin/points-mall' },
      { id: 129, name: '积分订单', path: '/admin/points-orders' },
    ],
  },
  {
    id: 144, name: '用户触达', sort: 50,
    items: [
      { id: 124, name: '用户消息', path: '/admin/messages' },
    ],
  },
  {
    id: 10, name: '系统管理', sort: 90,
    items: [
      { id: 11, name: '菜单管理', path: '/admin/system/menus' },
      { id: 12, name: '角色管理', path: '/admin/system/roles' },
      { id: 13, name: '用户管理', path: '/admin/system/users' },
      { id: 14, name: '管理员设置', path: '/admin/system/admins' },
      { id: 15, name: '系统设置', path: '/admin/system/configs' },
    ],
  },
]

const HIDDEN_DIR_IDS = new Set([1]) // 业务管理(旧)

function flattenMenus(tree) {
  const byPath = new Map()
  const byId = new Map()
  const walk = (nodes) => {
    for (const n of nodes || []) {
      if (n.path) byPath.set(n.path, n)
      if (n.id != null) byId.set(n.id, n)
      if (n.children?.length) walk(n.children)
    }
  }
  walk(tree)
  return { byPath, byId }
}

function pickChild(raw, item) {
  const hit = raw.byPath.get(item.path) || raw.byId.get(item.id)
  if (hit && hit.visible === 0) return null
  return {
    id: hit?.id ?? item.id,
    name: hit?.name || item.name,
    path: hit?.path || item.path,
    type: 'menu',
    visible: 1,
  }
}

/** 按固定分组重组侧栏；DB 有权限的项优先保留名称，缺省用兜底补齐 */
function buildGroupedMenus(tree) {
  const raw = flattenMenus(tree)
  const knownPaths = new Set()
  const groups = []

  for (const g of MENU_GROUPS) {
    const children = []
    for (const item of g.items) {
      const child = pickChild(raw, item)
      if (!child) continue
      children.push(child)
      knownPaths.add(child.path)
    }
    if (!children.length) continue
    const dir = raw.byId.get(g.id)
    groups.push({
      id: g.id,
      name: dir?.name || g.name,
      type: 'dir',
      visible: 1,
      sort: g.sort,
      children,
    })
  }

  // 未归类的可见菜单，挂到末尾「其他」
  const orphans = []
  for (const [, node] of raw.byPath) {
    if (node.type === 'button' || node.visible === 0) continue
    if (!node.path || knownPaths.has(node.path)) continue
    orphans.push({
      id: node.id,
      name: node.name,
      path: node.path,
      type: 'menu',
      visible: 1,
    })
  }
  if (orphans.length) {
    groups.push({
      id: 9990,
      name: '其他',
      type: 'dir',
      visible: 1,
      sort: 80,
      children: orphans,
    })
  }

  return groups
    .filter((g) => !HIDDEN_DIR_IDS.has(g.id))
    .sort((a, b) => (a.sort || 0) - (b.sort || 0))
}

const FALLBACK_TREE = MENU_GROUPS.map((g) => ({
  id: g.id,
  name: g.name,
  type: 'dir',
  visible: 1,
  children: g.items.map((c) => ({ ...c, type: 'menu', visible: 1 })),
}))

const menuNodes = computed(() => {
  const tree = auth.menuTree || []
  if (tree.length) {
    const filtered = tree
      .filter((n) => n.type !== 'button' && n.visible !== 0 && !HIDDEN_DIR_IDS.has(n.id))
      .map((n) => ({
        ...n,
        children: (n.children || []).filter((c) => c.type === 'menu' && c.visible !== 0),
      }))
    return buildGroupedMenus(filtered)
  }
  return FALLBACK_TREE
})

onMounted(async () => {
  if (auth.isAdmin) {
    try {
      await auth.loadAuthMe()
    } catch (_) {
      /* ignore */
    }
  }
})

function onLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout { min-height: 100vh; }
.aside { background: #1f2d3d; color: #fff; }
.brand { padding: 20px 16px; font-weight: 700; color: #fff; }
.aside :deep(.el-menu) { border-right: none; background: transparent; }
.aside :deep(.el-menu-item),
.aside :deep(.el-sub-menu__title) { color: #c0c4cc; }
.aside :deep(.el-menu-item.is-active) { background: #304156; color: #fff; }
.header { display: flex; justify-content: flex-end; align-items: center; gap: 12px; border-bottom: 1px solid #eee; }
</style>
