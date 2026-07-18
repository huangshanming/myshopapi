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

const EXTRA_MENUS = [
  { id: 5, name: '全站商品', path: '/admin/products', type: 'menu', visible: 1 },
  { id: 6, name: '商品分类', path: '/admin/products/categories', type: 'menu', visible: 1 },
  { id: 7, name: '售后管理', path: '/admin/after-sales', type: 'menu', visible: 1 },
  { id: 8, name: '物流管理', path: '/admin/logistics', type: 'menu', visible: 1 },
]

const EXTRA_ARTICLE_DIR = {
  id: 90,
  name: '文章管理',
  type: 'dir',
  visible: 1,
  children: [
    { id: 91, name: '文章列表', path: '/admin/articles', type: 'menu', visible: 1 },
    { id: 92, name: '分类管理', path: '/admin/articles/categories', type: 'menu', visible: 1 },
    { id: 93, name: '评论管理', path: '/admin/articles/comments', type: 'menu', visible: 1 },
    { id: 94, name: '文章回收站', path: '/admin/articles/recycle', type: 'menu', visible: 1 },
    { id: 95, name: '文章统计', path: '/admin/articles/stats', type: 'menu', visible: 1 },
  ],
}

function ensureBizMenus(tree) {
  const cloned = tree.map((n) => ({
    ...n,
    children: (n.children || []).map((c) => ({ ...c })),
  }))
  let biz = cloned.find((n) => n.id === 1 || n.name === '业务管理')
  if (!biz) {
    biz = { id: 1, name: '业务管理', type: 'dir', children: [] }
    cloned.unshift(biz)
  }
  biz.children = biz.children || []
  for (const m of EXTRA_MENUS) {
    if (!biz.children.some((c) => c.path === m.path || c.id === m.id)) {
      biz.children.push({ ...m })
    }
  }
  return ensureArticleMenus(cloned)
}

function ensureArticleMenus(tree) {
  const cloned = tree.map((n) => ({
    ...n,
    children: (n.children || []).map((c) => ({ ...c })),
  }))
  let art = cloned.find((n) => n.id === 90 || n.name === '文章管理')
  if (!art) {
    cloned.push({
      ...EXTRA_ARTICLE_DIR,
      children: EXTRA_ARTICLE_DIR.children.map((c) => ({ ...c })),
    })
    return cloned
  }
  art.children = art.children || []
  for (const m of EXTRA_ARTICLE_DIR.children) {
    if (!art.children.some((c) => c.path === m.path || c.id === m.id)) {
      art.children.push({ ...m })
    }
  }
  return cloned
}

const menuNodes = computed(() => {
  const tree = auth.menuTree || []
  if (tree.length) {
    const filtered = tree
      .filter((n) => n.type !== 'button' && n.visible !== 0)
      .map((n) => ({
        ...n,
        children: (n.children || []).filter((c) => c.type === 'menu' && c.visible !== 0),
      }))
    return ensureBizMenus(filtered)
  }
  return [
    {
      id: 1,
      name: '业务管理',
      children: [
        { id: 2, name: '入驻审核', path: '/admin/applications' },
        { id: 3, name: '店铺管理', path: '/admin/shops' },
        { id: 4, name: '全站订单', path: '/admin/orders' },
        { id: 5, name: '全站商品', path: '/admin/products' },
        { id: 6, name: '商品分类', path: '/admin/products/categories' },
      ],
    },
    {
      id: 90,
      name: '文章管理',
      children: [
        { id: 91, name: '文章列表', path: '/admin/articles' },
        { id: 92, name: '分类管理', path: '/admin/articles/categories' },
        { id: 93, name: '评论管理', path: '/admin/articles/comments' },
        { id: 94, name: '文章回收站', path: '/admin/articles/recycle' },
        { id: 95, name: '文章统计', path: '/admin/articles/stats' },
      ],
    },
  ]
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
