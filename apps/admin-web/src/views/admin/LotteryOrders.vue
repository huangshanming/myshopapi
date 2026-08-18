<template>
  <div>
    <div class="toolbar">
      <h2>抽奖实物订单</h2>
      <el-select v-model="q.fulfill_status" clearable placeholder="履约状态" style="width: 140px" @change="search">
        <el-option v-for="o in STATUS_OPTIONS" :key="o.value" :label="o.label" :value="o.value" />
      </el-select>
      <el-button type="primary" @click="search">查询</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="用户" width="100">
        <template #default="{ row }">#{{ row.user_id }}</template>
      </el-table-column>
      <el-table-column prop="prize_name" label="奖品" min-width="140" />
      <el-table-column label="收货信息" min-width="220">
        <template #default="{ row }">
          <div v-if="row.receiver_name || row.receiver_phone">
            {{ row.receiver_name }} {{ row.receiver_phone }}
          </div>
          <div class="sub">{{ row.receiver_address || (row.fulfill_status === 'need_address' ? '待用户填写地址' : '-') }}</div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusType(row.fulfill_status)" size="small">{{ statusLabel(row.fulfill_status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="物流" min-width="160">
        <template #default="{ row }">
          <span v-if="row.ship_no">{{ row.ship_company || '' }} {{ row.ship_no }}</span>
          <span v-else class="sub">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="中奖时间" width="170" />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.fulfill_status === 'pending'"
            v-permission="'marketing:lottery:order'"
            link
            type="warning"
            @click="openShip(row)"
          >发货</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      class="pager"
      layout="prev,pager,next,total"
      :total="total"
      v-model:current-page="page"
      :page-size="20"
      @current-change="load"
    />

    <el-dialog v-model="shipVisible" title="发货" width="420px">
      <el-form label-width="90px">
        <el-form-item label="物流公司"><el-input v-model="shipForm.ship_company" placeholder="如：顺丰速运" /></el-form-item>
        <el-form-item label="物流单号" required><el-input v-model="shipForm.ship_no" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitShip">确认发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { listLotteryOrders, shipLotteryOrder } from '../../api/lottery'

const STATUS_OPTIONS = [
  { value: 'need_address', label: '待填地址' },
  { value: 'pending', label: '待发货' },
  { value: 'shipped', label: '已发货' },
]

const list = ref([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)
const q = reactive({ fulfill_status: '' })

const shipVisible = ref(false)
const shipForm = reactive({ id: 0, ship_company: '', ship_no: '' })
const submitting = ref(false)

function statusLabel(s) {
  return STATUS_OPTIONS.find((o) => o.value === s)?.label || s || '-'
}
function statusType(s) {
  if (s === 'need_address') return 'info'
  if (s === 'pending') return 'warning'
  if (s === 'shipped') return 'success'
  return 'info'
}

function search() {
  page.value = 1
  load()
}

async function load() {
  loading.value = true
  try {
    const res = await listLotteryOrders({
      page: page.value,
      page_size: 20,
      fulfill_status: q.fulfill_status || undefined,
    })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function openShip(row) {
  shipForm.id = row.id
  shipForm.ship_company = ''
  shipForm.ship_no = ''
  shipVisible.value = true
}

async function submitShip() {
  if (!shipForm.ship_no.trim()) {
    ElMessage.warning('请填写物流单号')
    return
  }
  submitting.value = true
  try {
    await shipLotteryOrder(shipForm.id, {
      ship_company: shipForm.ship_company,
      ship_no: shipForm.ship_no,
    })
    ElMessage.success('已发货')
    shipVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message || '发货失败')
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 10px; align-items: center; margin-bottom: 12px; flex-wrap: wrap; }
.pager { margin-top: 16px; justify-content: flex-end; }
.sub { color: #94a3b8; font-size: 12px; margin-top: 2px; }
</style>
