<template>
  <div>
    <div class="toolbar">
      <h2>商品分类</h2>
      <el-button type="primary" @click="openForm(0)">新增根分类</el-button>
      <el-button @click="load">刷新</el-button>
    </div>
    <el-table
      :data="tree"
      v-loading="loading"
      row-key="id"
      default-expand-all
      :tree-props="{ children: 'children' }"
      stripe
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="200" />
      <el-table-column prop="parent_id" label="父级ID" width="90" />
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column prop="level" label="层级" width="70" />
      <el-table-column label="显示" width="90">
        <template #default="{ row }">
          <el-tag :type="row.is_show ? 'success' : 'info'" size="small">{{ row.is_show ? '显示' : '隐藏' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="product_count" label="商品数" width="90" />
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openForm(row.id, row)">编辑</el-button>
          <el-button link @click="openForm(0, { parent_id: row.id, level: (row.level || 1) + 1 })">加子类</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="editId ? '编辑分类' : '新增分类'" width="480px" destroy-on-close>
      <el-form label-width="90px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" maxlength="100" show-word-limit />
        </el-form-item>
        <el-form-item label="父级分类">
          <el-tree-select
            v-model="form.parent_id"
            :data="parentOptions"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            check-strictly
            clearable
            filterable
            placeholder="空=根分类"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="层级">
          <el-input-number v-model="form.level" :min="1" :max="5" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" placeholder="可选 URL / 图标名" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="前台显示">
          <el-switch v-model="form.is_show" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listAdminCategories, createAdminCategory, updateAdminCategory, deleteAdminCategory,
} from '../../../api/admin-product'
import { buildCategoryTree } from '../../../utils/categoryTree'

const loading = ref(false)
const saving = ref(false)
const flat = ref([])
const tree = ref([])
const visible = ref(false)
const editId = ref(0)
const form = reactive({
  name: '',
  parent_id: undefined,
  sort_order: 0,
  level: 1,
  icon: '',
  description: '',
  is_show: true,
})

/** 表格用：保留全部分类字段的树 */
function buildFullTree(list) {
  const nodes = (list || []).map((c) => ({ ...c, children: [] }))
  const byId = new Map(nodes.map((n) => [n.id, n]))
  const roots = []
  for (const n of nodes) {
    const pid = n.parent_id || 0
    const parent = pid ? byId.get(pid) : null
    if (parent) parent.children.push(n)
    else roots.push(n)
  }
  const prune = (arr) =>
    arr.map((n) => {
      const { children, ...rest } = n
      if (children?.length) return { ...rest, children: prune(children) }
      return { ...rest }
    })
  return prune(roots)
}

const parentOptions = computed(() => {
  const exclude = editId.value
  return buildCategoryTree(
    flat.value.filter((c) => c.id !== exclude && !isDescendantOf(c.id, exclude)),
  )
})

function isDescendantOf(id, ancestorId) {
  if (!ancestorId) return false
  let cur = flat.value.find((c) => c.id === id)
  while (cur?.parent_id) {
    if (cur.parent_id === ancestorId) return true
    cur = flat.value.find((c) => c.id === cur.parent_id)
  }
  return false
}

async function load() {
  loading.value = true
  try {
    const res = await listAdminCategories()
    const list = Array.isArray(res) ? res : (res?.list || [])
    flat.value = list
    tree.value = buildFullTree(list)
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openForm(id, row) {
  editId.value = id
  form.name = row?.name || ''
  form.parent_id = row?.parent_id || undefined
  form.sort_order = row?.sort_order ?? 0
  form.level = row?.level || 1
  form.icon = row?.icon || ''
  form.description = row?.description || ''
  form.is_show = row?.is_show !== false && row?.is_show !== 0
  visible.value = true
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写名称')
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.name.trim(),
      parent_id: form.parent_id || 0,
      sort_order: form.sort_order,
      level: form.level || 1,
      icon: form.icon,
      description: form.description,
      is_show: form.is_show,
    }
    if (editId.value) await updateAdminCategory(editId.value, payload)
    else await createAdminCategory(payload)
    ElMessage.success('已保存')
    visible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

async function remove(row) {
  await ElMessageBox.confirm(`确认删除分类「${row.name}」？`, '提示', { type: 'warning' })
  try {
    await deleteAdminCategory(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
</style>
