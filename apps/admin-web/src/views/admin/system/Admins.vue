<template>
  <div>
    <div class="toolbar">
      <h2>管理员设置</h2>
      <el-button v-permission="'system:admin:add'" type="primary" @click="openCreate">新增管理员</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="mobile" label="手机号" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="role" label="身份" width="140" />
      <el-table-column prop="status" label="状态" width="80" />
      <el-table-column label="操作" width="260">
        <template #default="{ row }">
          <el-button v-permission="'system:admin:assign'" link type="primary" @click="openRoles(row)">分配角色</el-button>
          <el-button v-permission="'system:admin:reset'" link @click="openReset(row)">重置密码</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="createVisible" title="新增管理员" width="460px">
      <el-form label-width="80px">
        <el-form-item label="手机号"><el-input v-model="createForm.mobile" maxlength="11" /></el-form-item>
        <el-form-item label="昵称"><el-input v-model="createForm.nickname" /></el-form-item>
        <el-form-item label="密码"><el-input v-model="createForm.password" type="password" show-password /></el-form-item>
        <el-form-item label="角色">
          <el-select v-model="createForm.role_ids" multiple style="width: 100%">
            <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCreate">创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="roleVisible" title="分配角色" width="420px">
      <el-select v-model="selectedRoles" multiple style="width: 100%">
        <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
      </el-select>
      <template #footer>
        <el-button @click="roleVisible = false">取消</el-button>
        <el-button type="primary" @click="saveRoles">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="pwdVisible" title="重置密码" width="420px">
      <el-input v-model="newPassword" type="password" show-password placeholder="新密码" />
      <template #footer>
        <el-button @click="pwdVisible = false">取消</el-button>
        <el-button type="primary" @click="savePwd">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  assignAdminRoles, createAdmin, fetchAdminRoles, fetchAdmins, fetchRoles, resetAdminPassword,
} from '../../../api/system'
import { pickList } from '../../../utils/list'

const list = ref([])
const roles = ref([])
const loading = ref(false)
const createVisible = ref(false)
const roleVisible = ref(false)
const pwdVisible = ref(false)
const createForm = ref({ mobile: '', nickname: '', password: '', role_ids: [] })
const selectedRoles = ref([])
const currentId = ref(0)
const newPassword = ref('')

async function load() {
  loading.value = true
  try {
    const [adminsRes, rolesRes] = await Promise.all([fetchAdmins({ page: 1, page_size: 50 }), fetchRoles()])
    list.value = adminsRes?.list || []
    roles.value = pickList(rolesRes)
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openCreate() {
  createForm.value = { mobile: '', nickname: '', password: '123456', role_ids: [] }
  createVisible.value = true
}

async function saveCreate() {
  try {
    await createAdmin(createForm.value)
    ElMessage.success('已创建')
    createVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function openRoles(row) {
  currentId.value = row.id
  const res = await fetchAdminRoles(row.id)
  selectedRoles.value = pickList(res)
  roleVisible.value = true
}

async function saveRoles() {
  try {
    await assignAdminRoles(currentId.value, selectedRoles.value)
    ElMessage.success('已分配')
    roleVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

function openReset(row) {
  currentId.value = row.id
  newPassword.value = '123456'
  pwdVisible.value = true
}

async function savePwd() {
  try {
    await resetAdminPassword(currentId.value, newPassword.value)
    ElMessage.success('密码已重置')
    pwdVisible.value = false
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
