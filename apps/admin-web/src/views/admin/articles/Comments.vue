<template>
  <div>
    <div class="toolbar">
      <h2>评论管理</h2>
      <el-input v-model="q.shop_id" placeholder="商家ID" clearable style="width: 110px" />
      <el-input v-model="q.article_id" placeholder="文章ID" clearable style="width: 110px" />
      <el-select v-model="q.status" clearable placeholder="状态" style="width: 120px">
        <el-option label="可见" value="visible" />
        <el-option label="隐藏" value="hidden" />
      </el-select>
      <el-button @click="load">查询</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="article_id" label="文章" width="80" />
      <el-table-column prop="shop_id" label="商家" width="80" />
      <el-table-column prop="user_id" label="用户" width="80" />
      <el-table-column prop="content" label="内容" min-width="220" />
      <el-table-column prop="status" label="状态" width="90" />
      <el-table-column prop="created_at" label="时间" width="160" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button link @click="setStatus(row, row.status === 'hidden' ? 'visible' : 'hidden')">
            {{ row.status === 'hidden' ? '显示' : '隐藏' }}
          </el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pager" layout="prev, pager, next, total" :total="total" v-model:current-page="page" :page-size="20" @current-change="load" />
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listArticleComments, patchArticleComment, deleteArticleComment } from '../../../api/admin-article'

const list = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const q = reactive({ shop_id: '', article_id: '', status: '' })

async function load() {
  loading.value = true
  try {
    const res = await listArticleComments({
      page: page.value, page_size: 20,
      shop_id: q.shop_id || undefined,
      article_id: q.article_id || undefined,
      status: q.status || undefined,
    })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

async function setStatus(row, status) {
  await patchArticleComment(row.id, status)
  load()
}

async function remove(row) {
  await ElMessageBox.confirm('删除该评论？')
  await deleteArticleComment(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.pager { margin-top: 12px; }
</style>
