<template>
  <div>
    <div class="toolbar">
      <h2>用户管理</h2>
      <el-input v-model="mobile" placeholder="手机号" clearable style="width: 200px" @clear="load" @keyup.enter="load" />
      <el-button @click="load">查询</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="mobile" label="手机号" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="role" label="角色" width="140" />
      <el-table-column prop="status" label="状态" width="80" />
      <el-table-column label="操作" width="520">
        <template #default="{ row }">
          <el-button v-permission="'system:user:edit'" link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button v-permission="'system:user:wallet'" link type="warning" @click="openWallet(row)">钱包</el-button>
          <el-button v-permission="'system:user:list'" link @click="openAddresses(row)">地址</el-button>
          <el-button v-permission="'system:user:list'" link @click="openFavorites(row)">收藏</el-button>
          <el-button v-permission="'system:user:reset'" link @click="openReset(row)">重置密码</el-button>
          <el-button v-permission="'system:user:list'" link type="success" @click="openToken(row)">生成Token</el-button>
          <el-button
            v-permission="'system:user:status'"
            size="small"
            :type="row.status === 1 ? 'warning' : 'success'"
            @click="toggle(row)"
          >{{ row.status === 1 ? '禁用' : '启用' }}</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pager"
      layout="prev, pager, next, total"
      :total="total"
      v-model:current-page="page"
      :page-size="pageSize"
      @current-change="load"
    />

    <el-dialog v-model="editVisible" title="编辑用户" width="460px">
      <el-form label-width="80px">
        <el-form-item label="手机号"><el-input v-model="editForm.mobile" maxlength="11" /></el-form-item>
        <el-form-item label="昵称"><el-input v-model="editForm.nickname" /></el-form-item>
        <el-form-item label="头像URL"><el-input v-model="editForm.avatar" /></el-form-item>
        <el-form-item label="性别">
          <el-select v-model="editForm.gender" style="width: 100%">
            <el-option label="未知" :value="0" />
            <el-option label="男" :value="1" />
            <el-option label="女" :value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="pwdVisible" title="重置密码" width="420px">
      <el-input v-model="newPassword" type="password" show-password placeholder="新密码" />
      <template #footer>
        <el-button @click="pwdVisible = false">取消</el-button>
        <el-button type="primary" @click="savePwd">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="walletVisible" :title="`用户钱包 #${walletUserId}`" width="560px" destroy-on-close>
      <div v-loading="walletLoading" class="wallet-box">
        <div class="wallet-nums">
          <div><span>可用余额</span><b>¥{{ wallet.balance ?? 0 }}</b></div>
          <div><span>冻结余额</span><b>¥{{ wallet.frozen_balance ?? 0 }}</b></div>
        </div>
        <el-form label-width="90px" class="adjust-form">
          <el-form-item label="调整项目">
            <el-radio-group v-model="adjustField">
              <el-radio-button value="balance">可用余额</el-radio-button>
              <el-radio-button value="frozen_balance">冻结余额</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="调账金额">
            <el-input-number v-model="adjustAmount" :precision="2" :step="10" />
            <span class="tip">正数增加，负数减少</span>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="adjustRemark" placeholder="调账说明" />
          </el-form-item>
          <el-button type="primary" @click="doAdjust">确认调账</el-button>
        </el-form>
        <h4>近期流水</h4>
        <el-table :data="walletLogs" size="small" max-height="240">
          <el-table-column prop="created_at" label="时间" width="150" />
          <el-table-column prop="change_type" label="类型" width="110" />
          <el-table-column prop="amount" label="金额" width="90" />
          <el-table-column prop="balance_after" label="余额后" width="80" />
          <el-table-column prop="frozen_after" label="冻结后" width="80" />
          <el-table-column prop="remark" label="备注" min-width="100" />
        </el-table>
      </div>
    </el-dialog>

    <el-dialog v-model="addrVisible" :title="`用户地址 #${addrUserId}`" width="720px" destroy-on-close>
      <el-table v-loading="addrLoading" :data="addrList" size="small" max-height="420">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="默认" width="70">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="warning" size="small">默认</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="receiver_name" label="收货人" width="100" />
        <el-table-column prop="receiver_phone" label="手机" width="120" />
        <el-table-column label="地址" min-width="220">
          <template #default="{ row }">
            {{ row.province }}{{ row.city }}{{ row.district }}{{ row.detail }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160" />
      </el-table>
      <el-empty v-if="!addrLoading && !addrList.length" description="暂无收货地址" />
    </el-dialog>

    <el-dialog v-model="favVisible" :title="`用户收藏 #${favUserId}`" width="720px" destroy-on-close>
      <el-table v-loading="favLoading" :data="favList" size="small" max-height="420">
        <el-table-column prop="product_id" label="商品ID" width="90" />
        <el-table-column label="封面" width="72">
          <template #default="{ row }">
            <el-image v-if="row.main_image" :src="row.main_image" style="width: 40px; height: 40px" fit="cover" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="160" />
        <el-table-column prop="sale_price" label="售价" width="90" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag v-if="row.invalid" type="info" size="small">失效</el-tag>
            <span v-else>{{ row.status }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="收藏时间" width="160" />
      </el-table>
      <el-empty v-if="!favLoading && !favList.length" description="暂无收藏" />
    </el-dialog>

    <el-dialog v-model="tokenVisible" title="用户 Token（并发压测可复制）" width="640px">
      <el-descriptions :column="1" border size="small" class="token-meta">
        <el-descriptions-item label="用户">{{ tokenInfo.nickname }}（{{ tokenInfo.mobile }}）</el-descriptions-item>
        <el-descriptions-item label="user_id">{{ tokenInfo.user_id }}</el-descriptions-item>
        <el-descriptions-item label="role">{{ tokenInfo.role }}</el-descriptions-item>
        <el-descriptions-item v-if="tokenInfo.shop_id" label="shop_id">{{ tokenInfo.shop_id }}</el-descriptions-item>
        <el-descriptions-item label="有效期">{{ tokenInfo.expire_hours }} 小时</el-descriptions-item>
      </el-descriptions>
      <el-input
        v-model="tokenInfo.token"
        type="textarea"
        :rows="5"
        readonly
        class="token-box"
      />
      <template #footer>
        <el-button @click="tokenVisible = false">关闭</el-button>
        <el-button type="primary" :loading="tokenLoading" @click="copyToken">复制 Token</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  adjustUserWallet, fetchUserAddresses, fetchUserFavorites, fetchUserWalletLogs, fetchUsers, generateUserToken, getUserWallet,
  resetUserPassword, setUserStatus, updateUser,
} from '../../../api/system'

const list = ref([])
const loading = ref(false)
const mobile = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)
const editVisible = ref(false)
const pwdVisible = ref(false)
const tokenVisible = ref(false)
const tokenLoading = ref(false)
const editForm = ref({ mobile: '', nickname: '', avatar: '', gender: 0 })
const currentId = ref(0)
const newPassword = ref('123456')
const tokenInfo = reactive({
  token: '',
  user_id: 0,
  mobile: '',
  nickname: '',
  role: '',
  shop_id: 0,
  expire_hours: 24,
})
const walletVisible = ref(false)
const walletLoading = ref(false)
const walletUserId = ref(0)
const wallet = ref({})
const walletLogs = ref([])
const adjustAmount = ref(100)
const adjustRemark = ref('')
const adjustField = ref('balance')
const addrVisible = ref(false)
const addrLoading = ref(false)
const addrUserId = ref(0)
const addrList = ref([])
const favVisible = ref(false)
const favLoading = ref(false)
const favUserId = ref(0)
const favList = ref([])

async function load() {
  loading.value = true
  try {
    const res = await fetchUsers({ page: page.value, page_size: pageSize, mobile: mobile.value || undefined })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openEdit(row) {
  currentId.value = row.id
  editForm.value = {
    mobile: row.mobile || '',
    nickname: row.nickname || '',
    avatar: row.avatar || '',
    gender: row.gender ?? 0,
  }
  editVisible.value = true
}

async function saveEdit() {
  try {
    await updateUser(currentId.value, editForm.value)
    ElMessage.success('已更新')
    editVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

function openReset(row) {
  currentId.value = row.id
  newPassword.value = '123456'
  pwdVisible.value = true
}

async function savePwd() {
  try {
    await resetUserPassword(currentId.value, newPassword.value)
    ElMessage.success('已重置')
    pwdVisible.value = false
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function openToken(row) {
  tokenLoading.value = true
  tokenVisible.value = true
  tokenInfo.token = ''
  try {
    const res = await generateUserToken(row.id)
    const d = res || {}
    tokenInfo.token = d.token || ''
    tokenInfo.user_id = d.user_id
    tokenInfo.mobile = d.mobile || ''
    tokenInfo.nickname = d.nickname || ''
    tokenInfo.role = d.role || ''
    tokenInfo.shop_id = d.shop_id || 0
    tokenInfo.expire_hours = d.expire_hours || 24
  } catch (e) {
    ElMessage.error(e.message)
    tokenVisible.value = false
  } finally {
    tokenLoading.value = false
  }
}

async function copyToken() {
  if (!tokenInfo.token) {
    ElMessage.warning('暂无 Token')
    return
  }
  try {
    await navigator.clipboard.writeText(tokenInfo.token)
    ElMessage.success('已复制到剪贴板')
  } catch {
    // fallback
    const ta = document.createElement('textarea')
    ta.value = tokenInfo.token
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
    ElMessage.success('已复制到剪贴板')
  }
}

async function toggle(row) {
  const status = row.status === 1 ? 0 : 1
  try {
    await setUserStatus(row.id, status)
    ElMessage.success('已更新')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

function openWallet(row) {
  walletUserId.value = row.id
  walletVisible.value = true
  adjustAmount.value = 100
  adjustRemark.value = ''
  adjustField.value = 'balance'
  loadWallet()
}

async function loadWallet() {
  walletLoading.value = true
  try {
    const [w, logs] = await Promise.all([
      getUserWallet(walletUserId.value),
      fetchUserWalletLogs(walletUserId.value, { page: 1, page_size: 20 }),
    ])
    wallet.value = w || {}
    walletLogs.value = logs?.list || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    walletLoading.value = false
  }
}

async function doAdjust() {
  try {
    await adjustUserWallet(walletUserId.value, {
      field: adjustField.value,
      amount: adjustAmount.value,
      remark: adjustRemark.value,
    })
    ElMessage.success('调账成功')
    loadWallet()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

function openAddresses(row) {
  addrUserId.value = row.id
  addrVisible.value = true
  loadAddresses()
}

async function loadAddresses() {
  addrLoading.value = true
  try {
    const res = await fetchUserAddresses(addrUserId.value)
    addrList.value = res || []
  } catch (e) {
    ElMessage.error(e.message)
    addrList.value = []
  } finally {
    addrLoading.value = false
  }
}

function openFavorites(row) {
  favUserId.value = row.id
  favVisible.value = true
  loadFavorites()
}

async function loadFavorites() {
  favLoading.value = true
  try {
    const res = await fetchUserFavorites(favUserId.value, { page: 1, page_size: 100 })
    favList.value = res?.list || []
  } catch (e) {
    ElMessage.error(e.message)
    favList.value = []
  } finally {
    favLoading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin-right: auto; }
.pager { margin-top: 12px; }
.token-meta { margin-bottom: 12px; }
.token-box { font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; }
.wallet-nums { display: flex; gap: 24px; margin-bottom: 16px; }
.wallet-nums div { display: flex; flex-direction: column; gap: 4px; }
.wallet-nums span { color: #64748b; font-size: 12px; }
.wallet-nums b { font-size: 18px; }
.adjust-form { margin-bottom: 16px; }
.tip { margin-left: 8px; color: #94a3b8; font-size: 12px; }
</style>
