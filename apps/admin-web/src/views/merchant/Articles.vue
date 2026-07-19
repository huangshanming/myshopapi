<template>
  <div>
    <div class="toolbar">
      <h2>我的文章</h2>
      <el-input v-model="q.title" placeholder="标题" clearable style="width: 160px" />
      <el-select v-model="q.audit_status" clearable placeholder="审核" style="width: 120px">
        <el-option label="待审" value="pending" />
        <el-option label="已通过" value="approved" />
        <el-option label="已驳回" value="rejected" />
      </el-select>
      <el-button @click="load">查询</el-button>
      <el-button v-permission="'article:add'" type="primary" @click="$router.push('/merchant/articles/edit')">发布</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="title" label="标题" min-width="180" />
      <el-table-column prop="audit_status" label="审核" width="90">
        <template #default="{ row }">
          <el-tag :type="auditType(row.audit_status)" size="small">{{ row.audit_status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column prop="schedule_publish_at" label="定时" width="160" />
      <el-table-column prop="reject_reason" label="驳回理由" min-width="140" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="$router.push('/merchant/articles/edit/' + row.id)">
            {{ row.audit_status === 'pending' ? '编辑' : '查看' }}
          </el-button>
          <el-button
            v-if="row.audit_status === 'pending'"
            v-permission="'article:delete'"
            link type="danger"
            @click="remove(row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pager" layout="prev, pager, next, total" :total="total" v-model:current-page="page" :page-size="20" @current-change="load" />
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listMyArticles, deleteMyArticle } from '../../api/merchant-article'

const list = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const q = reactive({ title: '', audit_status: '' })

function auditType(s) {
  if (s === 'approved') return 'success'
  if (s === 'rejected') return 'danger'
  return 'warning'
}

async function load() {
  loading.value = true
  try {
    const res = await listMyArticles({
      page: page.value, page_size: 20,
      title: q.title || undefined,
      audit_status: q.audit_status || undefined,
    })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function remove(row) {
  await ElMessageBox.confirm('删除待审文章？')
  await deleteMyArticle(row.id)
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.pager { margin-top: 12px; }
</style>
