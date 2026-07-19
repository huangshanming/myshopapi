<template>
  <div>
    <div class="toolbar">
      <h2>秒杀场次</h2>
      <el-button @click="load">刷新</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="rule_id" label="规则ID" width="80" />
      <el-table-column prop="start_at" label="开始" min-width="160" />
      <el-table-column prop="end_at" label="结束" min-width="160" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEntries(row)">报名明细</el-button>
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

    <el-drawer v-model="drawer" :title="`场次 #${sessionId} 报名`" size="640px">
      <el-table :data="entries" v-loading="entryLoading" size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="shop_id" label="店铺" width="70" />
        <el-table-column prop="product_name" label="商品" min-width="140" />
        <el-table-column prop="seckill_price" label="秒杀价" width="80" />
        <el-table-column prop="seckill_stock" label="库存" width="70" />
        <el-table-column prop="fee_amount" label="报名费" width="80" />
        <el-table-column prop="created_at" label="时间" width="150" />
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { fetchSeckillEntries, fetchSeckillSessions } from '../../api/seckill-wallet'

const list = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = 20
const total = ref(0)
const drawer = ref(false)
const sessionId = ref(0)
const entries = ref([])
const entryLoading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await fetchSeckillSessions({ page: page.value, page_size: pageSize })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function openEntries(row) {
  sessionId.value = row.id
  drawer.value = true
  entryLoading.value = true
  try {
    const res = await fetchSeckillEntries({ session_id: row.id, page: 1, page_size: 100 })
    entries.value = res?.list || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    entryLoading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; }
.toolbar h2 { margin-right: auto; }
.pager { margin-top: 12px; }
</style>
