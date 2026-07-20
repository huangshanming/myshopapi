<template>
  <div>
    <div class="toolbar">
      <h2>任务中心</h2>
    </div>
    <p class="tip">
      配置 C 端任务开关与奖励规则。发文奖励即「发布种草笔记」：可关闭、设每次积分、设一天最多领取次数（0=不限）。
      用户完成任务后需在任务中心手动领取积分。
    </p>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="code" label="编码" width="160" />
      <el-table-column prop="title" label="名称" min-width="120" />
      <el-table-column label="周期" width="90">
        <template #default="{ row }">{{ row.period === 'once' ? '一次性' : '每日' }}</template>
      </el-table-column>
      <el-table-column label="启用" width="90">
        <template #default="{ row }">
          <el-tag :type="row.enabled === 1 ? 'success' : 'info'" size="small">
            {{ row.enabled === 1 ? '开' : '关' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="reward_points" label="积分" width="80" />
      <el-table-column prop="target_count" label="目标次数" width="90" />
      <el-table-column prop="daily_limit" label="每日上限" width="90" />
      <el-table-column prop="sort" label="排序" width="70" />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'marketing:task:edit'" link type="primary" @click="openEdit(row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" title="编辑任务" width="520px" destroy-on-close>
      <el-form label-width="110px">
        <el-form-item label="名称">
          <el-input v-model="form.title" maxlength="100" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="form.description" type="textarea" :rows="2" maxlength="500" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="enabled" active-text="开" inactive-text="关" />
        </el-form-item>
        <el-form-item label="奖励积分">
          <el-input-number v-model="form.reward_points" :min="0" :max="999999" />
          <span class="hint">0 表示不发积分</span>
        </el-form-item>
        <el-form-item label="目标次数">
          <el-input-number v-model="form.target_count" :min="1" :max="999" />
        </el-form-item>
        <el-form-item v-if="form.period === 'daily'" label="每日上限">
          <el-input-number v-model="form.daily_limit" :min="0" :max="999" />
          <span class="hint">一天最多领取几次，0=不限</span>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item v-if="form.code === 'publish_article'" label="">
          <el-alert type="info" :closable="false" show-icon
            title="发文任务：用户笔记审核通过后计进度；可在此关闭奖励或限制每日次数。" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { listAdminTasks, updateAdminTask } from '../../api/task-center'

const list = ref([])
const loading = ref(false)
const visible = ref(false)
const saving = ref(false)
const form = reactive({
  id: 0, code: '', title: '', description: '', period: 'daily',
  enabled: 1, reward_points: 0, target_count: 1, daily_limit: 0, sort: 0,
})
const enabled = computed({
  get: () => form.enabled === 1,
  set: (v) => { form.enabled = v ? 1 : 0 },
})

async function load() {
  loading.value = true
  try {
    const res = await listAdminTasks()
    list.value = res?.list || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openEdit(row) {
  Object.assign(form, {
    id: row.id, code: row.code, title: row.title, description: row.description,
    period: row.period, enabled: row.enabled, reward_points: row.reward_points,
    target_count: row.target_count, daily_limit: row.daily_limit, sort: row.sort,
  })
  visible.value = true
}

async function save() {
  saving.value = true
  try {
    await updateAdminTask(form.id, {
      title: form.title,
      description: form.description,
      enabled: form.enabled,
      reward_points: form.reward_points,
      target_count: form.target_count,
      daily_limit: form.daily_limit,
      sort: form.sort,
    })
    ElMessage.success('已保存')
    visible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; }
.tip { color: #64748b; font-size: 13px; margin: 0 0 16px; max-width: 720px; line-height: 1.5; }
.hint { margin-left: 8px; color: #94a3b8; font-size: 12px; }
</style>
