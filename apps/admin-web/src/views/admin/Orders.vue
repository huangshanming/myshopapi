<template>
  <div>
    <div class="toolbar">
      <h2>全站订单</h2>
      <el-input v-model="q.order_no" placeholder="订单号" clearable style="width: 180px" @keyup.enter="search" />
      <el-select v-model="q.status" clearable placeholder="状态" style="width: 120px" @change="search">
        <el-option v-for="o in ORDER_STATUS_OPTIONS" :key="o.value" :label="o.label" :value="o.value" />
      </el-select>
      <el-button type="primary" @click="search">查询</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="order_no" label="订单号" min-width="170" />
      <el-table-column label="用户" min-width="140">
        <template #default="{ row }">
          <div>{{ row.user_name || '-' }}</div>
          <div class="sub">#{{ row.user_id }}</div>
        </template>
      </el-table-column>
      <el-table-column label="店铺" min-width="140">
        <template #default="{ row }">
          <div>{{ row.shop_name || '-' }}</div>
          <div class="sub">#{{ row.shop_id }}</div>
        </template>
      </el-table-column>
      <el-table-column label="商品" min-width="220">
        <template #default="{ row }">
          <div v-if="row.items?.length" class="items">
            <div v-for="it in row.items" :key="it.id" class="item-line">
              <a class="prod-link" href="javascript:;" @click.prevent="goProduct(it.product_id)">
                {{ it.product_name }}
              </a>
              <span class="sub"> ×{{ it.quantity }} · ¥{{ it.price }}</span>
              <div v-if="formatSkuSnapshot(it.sku_snapshot)" class="sub snap">
                {{ formatSkuSnapshot(it.sku_snapshot) }}
              </div>
            </div>
          </div>
          <span v-else class="sub">-</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="orderStatusType(row.status)" size="small">{{ orderStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="total_amount" label="金额" width="100" />
      <el-table-column prop="created_at" label="创建时间" min-width="160" />
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openDetail(row)">详情</el-button>
          <el-button
            v-if="row.status === 'confirmed'"
            v-permission="'business:order:ship'"
            link
            type="warning"
            @click="openShip(row)"
          >发货</el-button>
          <el-button
            v-if="row.status === 'shipped'"
            v-permission="'business:order:complete'"
            link
            type="success"
            @click="doComplete(row)"
          >完成</el-button>
          <el-button v-permission="'business:order:remark'" link @click="openRemark(row)">备注</el-button>
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

    <el-drawer v-model="detailVisible" title="订单详情" size="520px">
      <template v-if="detail">
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="订单号">{{ detail.order_no }}</el-descriptions-item>
          <el-descriptions-item label="用户">{{ detail.user_name }} (#{{ detail.user_id }})</el-descriptions-item>
          <el-descriptions-item label="店铺">{{ detail.shop_name }} (#{{ detail.shop_id }})</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="orderStatusType(detail.status)" size="small">{{ orderStatusLabel(detail.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="金额">{{ detail.total_amount }}</el-descriptions-item>
          <el-descriptions-item label="收货人">{{ detail.receiver_name || '-' }} {{ detail.receiver_phone }}</el-descriptions-item>
          <el-descriptions-item label="地址">{{ detail.receiver_address || '-' }}</el-descriptions-item>
          <el-descriptions-item label="物流">{{ detail.ship_company || '-' }} {{ detail.ship_no }}</el-descriptions-item>
          <el-descriptions-item label="备注">{{ detail.remark || '-' }}</el-descriptions-item>
        </el-descriptions>
        <h4 class="sec">商品明细</h4>
        <el-table :data="detail.items || []" size="small" stripe>
          <el-table-column label="商品" min-width="160">
            <template #default="{ row }">
              <a class="prod-link" href="javascript:;" @click.prevent="goProduct(row.product_id)">
                {{ row.product_name }}
              </a>
              <div v-if="formatSkuSnapshot(row.sku_snapshot)" class="sub snap">
                {{ formatSkuSnapshot(row.sku_snapshot) }}
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="sku_id" label="SKU" width="70" />
          <el-table-column prop="price" label="单价" width="80" />
          <el-table-column prop="quantity" label="数量" width="60" />
        </el-table>
        <h4 class="sec">售后记录</h4>
        <el-table :data="afterSales" size="small" stripe empty-text="暂无">
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="amount" label="金额" width="80" />
          <el-table-column label="状态" width="90">
            <template #default="{ row }">{{ afterSaleStatusLabel(row.status) }}</template>
          </el-table-column>
          <el-table-column prop="reason" label="原因" min-width="120" show-overflow-tooltip />
        </el-table>
      </template>
    </el-drawer>

    <el-dialog v-model="shipVisible" title="发货" width="420px">
      <el-form label-width="90px">
        <el-form-item label="物流公司" required>
          <el-select
            v-model="shipForm.ship_company"
            filterable
            remote
            clearable
            reserve-keyword
            placeholder="搜索并选择物流公司"
            :remote-method="searchLogistics"
            :loading="logisticsLoading"
            style="width: 100%"
          >
            <el-option
              v-for="c in logisticsOptions"
              :key="c.id"
              :label="`${c.name} (${c.code})`"
              :value="c.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="物流单号" required>
          <el-input v-model="shipForm.ship_no" placeholder="运单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitShip">确认发货</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="remarkVisible" title="订单备注" width="420px">
      <el-input v-model="remarkText" type="textarea" :rows="3" maxlength="255" show-word-limit />
      <template #footer>
        <el-button @click="remarkVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitRemark">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listOrders, getOrder, shipOrder, completeOrder, remarkOrder, listLogisticsOptions,
  ORDER_STATUS_OPTIONS, orderStatusLabel, orderStatusType, afterSaleStatusLabel, formatSkuSnapshot,
} from '../../api/order'

const router = useRouter()
const SCOPE = 'admin'
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const q = reactive({ order_no: '', status: '' })

const detailVisible = ref(false)
const detail = ref(null)
const afterSales = ref([])

const shipVisible = ref(false)
const shipRow = ref(null)
const shipForm = reactive({ ship_company: '', ship_no: '' })
const logisticsOptions = ref([])
const logisticsLoading = ref(false)

const remarkVisible = ref(false)
const remarkRow = ref(null)
const remarkText = ref('')
const submitting = ref(false)

function search() {
  page.value = 1
  load()
}

function goProduct(productId) {
  if (!productId) return
  router.push(`/admin/products/detail/${productId}`)
}

async function load() {
  loading.value = true
  try {
    const res = await listOrders(SCOPE, {
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

async function openDetail(row) {
  try {
    const res = await getOrder(SCOPE, row.id)
    detail.value = res.data?.order || res.data
    afterSales.value = res.data?.after_sales || []
    detailVisible.value = true
  } catch (e) {
    ElMessage.error(e.message)
  }
}

function openShip(row) {
  shipRow.value = row
  shipForm.ship_company = ''
  shipForm.ship_no = ''
  shipVisible.value = true
  searchLogistics('')
}

async function searchLogistics(keyword) {
  logisticsLoading.value = true
  try {
    const res = await listLogisticsOptions(keyword)
    logisticsOptions.value = Array.isArray(res.data) ? res.data : (res.data?.list || [])
  } catch (_) {
    logisticsOptions.value = []
  } finally {
    logisticsLoading.value = false
  }
}

async function submitShip() {
  if (!shipForm.ship_company || !shipForm.ship_no.trim()) {
    ElMessage.warning('请选择物流公司并填写单号')
    return
  }
  submitting.value = true
  try {
    await shipOrder(SCOPE, shipRow.value.id, {
      ship_company: shipForm.ship_company.trim(),
      ship_no: shipForm.ship_no.trim(),
    })
    ElMessage.success('发货成功')
    shipVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    submitting.value = false
  }
}

async function doComplete(row) {
  try {
    await ElMessageBox.confirm(`确认将订单 ${row.order_no} 标记为已完成？`, '完成订单')
    await completeOrder(SCOPE, row.id)
    ElMessage.success('已完成')
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message)
  }
}

function openRemark(row) {
  remarkRow.value = row
  remarkText.value = row.remark || ''
  remarkVisible.value = true
}

async function submitRemark() {
  submitting.value = true
  try {
    await remarkOrder(SCOPE, remarkRow.value.id, remarkText.value)
    ElMessage.success('已更新')
    remarkVisible.value = false
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
.sec { margin: 16px 0 8px; font-size: 14px; }
.items { display: flex; flex-direction: column; gap: 6px; }
.item-line { line-height: 1.35; }
.prod-link { color: var(--el-color-primary); text-decoration: none; }
.prod-link:hover { text-decoration: underline; }
.snap { margin-top: 2px; }
</style>
