<template>
  <div>
    <h2>入驻审核</h2>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="shop_name" label="店铺名" />
      <el-table-column prop="contact_name" label="联系人" />
      <el-table-column prop="contact_phone" label="手机" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" size="small" type="success" @click="approve(row)">通过</el-button>
          <el-button v-if="row.status === 'pending'" size="small" type="danger" @click="reject(row)">拒绝</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import http from '../../api/http'
import { pickList } from '../../utils/list'

const list = ref([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await http.get('/api/v1/admin/applications', { params: { page: 1, page_size: 50 } })
    list.value = pickList(res)
  } catch (e) {
    ElMessage.error(e.message)
    list.value = []
  } finally {
    loading.value = false
  }
}

async function approve(row) {
  await http.post(`/api/v1/admin/applications/${row.id}/approve`)
  ElMessage.success('已通过')
  load()
}

async function reject(row) {
  const { value } = await ElMessageBox.prompt('拒绝原因', '拒绝申请')
  await http.post(`/api/v1/admin/applications/${row.id}/reject`, { reason: value })
  ElMessage.success('已拒绝')
  load()
}

onMounted(load)
</script>
