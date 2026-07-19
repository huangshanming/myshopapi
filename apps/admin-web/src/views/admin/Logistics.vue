<template>
  <div>
    <div class="toolbar">
      <h2>物流管理</h2>
      <el-input v-model="q.name" placeholder="名称" clearable style="width: 140px" @keyup.enter="search" />
      <el-input v-model="q.code" placeholder="编码" clearable style="width: 120px" @keyup.enter="search" />
      <el-select v-model="q.status" clearable placeholder="状态" style="width: 110px" @change="search">
        <el-option label="启用" :value="1" />
        <el-option label="停用" :value="0" />
      </el-select>
      <el-button type="primary" @click="search">查询</el-button>
      <el-button v-permission="'business:logistics:add'" type="primary" @click="openForm()">新增</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="公司名称" min-width="160" />
      <el-table-column prop="code" label="编码" width="100" />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status === 1 ? '启用' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="更新时间" min-width="160" />
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'business:logistics:edit'" link type="primary" @click="openForm(row)">编辑</el-button>
          <el-button
            v-permission="'business:logistics:status'"
            link
            :type="row.status === 1 ? 'warning' : 'success'"
            @click="toggleStatus(row)"
          >{{ row.status === 1 ? '停用' : '启用' }}</el-button>
          <el-button v-permission="'business:logistics:delete'" link type="danger" @click="remove(row)">删除</el-button>
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

    <el-dialog v-model="visible" :title="editId ? '编辑物流公司' : '新增物流公司'" width="420px">
      <el-form label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="如：顺丰速运" maxlength="64" />
        </el-form-item>
        <el-form-item label="编码" required>
          <el-input v-model="form.code" placeholder="如：SF" maxlength="32" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listLogistics, createLogistics, updateLogistics, updateLogisticsStatus, deleteLogistics,
} from '../../api/order'

const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const q = reactive({ name: '', code: '', status: undefined })

const visible = ref(false)
const editId = ref(0)
const form = reactive({ name: '', code: '', sort: 0 })
const submitting = ref(false)

function search() {
  page.value = 1
  load()
}

async function load() {
  loading.value = true
  try {
    const res = await listLogistics({
      page: page.value,
      page_size: pageSize,
      name: q.name || undefined,
      code: q.code || undefined,
      status: q.status === undefined || q.status === null || q.status === '' ? undefined : q.status,
    })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
    list.value = []
  } finally {
    loading.value = false
  }
}

function openForm(row) {
  editId.value = row?.id || 0
  form.name = row?.name || ''
  form.code = row?.code || ''
  form.sort = row?.sort ?? 0
  visible.value = true
}

async function save() {
  if (!form.name.trim() || !form.code.trim()) {
    ElMessage.warning('请填写名称与编码')
    return
  }
  submitting.value = true
  try {
    const payload = { name: form.name.trim(), code: form.code.trim(), sort: form.sort }
    if (editId.value) await updateLogistics(editId.value, payload)
    else await createLogistics(payload)
    ElMessage.success('已保存')
    visible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  const next = row.status === 1 ? 0 : 1
  try {
    await updateLogisticsStatus(row.id, next)
    ElMessage.success(next === 1 ? '已启用' : '已停用')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确认删除「${row.name}」？`, '删除')
    await deleteLogistics(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.pager { margin-top: 16px; justify-content: flex-end; }
</style>
