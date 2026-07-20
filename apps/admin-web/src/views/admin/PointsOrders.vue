<template>
  <div>
    <div class="toolbar">
      <h2>积分订单</h2>
      <el-input v-model="q.order_no" placeholder="订单号" clearable style="width: 180px" @keyup.enter="search" />
      <el-input v-model="q.keyword" placeholder="商品名" clearable style="width: 160px" @keyup.enter="search" />
      <el-select v-model="q.status" clearable placeholder="状态" style="width: 120px" @change="search">
        <el-option v-for="o in STATUS_OPTIONS" :key="o.value" :label="o.label" :value="o.value" />
      </el-select>
      <el-button type="primary" @click="search">查询</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="order_no" label="订单号" min-width="180" />
      <el-table-column label="用户" min-width="140">
        <template #default="{ row }">
          <div>{{ row.user_name || '-' }}</div>
          <div class="sub">#{{ row.user_id }} {{ row.user_mobile || '' }}</div>
        </template>
      </el-table-column>
      <el-table-column label="商品" min-width="200">
        <template #default="{ row }">
          <div class="prod">
            <el-image v-if="row.product_cover" :src="row.product_cover" fit="cover" class="cover" />
            <div>
              <div>{{ row.product_name }}</div>
              <div class="sub">×{{ row.quantity }} · {{ row.points_cost }} 积分</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="收货信息" min-width="180">
        <template #default="{ row }">
          <div v-if="row.receiver_name || row.receiver_phone">
            {{ row.receiver_name }} {{ row.receiver_phone }}
          </div>
          <div class="sub">{{ row.receiver_address || '-' }}</div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="下单时间" min-width="160" />
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openDetail(row)">详情</el-button>
          <el-button
            v-if="row.status === 'pending'"
            v-permission="'marketing:points_mall:edit'"
            link
            type="warning"
            @click="openShip(row)"
          >发货</el-button>
          <el-button
            v-if="row.status === 'pending' || row.status === 'shipped'"
            v-permission="'marketing:points_mall:edit'"
            link
            type="success"
            @click="doComplete(row)"
          >完成</el-button>
          <el-button
            v-if="row.status === 'pending'"
            v-permission="'marketing:points_mall:edit'"
            link
            type="danger"
            @click="doCancel(row)"
          >取消退积分</el-button>
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

    <el-drawer v-model="detailVisible" title="兑换订单详情" size="480px">
      <template v-if="detail">
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="订单号">{{ detail.order_no }}</el-descriptions-item>
          <el-descriptions-item label="用户">{{ detail.user_name }} (#{{ detail.user_id }}) {{ detail.user_mobile }}</el-descriptions-item>
          <el-descriptions-item label="商品">{{ detail.product_name }} ×{{ detail.quantity }}</el-descriptions-item>
          <el-descriptions-item label="积分">{{ detail.points_cost }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusType(detail.status)" size="small">{{ statusLabel(detail.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="收货人">{{ detail.receiver_name || '-' }} {{ detail.receiver_phone }}</el-descriptions-item>
          <el-descriptions-item label="地址">{{ detail.receiver_address || '-' }}</el-descriptions-item>
          <el-descriptions-item label="物流">{{ detail.ship_company || '-' }} {{ detail.ship_no }}</el-descriptions-item>
          <el-descriptions-item label="备注">{{ detail.admin_remark || '-' }}</el-descriptions-item>
          <el-descriptions-item label="下单时间">{{ detail.created_at }}</el-descriptions-item>
        </el-descriptions>
        <div class="drawer-ops">
          <el-input v-model="remarkText" placeholder="管理员备注" />
          <el-button v-permission="'marketing:points_mall:edit'" type="primary" @click="saveRemark">保存备注</el-button>
        </div>
      </template>
    </el-drawer>

    <el-dialog v-model="shipVisible" title="发货" width="420px">
      <el-form label-width="90px">
        <el-form-item label="物流公司"><el-input v-model="shipForm.ship_company" /></el-form-item>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  cancelPointsOrder,
  completePointsOrder,
  getPointsOrder,
  listPointsOrders,
  remarkPointsOrder,
  shipPointsOrder,
} from '../../api/points-mall'

const STATUS_OPTIONS = [
  { value: 'pending', label: '待发货' },
  { value: 'shipped', label: '已发货' },
  { value: 'completed', label: '已完成' },
  { value: 'cancelled', label: '已取消' },
]

const list = ref([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)
const q = reactive({ order_no: '', keyword: '', status: '' })

const detailVisible = ref(false)
const detail = ref(null)
const remarkText = ref('')

const shipVisible = ref(false)
const shipForm = reactive({ id: 0, ship_company: '', ship_no: '' })
const submitting = ref(false)

function statusLabel(s) {
  return STATUS_OPTIONS.find((o) => o.value === s)?.label || s
}
function statusType(s) {
  if (s === 'pending') return 'warning'
  if (s === 'shipped') return ''
  if (s === 'completed') return 'success'
  if (s === 'cancelled') return 'info'
  return 'info'
}

function search() {
  page.value = 1
  load()
}

async function load() {
  loading.value = true
  try {
    const res = await listPointsOrders({
      page: page.value,
      page_size: 20,
      order_no: q.order_no || undefined,
      keyword: q.keyword || undefined,
      status: q.status || undefined,
    })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function openDetail(row) {
  try {
    detail.value = await getPointsOrder(row.id)
    remarkText.value = detail.value?.admin_remark || ''
    detailVisible.value = true
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
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
    await shipPointsOrder(shipForm.id, {
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

async function doComplete(row) {
  try {
    await ElMessageBox.confirm(`确认完成订单 ${row.order_no}？`, '完成订单')
    await completePointsOrder(row.id)
    ElMessage.success('已完成')
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || '操作失败')
  }
}

async function doCancel(row) {
  try {
    await ElMessageBox.confirm(`取消后将退回 ${row.points_cost} 积分并恢复库存，确认？`, '取消订单', { type: 'warning' })
    await cancelPointsOrder(row.id, { admin_remark: '后台取消' })
    ElMessage.success('已取消并退积分')
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || '操作失败')
  }
}

async function saveRemark() {
  if (!detail.value) return
  try {
    detail.value = await remarkPointsOrder(detail.value.id, { admin_remark: remarkText.value })
    ElMessage.success('备注已保存')
    load()
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 10px; align-items: center; margin-bottom: 16px; }
.toolbar h2 { margin: 0 12px 0 0; }
.pager { margin-top: 16px; justify-content: flex-end; }
.sub { color: #94a3b8; font-size: 12px; }
.prod { display: flex; gap: 10px; align-items: center; }
.cover { width: 44px; height: 44px; border-radius: 6px; flex-shrink: 0; }
.drawer-ops { margin-top: 16px; display: flex; gap: 8px; }
</style>
