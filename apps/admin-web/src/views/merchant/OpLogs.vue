<template>
  <div>
    <div class="toolbar">
      <h2>操作日志</h2>
      <el-button @click="load">刷新</el-button>
    </div>
    <el-alert
      v-if="!loading && !list.length"
      type="info"
      :closable="false"
      show-icon
      title="暂无日志。保存 / 上下架 / 复制商品后会出现在这里。"
      style="margin-bottom: 12px"
    />
    <el-table :data="list" v-loading="loading" stripe @row-click="onRowClick">
      <el-table-column label="时间" width="180">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作类型" width="120">
        <template #default="{ row }">
          <el-tag size="small" :type="tagType(row.action)">{{ actionLabel(row.action) || row.action_label || '-' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="对象" min-width="200">
        <template #default="{ row }">
          <div class="target">
            <span class="type">{{ row.target_type || '商品' }}</span>
            <el-button
              v-if="canOpen(row)"
              link
              type="primary"
              @click.stop="goProduct(row)"
            >{{ row.target_name || row.product_name || ('商品#' + row.product_id) }}</el-button>
            <span v-else>{{ row.target_name || row.product_name || '-' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="操作摘要" min-width="260" show-overflow-tooltip>
        <template #default="{ row }">{{ summaryText(row) }}</template>
      </el-table-column>
      <el-table-column label="操作人" width="120">
        <template #default="{ row }">{{ row.operator_name || (row.operator_id ? ('用户#' + row.operator_id) : '-') }}</template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button v-if="canOpen(row)" link type="primary" size="small" @click.stop="goProduct(row)">查看商品</el-button>
          <span v-else class="muted">-</span>
        </template>
      </el-table-column>
    </el-table>
    <div class="pager">
      <el-pagination
        v-model:current-page="page"
        layout="total, prev, pager, next"
        :total="total"
        @current-change="load"
      />
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { pickList } from '../../utils/list'
import { listOpLogs } from '../../api/merchant-product'

const router = useRouter()
const list = ref([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)

const ACTION_MAP = {
  create: '创建商品',
  update: '更新商品',
  copy: '复制商品',
  schedule: '设置定时',
  permanent_delete: '永久删除',
  'status:on_sale': '上架',
  'status:off_sale': '下架',
  'status:deleted': '移入回收站',
  'status:draft': '设为草稿',
  save: '保存商品',
  platform_off_sale: '平台强制下架',
  platform_delete: '平台删除',
}

function actionLabel(action) {
  return ACTION_MAP[action] || action || '-'
}

function tagType(action) {
  if (action === 'create' || action === 'status:on_sale') return 'success'
  if (action === 'status:off_sale' || action === 'status:draft' || action === 'platform_off_sale') return 'warning'
  if (action === 'status:deleted' || action === 'permanent_delete' || action === 'platform_delete') return 'danger'
  if (action === 'copy' || action === 'schedule') return 'info'
  return ''
}

function canOpen(row) {
  return !!row.product_id && row.action !== 'permanent_delete'
}

function goProduct(row) {
  if (!canOpen(row)) return
  router.push(`/merchant/products/edit/${row.product_id}`)
}

function onRowClick(row) {
  if (canOpen(row)) goProduct(row)
}

function formatTime(v) {
  if (!v) return '-'
  return String(v).replace('T', ' ').slice(0, 19)
}

function formatChange(row) {
  const raw = row.after_json
  if (raw == null || raw === '' || raw === 'null') return '-'
  if (typeof raw === 'object') return JSON.stringify(raw)
  try {
    const o = JSON.parse(raw)
    if (o.status) {
      const st = { on_sale: '上架', off_sale: '下架', deleted: '回收站', draft: '草稿' }[o.status] || o.status
      return `状态 → ${st}`
    }
    if (o.name) return `名称：${o.name}`
    return JSON.stringify(o)
  } catch {
    return String(raw)
  }
}

function summaryText(row) {
  const name = row.target_name || row.product_name || '商品'
  if (row.action === 'platform_off_sale') {
    let remark = ''
    try {
      remark = JSON.parse(row.after_json || '{}').remark || ''
    } catch (_) { /* ignore */ }
    return remark ? `平台强制下架「${name}」，备注：${remark}` : `平台强制下架「${name}」`
  }
  if (row.action === 'platform_delete') {
    let remark = ''
    try {
      remark = JSON.parse(row.after_json || '{}').remark || ''
    } catch (_) { /* ignore */ }
    return remark ? `平台将「${name}」移入回收站，备注：${remark}` : `平台将「${name}」移入回收站`
  }
  return row.summary || formatChange(row)
}

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
.target { display: flex; align-items: center; gap: 8px; }
.type {
  font-size: 12px; color: #64748b; background: #f1f5f9;
  padding: 2px 6px; border-radius: 4px; flex-shrink: 0;
}
.muted { color: #94a3b8; }
:deep(.el-table__row) { cursor: pointer; }
</style>
