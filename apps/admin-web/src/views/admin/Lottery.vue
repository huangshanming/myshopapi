<template>
  <div class="page">
    <div class="toolbar">
      <h2>九宫格抽奖</h2>
      <el-button v-permission="'marketing:lottery:list'" type="primary" @click="openCreate">新建活动</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="title" label="标题" min-width="140" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="cost_points" label="消耗积分" width="100" />
      <el-table-column prop="daily_limit" label="每日次数" width="100" />
      <el-table-column prop="start_at" label="开始" width="160" />
      <el-table-column prop="end_at" label="结束" width="160" />
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'marketing:lottery:list'" link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button v-permission="'marketing:lottery:list'" link type="primary" @click="openPrizes(row)">九宫格</el-button>
          <el-button v-permission="'marketing:lottery:list'" link @click="goRecords(row)">记录</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="formVisible" :title="form.id ? '编辑活动' : '新建活动'" width="520px" destroy-on-close>
      <el-form label-width="110px">
        <el-form-item label="标题" required>
          <el-input v-model="form.title" maxlength="100" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width: 100%">
            <el-option :value="0" label="草稿" />
            <el-option :value="1" label="上线" />
            <el-option :value="2" label="下线" />
          </el-select>
        </el-form-item>
        <el-form-item label="消耗积分">
          <el-input-number v-model="form.cost_points" :min="0" />
          <span class="hint">0 = 免费</span>
        </el-form-item>
        <el-form-item label="每日次数">
          <el-input-number v-model="form.daily_limit" :min="0" />
          <span class="hint">0 = 不限</span>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-input v-model="form.start_at" placeholder="2026-01-01 00:00:00" />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-input v-model="form.end_at" placeholder="2026-12-31 23:59:59" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveForm">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="prizeVisible" title="配置九宫格奖品" width="1100px" destroy-on-close>
      <p class="tip">必须 9 格（slot 0-8）。权重合计：{{ weightSum }}。库存 -1 表示无限；「强控库存」供后续业务使用。</p>
      <el-table :data="prizes" size="small" border>
        <el-table-column prop="slot" label="格" width="50" />
        <el-table-column label="名称" min-width="120">
          <template #default="{ row }"><el-input v-model="row.name" /></template>
        </el-table-column>
        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            <el-select v-model="row.prize_type">
              <el-option value="points" label="积分" />
              <el-option value="physical" label="实物" />
              <el-option value="thanks" label="谢谢参与" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="封面" width="110">
          <template #default="{ row }">
            <SingleImageField v-model="row.cover_url" :upload-fn="uploadLotteryImage" />
          </template>
        </el-table-column>
        <el-table-column label="积分数" width="110">
          <template #default="{ row }">
            <el-input-number v-model="row.points_amount" :min="0" :disabled="row.prize_type !== 'points'" controls-position="right" />
          </template>
        </el-table-column>
        <el-table-column label="权重" width="100">
          <template #default="{ row }"><el-input-number v-model="row.weight" :min="0" controls-position="right" /></template>
        </el-table-column>
        <el-table-column label="库存" width="110">
          <template #default="{ row }"><el-input-number v-model="row.stock" :min="-1" controls-position="right" /></template>
        </el-table-column>
        <el-table-column label="强控库存" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.stock_strict" :active-value="1" :inactive-value="0" />
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="prizeVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingPrizes" @click="savePrizes">保存九宫格</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import SingleImageField from '../../components/admin/SingleImageField.vue'
import {
  createLotteryActivity,
  getLotteryActivity,
  listLotteryActivities,
  saveLotteryPrizes,
  updateLotteryActivity,
  uploadLotteryImage,
} from '../../api/lottery'

const router = useRouter()

const list = ref([])
const loading = ref(false)
const formVisible = ref(false)
const saving = ref(false)
const form = reactive({
  id: 0,
  title: '',
  status: 0,
  cost_points: 10,
  daily_limit: 3,
  start_at: '',
  end_at: '',
})

const prizeVisible = ref(false)
const savingPrizes = ref(false)
const prizeActivityId = ref(0)
const prizes = ref([])

const weightSum = computed(() => prizes.value.reduce((s, p) => s + (Number(p.weight) || 0), 0))

function statusText(s) {
  return { 0: '草稿', 1: '上线', 2: '下线' }[s] || String(s)
}
function statusType(s) {
  return { 0: 'info', 1: 'success', 2: 'warning' }[s] || 'info'
}

async function load() {
  loading.value = true
  try {
    const res = await listLotteryActivities({ page: 1, page_size: 50 })
    list.value = res?.list || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, {
    id: 0, title: '', status: 0, cost_points: 10, daily_limit: 3, start_at: '', end_at: '',
  })
  formVisible.value = true
}

async function openEdit(row) {
  try {
    const res = await getLotteryActivity(row.id)
    const d = res?.data || res
    Object.assign(form, {
      id: d.id,
      title: d.title,
      status: d.status,
      cost_points: d.cost_points,
      daily_limit: d.daily_limit,
      start_at: d.start_at || '',
      end_at: d.end_at || '',
    })
    formVisible.value = true
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function saveForm() {
  saving.value = true
  try {
    const payload = {
      title: form.title,
      status: form.status,
      cost_points: form.cost_points,
      daily_limit: form.daily_limit,
      start_at: form.start_at,
      end_at: form.end_at,
    }
    if (form.id) await updateLotteryActivity(form.id, payload)
    else await createLotteryActivity(payload)
    ElMessage.success('已保存')
    formVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

function emptyPrizes() {
  return Array.from({ length: 9 }, (_, i) => ({
    slot: i,
    name: i % 2 === 1 ? '谢谢参与' : `${(i + 1) * 10}积分`,
    cover_url: '',
    prize_type: i % 2 === 1 ? 'thanks' : 'points',
    points_amount: i % 2 === 1 ? 0 : (i + 1) * 10,
    weight: i % 2 === 1 ? 40 : 10,
    stock: -1,
    stock_strict: 0,
  }))
}

async function openPrizes(row) {
  prizeActivityId.value = row.id
  try {
    const res = await getLotteryActivity(row.id)
    const d = res?.data || res
    const arr = d?.prizes || []
    if (arr.length === 9) {
      prizes.value = arr.map((p) => ({
        slot: p.slot,
        name: p.name,
        cover_url: p.cover_url || '',
        prize_type: p.prize_type,
        points_amount: p.points_amount || 0,
        weight: p.weight || 0,
        stock: p.stock ?? -1,
        stock_strict: p.stock_strict ? 1 : 0,
      }))
    } else {
      prizes.value = emptyPrizes()
    }
    prizeVisible.value = true
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function savePrizes() {
  savingPrizes.value = true
  try {
    await saveLotteryPrizes(prizeActivityId.value, prizes.value)
    ElMessage.success('九宫格已保存')
    prizeVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    savingPrizes.value = false
  }
}

async function goRecords(row) {
  router.push({ path: '/admin/lottery-records', query: { activity_id: String(row.id) } })
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.hint { margin-left: 8px; color: #94a3b8; font-size: 12px; }
.tip { color: #64748b; font-size: 13px; margin-bottom: 10px; }
:deep(.single-img .preview .el-image),
:deep(.single-img .add) {
  width: 56px !important;
  height: 56px !important;
}
:deep(.single-img .add) { font-size: 22px; }
</style>
