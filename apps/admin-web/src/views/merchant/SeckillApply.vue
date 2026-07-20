<template>
  <div>
    <h2>秒杀报名</h2>
    <el-alert
      v-if="rule"
      type="info"
      :closable="false"
      show-icon
      class="mb"
      :title="`报名费 ¥${rule.apply_fee} / 每店每场最多 ${rule.max_entries_per_shop} 个商品 / 场次时长 ${rule.duration_hours} 小时。开启自动续费后，本场到期且余额充足时将优先自动续报下一场。`"
    />

    <div v-for="s in sessions" :key="s.id" class="session">
      <div class="session-head">
        <div>
          <b>场次 #{{ s.id }}</b>
          <span class="tag" :class="s.status">{{ s.status === 'active' ? '进行中' : s.status }}</span>
        </div>
        <div class="time">{{ s.start_at }} ~ {{ s.end_at }}</div>
        <el-button
          v-permission="'seckill:apply'"
          type="primary"
          size="small"
          :disabled="s.status !== 'active'"
          @click="openApply(s)"
        >报名本场</el-button>
      </div>
    </div>

    <h3>我的报名</h3>
    <el-table :data="entries" v-loading="entryLoading" stripe>
      <el-table-column prop="session_id" label="场次" width="80" />
      <el-table-column prop="product_name" label="商品" min-width="160" />
      <el-table-column prop="seckill_price" label="秒杀价" width="90" />
      <el-table-column prop="seckill_stock" label="秒杀库存" width="90" />
      <el-table-column prop="fee_amount" label="扣费" width="80" />
      <el-table-column label="自动续费" width="110">
        <template #default="{ row }">
          <el-switch
            v-permission="'seckill:apply'"
            :model-value="row.auto_renew === 1"
            @change="(v) => onToggleRenew(row, v)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="报名时间" width="170" />
    </el-table>

    <el-dialog v-model="visible" title="报名秒杀" width="520px" destroy-on-close>
      <el-form label-width="100px">
        <el-form-item label="选择商品" required>
          <el-select
            v-model="form.product_id"
            filterable
            remote
            :remote-method="searchProducts"
            :loading="prodLoading"
            placeholder="搜索本店商品"
            style="width: 100%"
            @change="onPickProduct"
          >
            <el-option
              v-for="p in products"
              :key="p.id"
              :label="`${p.name} (¥${p.sale_price})`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="秒杀价" required>
          <el-input-number v-model="form.seckill_price" :min="0.01" :precision="2" :step="1" />
        </el-form-item>
        <el-form-item label="秒杀库存" required>
          <el-input-number v-model="form.seckill_stock" :min="1" :max="99999" />
        </el-form-item>
        <el-form-item label="自动续费">
          <el-switch v-model="form.auto_renew" />
          <span class="hint">到期后余额足够则自动续报下一场</span>
        </el-form-item>
        <el-form-item label="报名费">
          <span>¥{{ rule?.apply_fee ?? 0 }}（从可用余额扣除）</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">确认报名</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { listProducts } from '../../api/merchant-product'
import {
  applySeckill, fetchMerchantSeckillEntries, fetchMerchantSeckillSessions, setSeckillAutoRenew,
} from '../../api/seckill-wallet'

const rule = ref(null)
const sessions = ref([])
const entries = ref([])
const entryLoading = ref(false)
const visible = ref(false)
const submitting = ref(false)
const prodLoading = ref(false)
const products = ref([])
const picked = ref(null)
const form = reactive({
  session_id: 0,
  product_id: null,
  seckill_price: 1,
  seckill_stock: 10,
  auto_renew: true,
})

async function loadSessions() {
  try {
    const res = await fetchMerchantSeckillSessions()
    rule.value = res?.rule || null
    sessions.value = res?.sessions || []
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function loadEntries() {
  entryLoading.value = true
  try {
    const res = await fetchMerchantSeckillEntries({ page: 1, page_size: 50 })
    entries.value = res?.list || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    entryLoading.value = false
  }
}

async function searchProducts(q) {
  prodLoading.value = true
  try {
    const res = await listProducts({ page: 1, page_size: 30, keyword: q || undefined, status: 'on_sale' })
    products.value = res?.list || []
  } catch {
    products.value = []
  } finally {
    prodLoading.value = false
  }
}

function openApply(s) {
  form.session_id = s.id
  form.product_id = null
  form.seckill_price = 1
  form.seckill_stock = 10
  form.auto_renew = true
  picked.value = null
  visible.value = true
  searchProducts('')
}

function onPickProduct(id) {
  picked.value = products.value.find((p) => p.id === id) || null
  if (picked.value) {
    form.seckill_price = Number(picked.value.sale_price) || 1
  }
}

async function onToggleRenew(row, on) {
  const next = on ? 1 : 0
  try {
    await setSeckillAutoRenew(row.id, next)
    row.auto_renew = next
    ElMessage.success(next ? '已开启自动续费' : '已关闭自动续费')
  } catch (e) {
    ElMessage.error(e.message || '设置失败')
  }
}

async function submit() {
  if (!picked.value) {
    ElMessage.warning('请选择商品')
    return
  }
  submitting.value = true
  try {
    await applySeckill({
      session_id: form.session_id,
      product_id: picked.value.id,
      product_name: picked.value.name,
      product_image: picked.value.main_image || '',
      origin_price: Number(picked.value.market_price || picked.value.sale_price) || 0,
      seckill_price: form.seckill_price,
      seckill_stock: form.seckill_stock,
      auto_renew: form.auto_renew ? 1 : 0,
    })
    ElMessage.success('报名成功')
    visible.value = false
    loadEntries()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadSessions()
  loadEntries()
})
</script>

<style scoped>
.mb { margin: 12px 0 20px; }
.session {
  border: 1px solid #e2e8f0; border-radius: 8px; padding: 14px 16px; margin-bottom: 12px;
}
.session-head { display: flex; align-items: center; gap: 16px; flex-wrap: wrap; }
.session-head .time { color: #64748b; flex: 1; font-size: 13px; }
.tag {
  margin-left: 8px; font-size: 12px; padding: 2px 8px; border-radius: 4px;
  background: #ecfdf5; color: #059669;
}
.tag:not(.active) { background: #f1f5f9; color: #64748b; }
.hint { margin-left: 10px; color: #94a3b8; font-size: 12px; }
h3 { margin: 24px 0 12px; }
</style>
