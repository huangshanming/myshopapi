<template>
  <div>
    <div class="toolbar">
      <h2>用户管理</h2>
      <el-input v-model="mobile" placeholder="手机号" clearable style="width: 200px" @clear="load" @keyup.enter="load" />
      <el-button @click="load">查询</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="mobile" label="手机号" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="role" label="角色" width="140" />
      <el-table-column prop="status" label="状态" width="80" />
      <el-table-column label="操作" width="260">
        <template #default="{ row }">
          <el-button v-permission="'system:user:edit'" link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button v-permission="'system:user:reset'" link @click="openReset(row)">重置密码</el-button>
          <el-button
            v-permission="'system:user:status'"
            size="small"
            :type="row.status === 1 ? 'warning' : 'success'"
            @click="toggle(row)"
          >{{ row.status === 1 ? '禁用' : '启用' }}</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pager"
      layout="prev, pager, next, total"
      :total="total"
      v-model:current-page="page"
      :page-size="pageSize"
      @current-change="load"
    />

    <el-dialog v-model="editVisible" title="编辑用户" width="460px">
      <el-form label-width="80px">
        <el-form-item label="手机号"><el-input v-model="editForm.mobile" maxlength="11" /></el-form-item>
        <el-form-item label="昵称"><el-input v-model="editForm.nickname" /></el-form-item>
        <el-form-item label="头像URL"><el-input v-model="editForm.avatar" /></el-form-item>
        <el-form-item label="性别">
          <el-select v-model="editForm.gender" style="width: 100%">
            <el-option label="未知" :value="0" />
            <el-option label="男" :value="1" />
            <el-option label="女" :value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEdit">保存</el-button>
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
import { fetchUsers, resetUserPassword, setUserStatus, updateUser } from '../../../api/system'

const list = ref([])
const loading = ref(false)
const mobile = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)
const editVisible = ref(false)
const pwdVisible = ref(false)
const editForm = ref({ mobile: '', nickname: '', avatar: '', gender: 0 })
const currentId = ref(0)
const newPassword = ref('123456')

async function load() {
  loading.value = true
  try {
    const res = await fetchUsers({ page: page.value, page_size: pageSize, mobile: mobile.value || undefined })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openEdit(row) {
  currentId.value = row.id
  editForm.value = {
    mobile: row.mobile || '',
    nickname: row.nickname || '',
    avatar: row.avatar || '',
    gender: row.gender ?? 0,
  }
  editVisible.value = true
}

async function saveEdit() {
  try {
    await updateUser(currentId.value, editForm.value)
    ElMessage.success('已更新')
    editVisible.value = false
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
    await resetUserPassword(currentId.value, newPassword.value)
    ElMessage.success('已重置')
    pwdVisible.value = false
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function toggle(row) {
  const status = row.status === 1 ? 0 : 1
  try {
    await setUserStatus(row.id, status)
    ElMessage.success('已更新')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin-right: auto; }
.pager { margin-top: 12px; }
</style>
