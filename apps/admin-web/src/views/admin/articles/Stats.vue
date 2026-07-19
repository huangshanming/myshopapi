<template>
  <div>
    <div class="toolbar"><h2>文章数据统计</h2><el-button @click="load">刷新</el-button></div>
    <el-row :gutter="16" v-loading="loading">
      <el-col :span="6"><el-card shadow="never"><div class="n">{{ data.audit_pending || 0 }}</div><div class="l">待审核</div></el-card></el-col>
      <el-col :span="6"><el-card shadow="never"><div class="n">{{ data.view_sum || 0 }}</div><div class="l">总浏览</div></el-card></el-col>
      <el-col :span="6"><el-card shadow="never"><div class="n">{{ data.like_sum || 0 }}</div><div class="l">总点赞</div></el-card></el-col>
      <el-col :span="6"><el-card shadow="never"><div class="n">{{ published }}</div><div class="l">已发布</div></el-card></el-col>
    </el-row>
    <h3 style="margin-top: 24px">状态分布</h3>
    <el-table :data="statusRows" stripe style="max-width: 480px">
      <el-table-column prop="status" label="状态" />
      <el-table-column prop="cnt" label="数量" />
    </el-table>
    <h3 style="margin-top: 24px">近 7 日新建</h3>
    <el-table :data="data.recent_7d || []" stripe style="max-width: 480px">
      <el-table-column prop="Day" label="日期">
        <template #default="{ row }">{{ row.Day || row.day }}</template>
      </el-table-column>
      <el-table-column prop="Cnt" label="数量">
        <template #default="{ row }">{{ row.Cnt ?? row.cnt }}</template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { articleStats } from '../../../api/admin-article'

const loading = ref(false)
const data = ref({})

const published = computed(() => data.value.by_status?.published || 0)
const statusRows = computed(() => {
  const m = data.value.by_status || {}
  return Object.keys(m).map((k) => ({ status: k, cnt: m[k] }))
})

async function load() {
  loading.value = true
  try {
    const res = await articleStats()
    data.value = res || {}
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 16px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.n { font-size: 28px; font-weight: 700; }
.l { color: #64748b; margin-top: 4px; }
</style>
