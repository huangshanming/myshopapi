<template>
  <div>
    <div class="toolbar">
      <h2>文章回收站</h2>
      <el-button @click="load">刷新</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="shop_id" label="商家" width="80" />
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column prop="deleted_at" label="删除时间" width="160" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button link type="primary" @click="restore(row)">恢复</el-button>
          <el-button link type="danger" @click="purge(row)">永久删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pager" layout="prev, pager, next, total" :total="total" v-model:current-page="page" :page-size="20" @current-change="load" />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listArticleRecycle, restoreArticle, permanentDeleteArticle } from '../../../api/admin-article'

const list = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await listArticleRecycle({ page: page.value, page_size: 20 })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

async function restore(row) {
  await restoreArticle(row.id)
  ElMessage.success('已恢复为下架')
  load()
}

async function purge(row) {
  await ElMessageBox.confirm('永久删除将级联删除评论与配图，不可恢复')
  await permanentDeleteArticle(row.id)
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.pager { margin-top: 12px; }
</style>
