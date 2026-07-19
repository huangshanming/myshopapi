<template>
  <div>
    <div class="toolbar">
      <h2>用户消息</h2>
    </div>
    <p class="tip">向全部用户或指定用户发送站内信公告；订单状态变更消息由系统自动推送。</p>

    <el-card class="mb">
      <template #header>发送公告</template>
      <el-form label-width="100px" style="max-width:640px">
        <el-form-item label="标题" required>
          <el-input v-model="form.title" maxlength="100" />
        </el-form-item>
        <el-form-item label="正文">
          <el-input v-model="form.content" type="textarea" :rows="4" maxlength="1000" />
        </el-form-item>
        <el-form-item label="发送目标">
          <el-radio-group v-model="form.target">
            <el-radio-button value="all">全部用户</el-radio-button>
            <el-radio-button value="users">指定用户</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.target === 'users'" label="选择用户" required>
          <el-select
            v-model="selectedUserIds"
            multiple
            filterable
            remote
            clearable
            collapse-tags
            collapse-tags-tooltip
            :remote-method="searchUsers"
            :loading="userLoading"
            style="width:100%"
            placeholder="搜索手机号选择用户"
          >
            <el-option
              v-for="u in userOptions"
              :key="u.id"
              :label="userLabel(u)"
              :value="u.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关联订单">
          <el-select
            v-model="form.link_id"
            filterable
            remote
            clearable
            :remote-method="searchOrders"
            :loading="orderLoading"
            style="width:100%"
            placeholder="搜索订单号（可选，用于跳转订单详情）"
          >
            <el-option
              v-for="o in orderOptions"
              :key="o.id"
              :label="orderLabel(o)"
              :value="o.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button v-permission="'business:message:send'" type="primary" :loading="sending" @click="send">发送</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <template #header>最近发送</template>
      <el-table :data="sends" stripe>
        <el-table-column prop="title" label="标题" min-width="160" />
        <el-table-column label="正文" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">{{ row.content || '—' }}</template>
        </el-table-column>
        <el-table-column label="发送范围" width="110">
          <template #default="{ row }">{{ targetLabel(row.target) }}</template>
        </el-table-column>
        <el-table-column label="发送结果" width="130">
          <template #default="{ row }">成功 {{ row.success_count }} / {{ row.user_count }} 人</template>
        </el-table-column>
        <el-table-column label="跳转" width="100">
          <template #default="{ row }">{{ linkLabel(row) }}</template>
        </el-table-column>
        <el-table-column prop="created_at" label="发送时间" width="180" />
        <el-table-column label="操作" width="90" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-drawer v-model="detailVisible" title="发送详情" size="520px" destroy-on-close>
      <template v-if="detailBatch">
        <el-descriptions :column="1" border size="small" class="mb">
          <el-descriptions-item label="标题">{{ detailBatch.title }}</el-descriptions-item>
          <el-descriptions-item label="正文">{{ detailBatch.content || '—' }}</el-descriptions-item>
          <el-descriptions-item label="范围">{{ targetLabel(detailBatch.target) }}</el-descriptions-item>
          <el-descriptions-item label="结果">
            成功 {{ detailBatch.success_count }} / {{ detailBatch.user_count }} 人
          </el-descriptions-item>
          <el-descriptions-item label="跳转">{{ linkLabel(detailBatch) }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ detailBatch.created_at }}</el-descriptions-item>
        </el-descriptions>
        <h4 class="sec">接收用户</h4>
        <el-table :data="recipients" v-loading="detailLoading" stripe size="small">
          <el-table-column label="昵称" min-width="120">
            <template #default="{ row }">{{ row.nickname || '用户' }}</template>
          </el-table-column>
          <el-table-column prop="mobile" label="手机号" min-width="120" />
          <el-table-column label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.is_read ? 'info' : 'warning'" size="small">
                {{ row.is_read ? '已读' : '未读' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          class="pager"
          layout="prev, pager, next, total"
          :total="recipientTotal"
          v-model:current-page="recipientPage"
          :page-size="recipientPageSize"
          @current-change="loadRecipients"
        />
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listAdminNotificationRecipients, listAdminNotificationSends, sendAdminNotification,
} from '../../api/message'
import { fetchUsers } from '../../api/system'
import { listOrders } from '../../api/order'

const form = reactive({
  title: '', content: '', target: 'all', link_id: undefined,
})
const selectedUserIds = ref([])
const userOptions = ref([])
const userLoading = ref(false)
const orderOptions = ref([])
const orderLoading = ref(false)
const sending = ref(false)
const sends = ref([])

const detailVisible = ref(false)
const detailLoading = ref(false)
const detailBatch = ref(null)
const recipients = ref([])
const recipientTotal = ref(0)
const recipientPage = ref(1)
const recipientPageSize = 20
const detailBatchId = ref(0)

function userLabel(u) {
  const name = u.nickname || '用户'
  return u.mobile ? `${name} · ${u.mobile}` : name
}

function orderLabel(o) {
  const parts = [o.order_no]
  if (o.user_name) parts.push(o.user_name)
  if (o.total_amount != null) parts.push(`¥${o.total_amount}`)
  return parts.join(' · ')
}

function targetLabel(t) {
  return t === 'all' ? '全部用户' : '指定用户'
}

function linkLabel(row) {
  if (row.link_type === 'order' && row.link_id) return '订单详情'
  return '无'
}

async function searchUsers(kw) {
  const q = String(kw || '').trim()
  if (!q) {
    userOptions.value = []
    return
  }
  userLoading.value = true
  try {
    const res = await fetchUsers({ page: 1, page_size: 20, mobile: q })
    const list = res?.list || []
    const selected = userOptions.value.filter((u) => selectedUserIds.value.includes(u.id))
    const map = new Map(selected.map((u) => [u.id, u]))
    list.forEach((u) => map.set(u.id, u))
    userOptions.value = [...map.values()]
  } catch {
    userOptions.value = []
  } finally {
    userLoading.value = false
  }
}

async function searchOrders(kw) {
  const q = String(kw || '').trim()
  orderLoading.value = true
  try {
    const res = await listOrders('admin', {
      page: 1,
      page_size: 20,
      order_no: q || undefined,
    })
    orderOptions.value = res?.list || []
  } catch {
    orderOptions.value = []
  } finally {
    orderLoading.value = false
  }
}

async function reload() {
  const res = await listAdminNotificationSends({ page: 1, page_size: 20 })
  sends.value = res?.list || []
}

async function loadRecipients() {
  if (!detailBatchId.value) return
  detailLoading.value = true
  try {
    const res = await listAdminNotificationRecipients(detailBatchId.value, {
      page: recipientPage.value,
      page_size: recipientPageSize,
    })
    detailBatch.value = res?.batch || detailBatch.value
    recipients.value = res?.list || []
    recipientTotal.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
    recipients.value = []
  } finally {
    detailLoading.value = false
  }
}

async function openDetail(row) {
  detailBatchId.value = row.id
  detailBatch.value = row
  recipientPage.value = 1
  detailVisible.value = true
  await loadRecipients()
}

async function send() {
  if (!form.title.trim()) {
    ElMessage.warning('请填写标题')
    return
  }
  const linkId = form.link_id || 0
  const payload = {
    title: form.title.trim(),
    content: form.content.trim(),
    target: form.target,
    link_type: linkId > 0 ? 'order' : 'none',
    link_id: linkId,
  }
  if (form.target === 'users') {
    payload.user_ids = [...selectedUserIds.value]
    if (!payload.user_ids.length) {
      ElMessage.warning('请选择用户')
      return
    }
  }
  try {
    await ElMessageBox.confirm(
      form.target === 'all' ? '确认向全部用户发送？' : `确认向已选 ${payload.user_ids.length} 人发送？`,
      '发送确认',
    )
  } catch { return }
  sending.value = true
  try {
    const res = await sendAdminNotification(payload)
    ElMessage.success(`已发送 ${res.success_count} 人`)
    form.title = ''
    form.content = ''
    form.link_id = undefined
    selectedUserIds.value = []
    userOptions.value = []
    reload()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    sending.value = false
  }
}

onMounted(reload)
</script>

<style scoped>
.toolbar { margin-bottom: 8px; }
.tip { color: #64748b; font-size: 13px; margin: 0 0 16px; }
.mb { margin-bottom: 16px; }
.sec { margin: 16px 0 8px; font-size: 14px; }
.pager { margin-top: 12px; justify-content: flex-end; }
</style>
