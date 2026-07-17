<template>
  <div>
    <div class="toolbar">
      <h2>操作日志</h2>
      <el-button @click="load">刷新</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="product_id" label="商品ID" width="90" />
      <el-table-column prop="operator_id" label="操作人" width="90" />
      <el-table-column prop="action" label="动作" width="140" />
      <el-table-column prop="created_at" label="时间" width="180" />
      <el-table-column prop="after_json" label="变更" min-width="200" show-overflow-tooltip />
    </el-table>
    <div class="pager">
      <el-pagination
        v-model:current-page="page"
        layout="prev, pager, next"
        :total="total"
        @current-change="load"
      />
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { pickList } from '../../utils/list'
import { listOpLogs } from '../../api/merchant-product'

const list = ref([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)

async function load() {
  loading.value = true
  try {
    const res = await listOpLogs({ page: page.value, page_size: 20 })
    list.value = pickList(res.data)
    total.value = Number(res.data?.total || list.value.length)
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
