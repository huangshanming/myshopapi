<template>
  <div>
    <div class="toolbar">
      <h2>角色管理</h2>
      <el-button v-permission="'system:role:add'" type="primary" @click="openEdit()">新增角色</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="code" label="编码" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="status" label="状态" width="80" />
      <el-table-column prop="remark" label="备注" />
      <el-table-column label="操作" width="260">
        <template #default="{ row }">
          <el-button v-permission="'system:role:edit'" link @click="openEdit(row)">编辑</el-button>
          <el-button v-permission="'system:role:assign'" link type="primary" @click="openMenus(row)">分配菜单</el-button>
          <el-button v-if="row.code !== 'super_admin'" v-permission="'system:role:delete'" link type="danger" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="editVisible" :title="form.id ? '编辑角色' : '新增角色'" width="480px">
      <el-form label-width="80px">
        <el-form-item label="编码"><el-input v-model="form.code" :disabled="!!form.id" /></el-form-item>
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="menuVisible" title="分配菜单" width="480px">
      <el-tree
        ref="treeRef"
        :data="menuTree"
        node-key="id"
        show-checkbox
        default-expand-all
        :props="{ label: 'name', children: 'children' }"
      />
      <template #footer>
        <el-button @click="menuVisible = false">取消</el-button>
        <el-button type="primary" @click="saveMenus">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  assignRoleMenus, createRole, deleteRole, fetchMenus, fetchRoleMenus,
  fetchRoles, updateRole,
} from '../../../api/system'

const loading = ref(false)
const list = ref([])
const editVisible = ref(false)
const menuVisible = ref(false)
const form = ref({})
const menuTree = ref([])
const currentRoleId = ref(0)
const treeRef = ref()

async function load() {
  loading.value = true
  try {
    const res = await fetchRoles()
    list.value = res || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openEdit(row = {}) {
  form.value = {
    id: row.id,
    code: row.code || '',
    name: row.name || '',
    remark: row.remark || '',
    status: row.status ?? 1,
  }
  editVisible.value = true
}

async function save() {
  try {
    if (form.value.id) await updateRole(form.value.id, form.value)
    else await createRole(form.value)
    ElMessage.success('已保存')
    editVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function onDelete(row) {
  await ElMessageBox.confirm(`删除角色「${row.name}」？`, '确认')
  try {
    await deleteRole(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function openMenus(row) {
  currentRoleId.value = row.id
  const [menusRes, idsRes] = await Promise.all([fetchMenus(), fetchRoleMenus(row.id)])
  menuTree.value = menusRes || []
  menuVisible.value = true
  await nextTick()
  treeRef.value?.setCheckedKeys(idsRes || [], false)
}

async function saveMenus() {
  const keys = treeRef.value.getCheckedKeys(false)
  const half = treeRef.value.getHalfCheckedKeys()
  try {
    await assignRoleMenus(currentRoleId.value, [...keys, ...half])
    ElMessage.success('已分配')
    menuVisible.value = false
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
