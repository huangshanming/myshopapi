<template>
  <div>
    <div class="toolbar">
      <h2>菜单管理</h2>
      <el-button v-permission="'system:menu:add'" type="primary" @click="openEdit()">新增</el-button>
    </div>
    <el-table
      :data="tree"
      v-loading="loading"
      row-key="id"
      default-expand-all
      :tree-props="{ children: 'children' }"
      stripe
    >
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column prop="type" label="类型" width="90" />
      <el-table-column prop="path" label="路径" min-width="140" />
      <el-table-column prop="perms" label="权限码" min-width="160" />
      <el-table-column prop="sort" label="排序" width="70" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button v-permission="'system:menu:add'" link type="primary" @click="openEdit({ parent_id: row.id })">子级</el-button>
          <el-button v-permission="'system:menu:edit'" link @click="openEdit(row)">编辑</el-button>
          <el-button v-permission="'system:menu:delete'" link type="danger" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="form.id ? '编辑菜单' : '新增菜单'" width="520px">
      <el-form label-width="90px">
        <el-form-item label="上级ID"><el-input-number v-model="form.parent_id" :min="0" /></el-form-item>
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type">
            <el-option label="目录" value="dir" />
            <el-option label="菜单" value="menu" />
            <el-option label="按钮" value="button" />
          </el-select>
        </el-form-item>
        <el-form-item label="路径"><el-input v-model="form.path" /></el-form-item>
        <el-form-item label="组件"><el-input v-model="form.component" /></el-form-item>
        <el-form-item label="图标"><el-input v-model="form.icon" /></el-form-item>
        <el-form-item label="权限码"><el-input v-model="form.perms" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" /></el-form-item>
        <el-form-item label="可见"><el-switch v-model="form.visible" :active-value="1" :inactive-value="0" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createMenu, deleteMenu, fetchMenus, updateMenu } from '../../../api/system'

const loading = ref(false)
const tree = ref([])
const visible = ref(false)
const form = ref({})

async function load() {
  loading.value = true
  try {
    const res = await fetchMenus()
    tree.value = res.data || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openEdit(row = {}) {
  form.value = {
    id: row.id,
    parent_id: row.parent_id || 0,
    name: row.name || '',
    type: row.type || 'menu',
    path: row.path || '',
    component: row.component || '',
    icon: row.icon || '',
    perms: row.perms || '',
    sort: row.sort || 0,
    visible: row.visible ?? 1,
    status: row.status ?? 1,
  }
  visible.value = true
}

async function save() {
  try {
    if (form.value.id) {
      await updateMenu(form.value.id, form.value)
    } else {
      await createMenu(form.value)
    }
    ElMessage.success('已保存')
    visible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function onDelete(row) {
  await ElMessageBox.confirm(`删除菜单「${row.name}」？`, '确认')
  try {
    await deleteMenu(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
