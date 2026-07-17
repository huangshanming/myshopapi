<template>
  <div>
    <div class="toolbar">
      <h2>文章分类</h2>
      <el-button type="primary" @click="openForm(0)">新增根分类</el-button>
    </div>
    <el-table :data="flat" row-key="id" default-expand-all stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="200" />
      <el-table-column prop="parent_id" label="父级" width="80" />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column prop="status" label="状态" width="80" />
      <el-table-column label="操作" width="240">
        <template #default="{ row }">
          <el-button link type="primary" @click="openForm(row.id, row)">编辑</el-button>
          <el-button link @click="openForm(0, { parent_id: row.id })">加子类</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="editId ? '编辑分类' : '新增分类'" width="420px">
      <el-form label-width="80px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="父级ID"><el-input-number v-model="form.parent_id" :min="0" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" /></el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listArticleCategories, createArticleCategory, updateArticleCategory, deleteArticleCategory,
} from '../../../api/admin-article'

const tree = ref([])
const flat = ref([])
const visible = ref(false)
const editId = ref(0)
const form = reactive({ name: '', parent_id: 0, sort: 0, status: 1 })

function flatten(nodes, out = []) {
  for (const n of nodes || []) {
    out.push(n)
    if (n.children?.length) flatten(n.children, out)
  }
  return out
}

async function load() {
  const res = await listArticleCategories()
  tree.value = res.data || []
  flat.value = flatten(tree.value)
}

function openForm(id, row) {
  editId.value = id
  form.name = row?.name || ''
  form.parent_id = row?.parent_id || 0
  form.sort = row?.sort || 0
  form.status = row?.status ?? 1
  visible.value = true
}

async function save() {
  if (!form.name) return ElMessage.warning('请填写名称')
  if (editId.value) await updateArticleCategory(editId.value, { ...form })
  else await createArticleCategory({ ...form })
  visible.value = false
  ElMessage.success('已保存')
  load()
}

async function remove(row) {
  await ElMessageBox.confirm('确认删除？')
  await deleteArticleCategory(row.id)
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
</style>
