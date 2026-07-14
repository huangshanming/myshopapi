<template>
  <div>
    <h2>店铺管理</h2>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column prop="owner_user_id" label="店主UID" width="100" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button v-if="row.status !== 'disabled'" size="small" type="warning" @click="disable(row)">停用</el-button>
          <el-button v-else size="small" type="success" @click="enable(row)">启用</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import http from '../../api/http'

const list = ref([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await http.get('/api/v1/admin/shops', { params: { page: 1, page_size: 50 } })
    list.value = res.data?.list || res.data || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function disable(row) {
  await http.put(`/api/v1/admin/shops/${row.id}/disable`)
  ElMessage.success('已停用')
  load()
}

async function enable(row) {
  await http.put(`/api/v1/admin/shops/${row.id}/enable`)
  ElMessage.success('已启用')
  load()
}

onMounted(load)
</script>
