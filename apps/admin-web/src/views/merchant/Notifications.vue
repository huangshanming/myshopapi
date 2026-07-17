<template>
  <div>
    <div class="toolbar">
      <h2>消息通知</h2>
      <el-select v-model="isRead" clearable placeholder="已读状态" style="width: 120px" @change="load">
        <el-option label="未读" :value="0" />
        <el-option label="已读" :value="1" />
      </el-select>
      <el-button @click="load">刷新</el-button>
      <el-button type="primary" @click="readAll">全部已读</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe @row-click="onRow">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="title" label="标题" min-width="160" />
      <el-table-column prop="content" label="内容" min-width="280" />
      <el-table-column prop="is_read" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_read ? 'info' : 'danger'" size="small">{{ row.is_read ? '已读' : '未读' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="160" />
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button v-if="!row.is_read" link type="primary" @click.stop="markOne(row)">标已读</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pager"
      layout="prev, pager, next, total"
      :total="total"
      v-model:current-page="page"
      :page-size="20"
      @current-change="load"
    />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  listNotifications, markNotificationRead, markAllNotificationsRead,
} from '../../api/merchant-notification'

const router = useRouter()
const list = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const isRead = ref(undefined)

async function load() {
  loading.value = true
  try {
    const res = await listNotifications({
      page: page.value,
      page_size: 20,
      is_read: isRead.value === undefined || isRead.value === null || isRead.value === '' ? undefined : isRead.value,
    })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function markOne(row) {
  await markNotificationRead(row.id)
  load()
}

async function readAll() {
  await markAllNotificationsRead()
  ElMessage.success('已全部标为已读')
  load()
}

async function onRow(row) {
  if (!row.is_read) {
    try {
      await markNotificationRead(row.id)
    } catch (_) { /* ignore */ }
  }
  if (row.link) {
    router.push(row.link)
  }
}

onMounted(load)

defineExpose({ load })
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.pager { margin-top: 12px; }
</style>
