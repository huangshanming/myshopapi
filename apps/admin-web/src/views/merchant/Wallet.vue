<template>
  <div>
    <h2>店铺钱包</h2>
    <div v-loading="loading" class="nums">
      <div class="card"><span>可用余额</span><b>¥{{ wallet.balance ?? 0 }}</b></div>
      <div class="card"><span>冻结余额</span><b>¥{{ wallet.frozen_balance ?? 0 }}</b></div>
      <div class="card"><span>保证金</span><b>¥{{ wallet.deposit ?? 0 }}</b></div>
    </div>
    <h3>变动记录</h3>
    <el-table :data="logs" stripe>
      <el-table-column prop="created_at" label="时间" width="170" />
      <el-table-column prop="change_type" label="类型" width="130">
        <template #default="{ row }">{{ typeLabel(row.change_type) }}</template>
      </el-table-column>
      <el-table-column prop="amount" label="金额" width="100" />
      <el-table-column prop="balance_after" label="余额后" width="100" />
      <el-table-column prop="remark" label="备注" min-width="160" />
    </el-table>
    <el-pagination
      class="pager"
      layout="prev, pager, next, total"
      :total="total"
      v-model:current-page="page"
      :page-size="pageSize"
      @current-change="loadLogs"
    />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { fetchMerchantWalletLogs, getMerchantWallet } from '../../api/seckill-wallet'

const loading = ref(false)
const wallet = ref({})
const logs = ref([])
const page = ref(1)
const pageSize = 20
const total = ref(0)

function typeLabel(t) {
  return { admin_adjust: '平台调账', seckill_apply: '秒杀报名' }[t] || t
}

async function loadWallet() {
  loading.value = true
  try {
    const res = await getMerchantWallet()
    wallet.value = res || {}
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function loadLogs() {
  try {
    const res = await fetchMerchantWalletLogs({ page: page.value, page_size: pageSize })
    logs.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(() => {
  loadWallet()
  loadLogs()
})
</script>

<style scoped>
.nums { display: flex; gap: 16px; margin: 16px 0 24px; flex-wrap: wrap; }
.card {
  min-width: 140px; padding: 16px 20px; background: #f8fafc; border-radius: 8px;
  display: flex; flex-direction: column; gap: 6px;
}
.card span { color: #64748b; font-size: 13px; }
.card b { font-size: 22px; }
.pager { margin-top: 12px; }
</style>
