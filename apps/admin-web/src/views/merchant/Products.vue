<template>
  <div>
    <div class="toolbar">
      <h2>商品管理</h2>
      <div class="actions">
        <el-button v-permission="'product:add'" type="primary" @click="$router.push('/merchant/products/edit')">新建商品</el-button>
        <el-button v-permission="'product:import'" @click="triggerImport">导入</el-button>
        <el-button v-permission="'product:import'" @click="doExport">导出</el-button>
        <input ref="fileRef" type="file" accept=".csv,.xlsx,.xls" hidden @change="onImport" />
      </div>
    </div>

    <el-form :inline="true" class="filter" @submit.prevent="load">
      <el-form-item label="名称"><el-input v-model="query.name" clearable /></el-form-item>
      <el-form-item label="货号"><el-input v-model="query.product_no" clearable /></el-form-item>
      <el-form-item label="状态">
        <el-select v-model="query.status" clearable style="width: 120px">
          <el-option label="草稿" value="draft" />
          <el-option label="上架" value="on_sale" />
          <el-option label="下架" value="off_sale" />
        </el-select>
      </el-form-item>
      <el-form-item label="类型">
        <el-select v-model="query.product_type" clearable style="width: 120px">
          <el-option label="实物" value="physical" />
          <el-option label="生鲜" value="fresh" />
          <el-option label="虚拟" value="virtual" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="load">查询</el-button>
      </el-form-item>
    </el-form>

    <div class="batch" v-if="selected.length">
      <span>已选 {{ selected.length }} 件</span>
      <el-button v-permission="'product:batch'" size="small" @click="batch('on_sale')">批量上架</el-button>
      <el-button v-permission="'product:batch'" size="small" @click="batch('off_sale')">批量下架</el-button>
      <el-button v-permission="'product:batch'" size="small" type="danger" @click="batch('recycle')">移入回收站</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe @selection-change="onSelect">
      <el-table-column type="selection" width="48" />
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="主图" width="72">
        <template #default="{ row }">
          <el-image v-if="row.main_image" :src="row.main_image" style="width: 40px; height: 40px" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column prop="product_no" label="货号" width="120" />
      <el-table-column prop="product_type" label="类型" width="80" />
      <el-table-column prop="sale_price" label="售价" width="90" />
      <el-table-column prop="stock" label="库存" width="80" />
      <el-table-column prop="collect_count" label="收藏" width="70" />
      <el-table-column prop="avg_rating" label="均分" width="70" />
      <el-table-column prop="review_count" label="评价" width="70" />
      <el-table-column prop="good_rate" label="好评率" width="80" />
      <el-table-column prop="status" label="状态" width="90" />
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button size="small" link @click="$router.push(`/merchant/products/edit/${row.id}`)">编辑</el-button>
          <el-button size="small" link @click="copy(row)">复制</el-button>
          <el-button
            v-if="row.status !== 'on_sale'"
            v-permission="'product:status'"
            size="small"
            link
            type="success"
            @click="setStatus(row, 'on_sale')"
          >上架</el-button>
          <el-button
            v-else
            v-permission="'product:status'"
            size="small"
            link
            type="warning"
            @click="setStatus(row, 'off_sale')"
          >下架</el-button>
          <el-button size="small" link @click="openSchedule(row)">定时</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pager">
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.page_size"
        layout="total, prev, pager, next"
        :total="total"
        @current-change="load"
      />
    </div>

    <el-dialog v-model="scheduleVisible" title="定时上下架" width="420px">
      <el-form label-width="90px">
        <el-form-item label="动作">
          <el-radio-group v-model="schedule.action">
            <el-radio value="on_sale">上架</el-radio>
            <el-radio value="off_sale">下架</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="执行时间">
          <el-date-picker v-model="schedule.run_at" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scheduleVisible = false">取消</el-button>
        <el-button type="primary" @click="submitSchedule">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { pickList } from '../../utils/list'
import {
  batchProducts,
  copyProduct,
  getBatchJob,
  importProducts,
  listProducts,
  scheduleProduct,
  setProductStatus,
} from '../../api/merchant-product'
import http from '../../api/http'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const selected = ref([])
const fileRef = ref()
const query = reactive({
  page: 1,
  page_size: 20,
  name: '',
  product_no: '',
  status: '',
  product_type: '',
})

const scheduleVisible = ref(false)
const schedule = reactive({ product_id: 0, action: 'on_sale', run_at: '' })

async function load() {
  loading.value = true
  try {
    const res = await listProducts({ ...query })
    list.value = pickList(res)
    total.value = Number(res?.total || list.value.length)
  } catch (e) {
    ElMessage.error(e.message)
    list.value = []
  } finally {
    loading.value = false
  }
}

function onSelect(rows) {
  selected.value = rows
}

async function setStatus(row, status) {
  await setProductStatus(row.id, status)
  ElMessage.success('已更新')
  load()
}

async function copy(row) {
  await copyProduct(row.id)
  ElMessage.success('已复制为草稿')
  load()
}

async function batch(action) {
  const product_ids = selected.value.map((r) => r.id)
  const res = await batchProducts({ action, product_ids })
  const jobId = res?.id || res?.job_id
  ElMessage.success(jobId ? `任务已提交 #${jobId}` : '已提交')
  if (jobId) {
    const poll = setInterval(async () => {
      try {
        const j = await getBatchJob(jobId)
        const st = j?.status
        if (st === 'success' || st === 'failed' || st === 'partial') {
          clearInterval(poll)
          ElMessage.info(`批量任务 ${st}`)
          load()
        }
      } catch {
        clearInterval(poll)
      }
    }, 1500)
  } else {
    load()
  }
}

function openSchedule(row) {
  schedule.product_id = row.id
  schedule.action = 'on_sale'
  schedule.run_at = ''
  scheduleVisible.value = true
}

async function submitSchedule() {
  if (!schedule.run_at) {
    ElMessage.warning('请选择时间')
    return
  }
  await scheduleProduct(schedule.product_id, { action: schedule.action, run_at: schedule.run_at })
  ElMessage.success('已设置定时')
  scheduleVisible.value = false
}

function triggerImport() {
  fileRef.value?.click()
}

async function onImport(e) {
  const file = e.target.files?.[0]
  e.target.value = ''
  if (!file) return
  await importProducts(file)
  ElMessage.success('导入任务已提交')
  load()
}

async function doExport() {
  const res = await http.get('/api/v1/merchant/products/export')
  const url = res?.url
  if (!url) {
    ElMessage.warning('未返回导出文件')
    return
  }
  window.open(url, '_blank')
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.actions { display: flex; gap: 8px; }
.filter { margin: 12px 0; }
.batch { margin-bottom: 8px; display: flex; gap: 8px; align-items: center; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
