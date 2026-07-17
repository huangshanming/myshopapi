<template>
  <div>
    <h2>全站订单</h2>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="order_no" label="订单号" min-width="160" />
      <el-table-column prop="shop_id" label="店铺" width="80" />
      <el-table-column prop="user_id" label="用户" width="80" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column prop="total_amount" label="金额" width="100" />
      <el-table-column prop="created_at" label="创建时间" min-width="160" />
    </el-table>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import http from '../../api/http'
import { pickList } from '../../utils/list'

const list = ref([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await http.get('/api/v1/admin/orders', { params: { page: 1, page_size: 50 } })
    list.value = pickList(res.data)
  } catch (e) {
    ElMessage.error(e.message)
    list.value = []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
