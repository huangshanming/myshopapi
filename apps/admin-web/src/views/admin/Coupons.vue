<template>
  <div>
    <div class="toolbar">
      <h2>优惠券</h2>
      <el-button v-permission="'marketing:coupon:edit'" type="primary" @click="openCreate">创建</el-button>
    </div>
    <el-form inline class="mb">
      <el-form-item label="状态">
        <el-select v-model="query.status" clearable style="width:140px" @change="reload">
          <el-option label="上架" value="on" />
          <el-option label="下架" value="off" />
          <el-option label="草稿" value="draft" />
        </el-select>
      </el-form-item>
      <el-form-item label="关键词">
        <el-input v-model="query.keyword" clearable @keyup.enter="reload" />
      </el-form-item>
      <el-button @click="reload">查询</el-button>
    </el-form>

    <el-table :data="list" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column label="类型" width="110">
        <template #default="{ row }">{{ couponTypeLabel(row.coupon_type) }}</template>
      </el-table-column>
      <el-table-column label="优惠" min-width="140">
        <template #default="{ row }">{{ benefitText(row) }}</template>
      </el-table-column>
      <el-table-column label="库存" width="100">
        <template #default="{ row }">{{ row.claimed_count }}/{{ row.total_count || '∞' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }">{{ displayStatusLabel(row.display_status || row.status) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="360" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link @click="doCopy(row)">复制</el-button>
          <el-button v-permission="'marketing:coupon:grant'" link @click="openGrant(row)">发放</el-button>
          <el-button link @click="openStats(row)">统计</el-button>
          <el-button link @click="openClaims(row)">领取</el-button>
          <el-button link type="danger" @click="doOff(row)">下架</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="mt" layout="total, prev, pager, next" :total="total" v-model:current-page="query.page" :page-size="20" @current-change="reload" />

    <el-dialog v-model="formVisible" :title="form.id ? '编辑优惠券' : '创建优惠券'" width="640px">
      <el-form label-width="110px">
        <el-form-item label="名称" required><el-input v-model="form.name" maxlength="100" /></el-form-item>
        <el-form-item label="类型" required>
          <el-select v-model="form.coupon_type" style="width:100%" @change="onType">
            <el-option v-for="t in COUPON_TYPES" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.coupon_type !== 'discount' && form.coupon_type !== 'no_threshold'" label="门槛">
          <el-input-number v-model="form.threshold_amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item v-if="form.coupon_type === 'discount'" label="折扣比例" required>
          <el-input-number v-model="form.discount_rate" :min="0.01" :max="0.99" :step="0.05" :precision="2" />
          <span class="tip">如 0.8 = 八折</span>
        </el-form-item>
        <el-form-item v-if="form.coupon_type === 'discount'" label="最高抵扣">
          <el-input-number v-model="form.max_discount_amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item v-else label="减免金额" required>
          <el-input-number v-model="form.discount_amount" :min="0.01" :precision="2" />
        </el-form-item>
        <el-form-item v-if="form.coupon_type === 'category'" label="分类ID" required>
          <el-input v-model="scopeIdsText" placeholder="逗号分隔分类ID，如 3,5" />
        </el-form-item>
        <el-form-item v-if="form.coupon_type === 'product'" label="商品ID" required>
          <el-input v-model="scopeIdsText" placeholder="逗号分隔商品ID" />
        </el-form-item>
        <el-form-item label="发放总量"><el-input-number v-model="form.total_count" :min="0" /> <span class="tip">0=不限</span></el-form-item>
        <el-form-item label="每人限领"><el-input-number v-model="form.per_user_limit" :min="1" /></el-form-item>
        <el-form-item label="有效期">
          <el-radio-group v-model="form.valid_type">
            <el-radio-button value="fixed">固定日期</el-radio-button>
            <el-radio-button value="relative">领取后N天</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.valid_type === 'fixed'" label="起止">
          <el-date-picker v-model="validRange" type="datetimerange" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
        <el-form-item v-else label="天数"><el-input-number v-model="form.valid_days" :min="1" /></el-form-item>
        <el-form-item label="渠道">
          <el-checkbox-group v-model="form.channels">
            <el-checkbox value="direct">直接领取</el-checkbox>
            <el-checkbox value="popup">弹窗</el-checkbox>
            <el-checkbox value="order_gift">下单赠送</el-checkbox>
            <el-checkbox value="targeted">定向发放</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="用户身份">
          <el-select v-model="form.user_identity" style="width:160px">
            <el-option label="全部" value="all" />
            <el-option label="新用户" value="new" />
            <el-option label="老用户" value="old" />
          </el-select>
        </el-form-item>
        <el-form-item label="可叠秒杀"><el-switch v-model="stackableBool" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width:140px">
            <el-option label="上架" value="on" />
            <el-option label="草稿" value="draft" />
            <el-option label="下架" value="off" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="grantVisible" title="定向发放" width="480px">
      <el-input v-model="grantUserIds" type="textarea" :rows="4" placeholder="用户ID，逗号或换行分隔" />
      <template #footer>
        <el-button @click="grantVisible = false">取消</el-button>
        <el-button type="primary" :loading="granting" @click="doGrant">发放</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" :title="detailTitle" width="720px">
      <el-descriptions v-if="stats" :column="2" border class="mb">
        <el-descriptions-item label="领取">{{ stats.claimed_count }}</el-descriptions-item>
        <el-descriptions-item label="核销">{{ stats.redeemed_count }}</el-descriptions-item>
        <el-descriptions-item label="核销率">{{ stats.redeem_rate }}%</el-descriptions-item>
        <el-descriptions-item label="优惠总额">¥{{ stats.discount_total }}</el-descriptions-item>
      </el-descriptions>
      <el-table :data="detailRows" stripe max-height="360">
        <el-table-column v-for="c in detailCols" :key="c.prop" :prop="c.prop" :label="c.label" />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  COUPON_TYPES, adminCouponClaims, adminCouponRedeems, adminCouponStats,
  copyAdminCoupon, couponTypeLabel, createAdminCoupon, displayStatusLabel,
  grantAdminCoupon, listAdminCoupons, offAdminCoupon, updateAdminCoupon,
} from '../../api/coupon'

const list = ref([])
const total = ref(0)
const query = reactive({ page: 1, status: '', keyword: '' })
const formVisible = ref(false)
const saving = ref(false)
const scopeIdsText = ref('')
const validRange = ref([])
const stackableBool = ref(false)
const form = reactive(emptyForm())
const grantVisible = ref(false)
const granting = ref(false)
const grantUserIds = ref('')
const grantCouponId = ref(0)
const detailVisible = ref(false)
const detailTitle = ref('')
const detailRows = ref([])
const detailCols = ref([])
const stats = ref(null)

function emptyForm() {
  return {
    id: 0, name: '', coupon_type: 'full_reduce', threshold_amount: 100, discount_amount: 10,
    discount_rate: 0.8, max_discount_amount: 0, total_count: 100, per_user_limit: 1,
    valid_type: 'relative', valid_days: 7, channels: ['direct'], user_identity: 'all',
    status: 'on', remark: '', scope_type: 'all',
  }
}

function benefitText(row) {
  if (row.coupon_type === 'discount') return `${(row.discount_rate * 10).toFixed(1)}折` + (row.max_discount_amount ? ` 封顶¥${row.max_discount_amount}` : '')
  if (row.coupon_type === 'no_threshold') return `减¥${row.discount_amount}`
  return `满${row.threshold_amount}减${row.discount_amount}`
}

function onType() {
  if (form.coupon_type === 'no_threshold') form.threshold_amount = 0
}

async function reload() {
  const res = await listAdminCoupons({ ...query, page_size: 20 })
  list.value = res?.list || []
  total.value = res?.total || 0
}

function openCreate() {
  Object.assign(form, emptyForm())
  scopeIdsText.value = ''
  validRange.value = []
  stackableBool.value = false
  formVisible.value = true
}

function openEdit(row) {
  Object.assign(form, {
    ...emptyForm(), ...row,
    channels: row.channels?.length ? [...row.channels] : ['direct'],
  })
  scopeIdsText.value = (row.scopes || []).map((s) => s.ref_id).join(',')
  validRange.value = row.valid_start && row.valid_end ? [row.valid_start, row.valid_end] : []
  stackableBool.value = !!row.stackable
  formVisible.value = true
}

function buildPayload() {
  const scopes = scopeIdsText.value.split(/[,，\s]+/).filter(Boolean).map((id) => ({
    ref_type: form.coupon_type === 'product' ? 'product' : 'category',
    ref_id: Number(id),
  })).filter((s) => s.ref_id > 0)
  const payload = {
    ...form,
    stackable: stackableBool.value ? 1 : 0,
    scopes,
    scope_type: form.coupon_type === 'category' ? 'category' : form.coupon_type === 'product' ? 'product' : 'all',
  }
  if (form.valid_type === 'fixed' && validRange.value?.length === 2) {
    payload.valid_start = validRange.value[0]
    payload.valid_end = validRange.value[1]
  }
  delete payload.id
  delete payload.claimed_count
  delete payload.display_status
  return payload
}

async function save() {
  saving.value = true
  try {
    const payload = buildPayload()
    if (form.id) await updateAdminCoupon(form.id, payload)
    else await createAdminCoupon(payload)
    ElMessage.success('已保存')
    formVisible.value = false
    reload()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

async function doOff(row) {
  await ElMessageBox.confirm('确认下架？')
  await offAdminCoupon(row.id)
  ElMessage.success('已下架')
  reload()
}

async function doCopy(row) {
  await copyAdminCoupon(row.id)
  ElMessage.success('已复制为草稿')
  reload()
}

function openGrant(row) {
  grantCouponId.value = row.id
  grantUserIds.value = ''
  grantVisible.value = true
}

async function doGrant() {
  const ids = grantUserIds.value.split(/[,，\s]+/).map(Number).filter((n) => n > 0)
  if (!ids.length) return ElMessage.warning('请填写用户ID')
  granting.value = true
  try {
    const g = await grantAdminCoupon({ coupon_id: grantCouponId.value, user_ids: ids })
    ElMessage.success(`成功发放 ${g.success_count}/${g.user_count}`)
    grantVisible.value = false
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    granting.value = false
  }
}

async function openStats(row) {
  stats.value = await adminCouponStats(row.id)
  detailTitle.value = `统计 · ${row.name}`
  detailRows.value = []
  detailCols.value = []
  detailVisible.value = true
}

async function openClaims(row) {
  stats.value = null
  const res = await adminCouponClaims(row.id, { page: 1, page_size: 50 })
  detailTitle.value = `领取明细 · ${row.name}`
  detailCols.value = [
    { prop: 'id', label: 'ID' }, { prop: 'user_id', label: '用户' },
    { prop: 'status', label: '状态' }, { prop: 'source', label: '来源' }, { prop: 'valid_end', label: '有效期至' },
  ]
  detailRows.value = res?.list || []
  detailVisible.value = true
}

onMounted(reload)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.mb { margin-bottom: 12px; }
.mt { margin-top: 12px; }
.tip { margin-left: 8px; color: #94a3b8; font-size: 12px; }
</style>
