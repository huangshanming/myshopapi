<template>
  <div>
    <div class="toolbar">
      <h2>文章列表</h2>
      <el-select
        v-model="q.shop_id"
        filterable
        remote
        clearable
        reserve-keyword
        placeholder="归属"
        :remote-method="searchShops"
        :loading="shopLoading"
        style="width: 200px"
        @change="load"
      >
        <el-option label="平台官方" :value="0" />
        <el-option v-for="s in shops" :key="s.id" :label="`${s.name} (#${s.id})`" :value="s.id" />
      </el-select>
      <el-input v-model="q.title" placeholder="标题" clearable style="width: 160px" @keyup.enter="load" />
      <el-select v-model="q.audit_status" clearable placeholder="审核" style="width: 120px" @change="load">
        <el-option label="待审" value="pending" />
        <el-option label="已通过" value="approved" />
        <el-option label="已驳回" value="rejected" />
      </el-select>
      <el-select v-model="q.has_schedule" clearable placeholder="定时" style="width: 110px" @change="load">
        <el-option label="有定时" value="1" />
        <el-option label="无定时" value="0" />
      </el-select>
      <el-date-picker v-model="dateRange" type="daterange" value-format="YYYY-MM-DD" start-placeholder="开始" end-placeholder="结束" />
      <el-button @click="load">查询</el-button>
      <el-button v-permission="'community:article:add'" type="primary" @click="$router.push('/admin/articles/edit')">发布文章</el-button>
      <el-button v-permission="'community:article:audit'" :disabled="!selected.length" @click="batchPass">批量通过</el-button>
      <el-button v-permission="'community:article:audit'" :disabled="!selected.length" @click="batchReject">批量驳回</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe @selection-change="(rows) => (selected = rows)">
      <el-table-column type="selection" width="45" />
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="归属" width="120">
        <template #default="{ row }">
          <el-tag v-if="!row.shop_id" type="warning" size="small">平台官方</el-tag>
          <span v-else>#{{ row.shop_id }}</span>
        </template>
      </el-table-column>
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
          <el-button link @click="openRemark(row, 'offline')">下架</el-button>
          <el-button link type="danger" @click="openRemark(row, 'delete')">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pager" layout="prev, pager, next, total" :total="total" v-model:current-page="page" :page-size="pageSize" @current-change="load" />

    <el-dialog v-model="remarkVisible" :title="remarkTitle" width="480px">
      <el-form label-width="80px">
        <el-form-item label="文章">
          <span>{{ remarkRow?.title }} (#{{ remarkRow?.id }})</span>
        </el-form-item>
        <el-form-item v-if="needRemark" label="备注" required>
          <el-input v-model="remark" type="textarea" :rows="3" placeholder="将展示给商家，请说明原因" maxlength="500" show-word-limit />
        </el-form-item>
        <el-alert v-else type="info" :closable="false" title="平台官方文章，确认后直接处理，无需通知商家。" />
      </el-form>
      <template #footer>
        <el-button @click="remarkVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitRemark">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { fetchShops } from '../../../api/merchant'
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
const shops = ref([])
const shopLoading = ref(false)
const q = reactive({ shop_id: undefined, title: '', audit_status: '', has_schedule: '' })

const remarkVisible = ref(false)
const remarkAction = ref('delete')
const remarkRow = ref(null)
const remark = ref('')
const submitting = ref(false)

const needRemark = computed(() => !!(remarkRow.value?.shop_id))
const remarkTitle = computed(() => {
  if (remarkAction.value === 'delete') {
    return needRemark.value ? '删除并通知商家' : '删除平台文章'
  }
  return needRemark.value ? '下架并通知商家' : '下架平台文章'
})

async function searchShops(keyword) {
  shopLoading.value = true
  try {
    const res = await fetchShops({ name: keyword || undefined, page: 1, page_size: 50 })
    shops.value = res.data?.list || res.data?.items || []
  } finally {
    shopLoading.value = false
  }
}

async function load() {
  loading.value = true
  try {
    const params = {
      page: page.value, page_size: pageSize,
      title: q.title || undefined,
      audit_status: q.audit_status || undefined,
      has_schedule: q.has_schedule || undefined,
      created_from: dateRange.value?.[0],
      created_to: dateRange.value?.[1],
    }
    if (q.shop_id !== undefined && q.shop_id !== null && q.shop_id !== '') {
      params.shop_id = q.shop_id
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
    const { value } = await ElMessageBox.prompt('驳回理由（将通知商家）', '审核驳回', { inputPattern: /.+/, inputErrorMessage: '必填' })
    reject_reason = value
  }
  await auditArticle(row.id, { pass, reject_reason })
  ElMessage.success(pass ? '已通过' : '已驳回并通知商家')
  load()
}

async function batchPass() {
  await batchAuditArticles({ ids: selected.value.map((r) => r.id), pass: true })
  ElMessage.success('批量通过')
  load()
}

async function batchReject() {
  const { value } = await ElMessageBox.prompt('驳回理由（将通知商家）', '批量驳回', { inputPattern: /.+/, inputErrorMessage: '必填' })
  await batchAuditArticles({ ids: selected.value.map((r) => r.id), pass: false, reject_reason: value })
  ElMessage.success('批量驳回并通知商家')
  load()
}

async function toggleTop(row) {
  await topArticle(row.id, row.is_top ? 0 : 1)
  load()
}

function openRemark(row, action) {
  remarkRow.value = row
  remarkAction.value = action
  remark.value = ''
  remarkVisible.value = true
}

async function submitRemark() {
  if (needRemark.value && !remark.value.trim()) {
    ElMessage.warning('请填写备注')
    return
  }
  submitting.value = true
  try {
    const text = remark.value.trim()
    if (remarkAction.value === 'delete') {
      await deleteArticle(remarkRow.value.id, text)
      ElMessage.success(needRemark.value ? '已删除并通知商家' : '已删除')
    } else {
      await offlineArticle(remarkRow.value.id, text)
      ElMessage.success(needRemark.value ? '已下架并通知商家' : '已下架')
    }
    remarkVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  await searchShops('')
  load()
})
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.pager { margin-top: 12px; justify-content: flex-end; }
</style>
