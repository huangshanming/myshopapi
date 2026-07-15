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
      <el-table-column label="操作" width="140">
        <template #default="{ row }">
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
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { fetchUsers, setUserStatus } from '../../../api/system'

const list = ref([])
const loading = ref(false)
const mobile = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)

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
