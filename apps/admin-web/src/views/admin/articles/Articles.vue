<template>
  <div>
    <div class="toolbar">
      <h2>文章列表</h2>
      <el-input v-model="q.shop_id" placeholder="商家ID" clearable style="width: 110px" />
      <el-input v-model="q.title" placeholder="标题" clearable style="width: 160px" />
      <el-select v-model="q.audit_status" clearable placeholder="审核" style="width: 120px">
        <el-option label="待审" value="pending" />
        <el-option label="已通过" value="approved" />
        <el-option label="已驳回" value="rejected" />
      </el-select>
      <el-select v-model="q.has_schedule" clearable placeholder="定时" style="width: 110px">
        <el-option label="有定时" value="1" />
        <el-option label="无定时" value="0" />
      </el-select>
      <el-date-picker v-model="dateRange" type="daterange" value-format="YYYY-MM-DD" start-placeholder="开始" end-placeholder="结束" />
      <el-button @click="load">查询</el-button>
      <el-button v-permission="'community:article:add'" type="primary" @click="$router.push('/admin/articles/edit')">新增</el-button>
      <el-button v-permission="'community:article:audit'" :disabled="!selected.length" @click="batchPass">批量通过</el-button>
      <el-button v-permission="'community:article:audit'" :disabled="!selected.length" @click="batchReject">批量驳回</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe @selection-change="(rows) => (selected = rows)">
      <el-table-column type="selection" width="45" />
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="shop_id" label="商家" width="80" />
      <el-table-column prop="title" label="标题" min-width="180" />
      <el-table-column prop="audit_status" label="审核" width="90" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column prop="is_top" label="置顶" width="70">
        <template #default="{ row }">{{ row.is_top ? '是' : '' }}</template>
      </el-table-column>
      <el-table-column prop="schedule_publish_at" label="定时发布" width="160" />
      <el-table-column prop="created_at" label="创建时间" width="160" />
      <el-table-column label="操作" width="320" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="$router.push('/admin/articles/edit/' + row.id)">编辑</el-button>
          <el-button v-if="row.audit_status === 'pending'" v-permission="'community:article:audit'" link type="success" @click="doAudit(row, true)">通过</el-button>
          <el-button v-if="row.audit_status === 'pending'" v-permission="'community:article:audit'" link type="warning" @click="doAudit(row, false)">驳回</el-button>
          <el-button v-permission="'community:article:top'" link @click="toggleTop(row)">{{ row.is_top ? '取消置顶' : '置顶' }}</el-button>
          <el-button link @click="doOffline(row)">下架</el-button>
          <el-button link type="danger" @click="doDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pager" layout="prev, pager, next, total" :total="total" v-model:current-page="page" :page-size="pageSize" @current-change="load" />
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listArticles, auditArticle, batchAuditArticles, topArticle, offlineArticle, deleteArticle,
} from '../../../api/admin-article'

const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const selected = ref([])
const dateRange = ref([])
const q = reactive({ shop_id: '', title: '', audit_status: '', has_schedule: '' })

async function load() {
  loading.value = true
  try {
    const params = {
      page: page.value, page_size: pageSize,
      shop_id: q.shop_id || undefined,
      title: q.title || undefined,
      audit_status: q.audit_status || undefined,
      has_schedule: q.has_schedule || undefined,
      created_from: dateRange.value?.[0],
      created_to: dateRange.value?.[1],
    }
    const res = await listArticles(params)
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function doAudit(row, pass) {
  let reject_reason = ''
  if (!pass) {
    const { value } = await ElMessageBox.prompt('驳回理由', '审核驳回', { inputPattern: /.+/, inputErrorMessage: '必填' })
    reject_reason = value
  }
  await auditArticle(row.id, { pass, reject_reason })
  ElMessage.success('已处理')
  load()
}

async function batchPass() {
  await batchAuditArticles({ ids: selected.value.map((r) => r.id), pass: true })
  ElMessage.success('批量通过')
  load()
}

async function batchReject() {
  const { value } = await ElMessageBox.prompt('驳回理由', '批量驳回', { inputPattern: /.+/, inputErrorMessage: '必填' })
  await batchAuditArticles({ ids: selected.value.map((r) => r.id), pass: false, reject_reason: value })
  ElMessage.success('批量驳回')
  load()
}

async function toggleTop(row) {
  await topArticle(row.id, row.is_top ? 0 : 1)
  load()
}

async function doOffline(row) {
  await offlineArticle(row.id)
  ElMessage.success('已下架')
  load()
}

async function doDelete(row) {
  await ElMessageBox.confirm('移入回收站？')
  await deleteArticle(row.id)
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.pager { margin-top: 12px; justify-content: flex-end; }
</style>
