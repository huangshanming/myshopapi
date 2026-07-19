<template>
  <div>
    <div class="toolbar">
      <h2>店铺优惠券</h2>
      <el-button v-permission="'coupon:edit'" type="primary" @click="openCreate">创建</el-button>
    </div>
    <p class="tip">本店券仅在本店订单结算时可使用；可配置满减/无门槛/品类/商品/折扣。</p>
    <el-table :data="list" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="名称" />
      <el-table-column label="类型" width="110">
        <template #default="{ row }">{{ couponTypeLabel(row.coupon_type) }}</template>
      </el-table-column>
      <el-table-column label="库存" width="100">
        <template #default="{ row }">{{ row.claimed_count }}/{{ row.total_count || '∞' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }">{{ displayStatusLabel(row.display_status || row.status) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="280">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link @click="doCopy(row)">复制</el-button>
          <el-button v-permission="'coupon:grant'" link @click="openGrant(row)">发放</el-button>
          <el-button link type="danger" @click="doOff(row)">下架</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="formVisible" :title="form.id ? '编辑' : '创建'" width="600px">
      <el-form label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.coupon_type" style="width:100%">
            <el-option v-for="t in COUPON_TYPES" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.coupon_type !== 'no_threshold' && form.coupon_type !== 'discount'" label="门槛">
          <el-input-number v-model="form.threshold_amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item v-if="form.coupon_type === 'discount'" label="折扣">
          <el-input-number v-model="form.discount_rate" :min="0.01" :max="0.99" :step="0.05" :precision="2" />
        </el-form-item>
        <el-form-item v-else label="减免">
          <el-input-number v-model="form.discount_amount" :min="0.01" :precision="2" />
        </el-form-item>
        <el-form-item v-if="form.coupon_type === 'category' || form.coupon_type === 'product'" label="范围ID">
          <el-input v-model="scopeIdsText" placeholder="逗号分隔" />
        </el-form-item>
        <el-form-item label="总量"><el-input-number v-model="form.total_count" :min="0" /></el-form-item>
        <el-form-item label="限领"><el-input-number v-model="form.per_user_limit" :min="1" /></el-form-item>
        <el-form-item label="领后天数"><el-input-number v-model="form.valid_days" :min="1" /></el-form-item>
        <el-form-item label="渠道">
          <el-checkbox-group v-model="form.channels">
            <el-checkbox value="direct">直接领取</el-checkbox>
            <el-checkbox value="order_gift">下单赠送</el-checkbox>
            <el-checkbox value="targeted">定向</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status"><el-option label="上架" value="on" /><el-option label="草稿" value="draft" /></el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="grantVisible" title="定向发放" width="420px">
      <el-input v-model="grantUserIds" type="textarea" :rows="3" placeholder="用户ID，逗号分隔" />
      <template #footer>
        <el-button @click="grantVisible = false">取消</el-button>
        <el-button type="primary" @click="doGrant">发放</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  COUPON_TYPES, copyMerchantCoupon, couponTypeLabel, createMerchantCoupon,
  displayStatusLabel, grantMerchantCoupon, listMerchantCoupons, offMerchantCoupon, updateMerchantCoupon,
} from '../../api/coupon'

const list = ref([])
const formVisible = ref(false)
const saving = ref(false)
const scopeIdsText = ref('')
const form = reactive({
  id: 0, name: '', coupon_type: 'full_reduce', threshold_amount: 50, discount_amount: 5,
  discount_rate: 0.9, total_count: 50, per_user_limit: 1, valid_type: 'relative', valid_days: 7,
  channels: ['direct'], status: 'on', user_identity: 'all',
})
const grantVisible = ref(false)
const grantUserIds = ref('')
const grantCouponId = ref(0)

async function reload() {
  const res = await listMerchantCoupons({ page: 1, page_size: 50 })
  list.value = res?.list || []
}

function openCreate() {
  Object.assign(form, {
    id: 0, name: '', coupon_type: 'full_reduce', threshold_amount: 50, discount_amount: 5,
    discount_rate: 0.9, total_count: 50, per_user_limit: 1, valid_type: 'relative', valid_days: 7,
    channels: ['direct'], status: 'on', user_identity: 'all',
  })
  scopeIdsText.value = ''
  formVisible.value = true
}

function openEdit(row) {
  Object.assign(form, { ...row, channels: row.channels?.length ? [...row.channels] : ['direct'] })
  scopeIdsText.value = (row.scopes || []).map((s) => s.ref_id).join(',')
  formVisible.value = true
}

function payload() {
  const scopes = scopeIdsText.value.split(/[,，\s]+/).filter(Boolean).map((id) => ({
    ref_type: form.coupon_type === 'product' ? 'product' : 'category',
    ref_id: Number(id),
  })).filter((s) => s.ref_id > 0)
  return {
    name: form.name, coupon_type: form.coupon_type, threshold_amount: form.threshold_amount,
    discount_amount: form.discount_amount, discount_rate: form.discount_rate,
    total_count: form.total_count, per_user_limit: form.per_user_limit,
    valid_type: 'relative', valid_days: form.valid_days, channels: form.channels,
    status: form.status, user_identity: form.user_identity, stackable: 0, scopes,
    scope_type: form.coupon_type === 'category' ? 'category' : form.coupon_type === 'product' ? 'product' : 'all',
  }
}

async function save() {
  saving.value = true
  try {
    if (form.id) await updateMerchantCoupon(form.id, payload())
    else await createMerchantCoupon(payload())
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
  await offMerchantCoupon(row.id)
  reload()
}

async function doCopy(row) {
  await copyMerchantCoupon(row.id)
  ElMessage.success('已复制')
  reload()
}

function openGrant(row) {
  grantCouponId.value = row.id
  grantUserIds.value = ''
  grantVisible.value = true
}

async function doGrant() {
  const ids = grantUserIds.value.split(/[,，\s]+/).map(Number).filter((n) => n > 0)
  const g = await grantMerchantCoupon({ coupon_id: grantCouponId.value, user_ids: ids })
  ElMessage.success(`成功 ${g.success_count}`)
  grantVisible.value = false
}

onMounted(reload)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; margin-bottom: 8px; }
.tip { color: #64748b; font-size: 13px; margin: 0 0 12px; }
</style>
