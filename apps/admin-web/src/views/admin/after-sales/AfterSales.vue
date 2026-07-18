<template>
  <div>
    <div class="toolbar">
      <h2>售后管理</h2>
      <el-input v-model="q.order_no" placeholder="订单号" clearable style="width: 180px" @keyup.enter="search" />
      <el-select v-model="q.status" clearable placeholder="状态" style="width: 120px" @change="search">
        <el-option v-for="o in AFTER_SALE_STATUS_OPTIONS" :key="o.value" :label="o.label" :value="o.value" />
      </el-select>
      <el-button type="primary" @click="search">查询</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="order_no" label="订单号" min-width="160" />
      <el-table-column label="用户" min-width="120">
        <template #default="{ row }">
          <div>{{ row.user_name || '-' }}</div>
          <div class="sub">#{{ row.user_id }}</div>
        </template>
      </el-table-column>
      <el-table-column label="店铺" min-width="120">
        <template #default="{ row }">
          <div>{{ row.shop_name || '-' }}</div>
          <div class="sub">#{{ row.shop_id }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="类型" width="110" />
      <el-table-column prop="amount" label="金额" width="90" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag size="small">{{ afterSaleStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="reason" label="原因" min-width="140" show-overflow-tooltip />
      <el-table-column prop="created_at" label="申请时间" min-width="160" />
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <template v-if="row.status === 'pending'">
            <el-button v-permission="'business:aftersale:handle'" link type="success" @click="doHandle(row, 'approve')">同意</el-button>
            <el-button v-permission="'business:aftersale:handle'" link type="danger" @click="doHandle(row, 'reject')">拒绝</el-button>
          </template>
          <el-button
            v-if="row.status === 'pending' || row.status === 'approved'"
            v-permission="'business:aftersale:handle'"
            link
            type="warning"
            @click="doHandle(row, 'refunded')"
          >退款完成</el-button>
          <el-button
            v-if="row.status !== 'closed' && row.status !== 'refunded'"
            v-permission="'business:aftersale:handle'"
            link
            @click="doHandle(row, 'closed')"
          >关闭</el-button>
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

    <el-dialog v-model="handleVisible" :title="handleTitle" width="420px">
      <el-form label-width="80px">
        <el-form-item label="备注">
          <el-input v-model="adminRemark" type="textarea" :rows="3" maxlength="500" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="handleVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitHandle">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  listAfterSales, handleAfterSale,
  AFTER_SALE_STATUS_OPTIONS, afterSaleStatusLabel,
} from '../../../api/order'

const SCOPE = 'admin'
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const q = reactive({ order_no: '', status: '' })

const handleVisible = ref(false)
const handleRow = ref(null)
const handleAction = ref('')
const adminRemark = ref('')
const submitting = ref(false)

const handleTitle = computed(() => {
  const map = { approve: '同意售后', reject: '拒绝售后', refunded: '确认退款完成', closed: '关闭售后' }
  return map[handleAction.value] || '处理售后'
})

function search() {
  page.value = 1
  load()
}

async function load() {
  loading.value = true
  try {
    const res = await listAfterSales(SCOPE, {
      page: page.value,
      page_size: pageSize,
      order_no: q.order_no || undefined,
      status: q.status || undefined,
    })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
    list.value = []
  } finally {
    loading.value = false
  }
}

function doHandle(row, action) {
  handleRow.value = row
  handleAction.value = action
  adminRemark.value = ''
  handleVisible.value = true
}

async function submitHandle() {
  submitting.value = true
  try {
    await handleAfterSale(SCOPE, handleRow.value.id, {
      action: handleAction.value,
      admin_remark: adminRemark.value,
    })
    ElMessage.success('处理成功')
    handleVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.pager { margin-top: 16px; justify-content: flex-end; }
.sub { color: #909399; font-size: 12px; }
</style>
