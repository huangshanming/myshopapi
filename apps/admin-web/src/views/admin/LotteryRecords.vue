<template>
  <div>
    <div class="toolbar">
      <h2>抽奖记录</h2>
      <el-input
        v-model="q.activity_id"
        clearable
        placeholder="活动 ID"
        style="width: 120px"
        @keyup.enter="search"
      />
      <el-select v-model="q.prize_type" clearable placeholder="奖品类型" style="width: 140px" @change="search">
        <el-option value="physical" label="实物" />
        <el-option value="points" label="积分" />
        <el-option value="thanks" label="谢谢参与" />
      </el-select>
      <el-button type="primary" @click="search">查询</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="activity_id" label="活动" width="80" />
      <el-table-column label="用户" width="100">
        <template #default="{ row }">#{{ row.user_id }}</template>
      </el-table-column>
      <el-table-column prop="slot" label="格" width="50" />
      <el-table-column prop="prize_name" label="奖品" min-width="140" />
      <el-table-column label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="typeTag(row.prize_type)" size="small">{{ typeLabel(row.prize_type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="积分" width="80">
        <template #default="{ row }">
          <span v-if="row.prize_type === 'points'">+{{ row.points_amount }}</span>
          <span v-else class="sub">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="cost_points" label="消耗" width="70" />
      <el-table-column label="履约" width="110">
        <template #default="{ row }">
          <template v-if="row.prize_type === 'physical'">
            <el-tag :type="fulfillType(row.fulfill_status)" size="small">{{ fulfillLabel(row.fulfill_status) }}</el-tag>
          </template>
          <span v-else class="sub">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="170" />
    </el-table>

    <el-pagination
      class="pager"
      layout="prev,pager,next,total"
      :total="total"
      v-model:current-page="page"
      :page-size="20"
      @current-change="load"
    />
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { listLotteryRecords } from '../../api/lottery'

const route = useRoute()

const TYPE_OPTIONS = [
  { value: 'points', label: '积分' },
  { value: 'thanks', label: '谢谢参与' },
  { value: 'physical', label: '实物' },
]
const FULFILL_OPTIONS = [
  { value: 'need_address', label: '待填地址' },
  { value: 'pending', label: '待发货' },
  { value: 'shipped', label: '已发货' },
  { value: 'none', label: '-' },
]

const list = ref([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)
const q = reactive({
  activity_id: route.query.activity_id ? String(route.query.activity_id) : '',
  prize_type: route.query.prize_type === 'physical' ? 'physical' : '',
})

function typeLabel(t) {
  return TYPE_OPTIONS.find((o) => o.value === t)?.label || t || '-'
}
function typeTag(t) {
  if (t === 'physical') return 'warning'
  if (t === 'points') return 'success'
  return 'info'
}
function fulfillLabel(s) {
  return FULFILL_OPTIONS.find((o) => o.value === s)?.label || s || '-'
}
function fulfillType(s) {
  if (s === 'need_address') return 'info'
  if (s === 'pending') return 'warning'
  if (s === 'shipped') return 'success'
  return 'info'
}

function search() {
  page.value = 1
  load()
}

async function load() {
  loading.value = true
  try {
    const aid = Number(q.activity_id)
    const res = await listLotteryRecords({
      page: page.value,
      page_size: 20,
      activity_id: aid > 0 ? aid : undefined,
      prize_type: q.prize_type || undefined,
    })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 10px; align-items: center; margin-bottom: 12px; flex-wrap: wrap; }
.pager { margin-top: 16px; justify-content: flex-end; }
.sub { color: #94a3b8; font-size: 12px; }
</style>
