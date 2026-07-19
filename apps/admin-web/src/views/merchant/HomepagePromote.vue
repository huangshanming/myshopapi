<template>
  <div>
    <div class="toolbar">
      <h2>首页推广</h2>
      <el-radio-group v-model="slotType" @change="reload">
        <el-radio-button v-for="o in SLOT_TYPE_OPTIONS" :key="o.value" :value="o.value">{{ o.label }}</el-radio-button>
      </el-radio-group>
    </div>

    <h3>可选套餐</h3>
    <el-table :data="pkgs" v-loading="loading" stripe>
      <el-table-column prop="name" label="套餐" />
      <el-table-column prop="price" label="价格" width="100" />
      <el-table-column prop="duration_days" label="天数" width="80" />
      <el-table-column label="操作" width="140">
        <template #default="{ row }">
          <el-button v-permission="'homepage:buy'" type="primary" size="small" @click="buy(row)">购买</el-button>
        </template>
      </el-table-column>
    </el-table>

    <h3 class="mt">我的订单</h3>
    <el-table :data="orders" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="套餐" min-width="120">
        <template #default="{ row }">{{ row.package_name || ('套餐 #' + row.package_id) }}</template>
      </el-table-column>
      <el-table-column label="类型" width="100">
        <template #default="{ row }">{{ typeLabel(row.slot_type) }}</template>
      </el-table-column>
      <el-table-column label="推广目标" min-width="160">
        <template #default="{ row }">{{ targetLabel(row) }}</template>
      </el-table-column>
      <el-table-column prop="amount" label="金额" width="90" />
      <el-table-column label="来源" width="110">
        <template #default="{ row }">{{ paySourceLabel(row.pay_source) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }">{{ statusLabel(row.status) }}</template>
      </el-table-column>
      <el-table-column prop="start_at" label="开始" width="160" />
      <el-table-column prop="end_at" label="结束" width="160" />
    </el-table>

    <el-dialog v-model="buyVisible" title="确认购买" width="480px">
      <p>套餐：{{ buyPkg?.name }} · ¥{{ buyPkg?.price }} · {{ buyPkg?.duration_days }}天</p>
      <p class="tip">将从店铺钱包扣款，购买后立即生效（续费顺延）。</p>
      <el-form v-if="slotType === 'article'" label-width="80px" class="mt">
        <el-form-item label="文章" required>
          <el-select
            v-model="articleId"
            filterable
            remote
            clearable
            :remote-method="searchArticles"
            :loading="articleLoading"
            style="width: 100%"
            placeholder="搜索本店已通过文章"
          >
            <el-option
              v-for="a in articles"
              :key="a.id"
              :label="`#${a.id} ${a.title}`"
              :value="a.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="buyVisible = false">取消</el-button>
        <el-button type="primary" :loading="buying" @click="confirmBuy">确认扣款</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  SLOT_TYPE_OPTIONS, merchantBuyHomepage, merchantListHomepageOrders, merchantListHomepagePackages,
  slotPaySourceLabel, slotStatusLabel, slotTargetLabel, slotTypeLabel,
} from '../../api/homepage'
import { listMyArticles } from '../../api/merchant-article'

function typeLabel(t) {
  return slotTypeLabel(t)
}
function statusLabel(s) {
  return slotStatusLabel(s)
}
function paySourceLabel(s) {
  return slotPaySourceLabel(s)
}
function targetLabel(row) {
  return slotTargetLabel(row)
}

const slotType = ref('brand_shop')
const pkgs = ref([])
const orders = ref([])
const loading = ref(false)
const buyVisible = ref(false)
const buyPkg = ref(null)
const articleId = ref(undefined)
const articles = ref([])
const articleLoading = ref(false)
const buying = ref(false)

async function reload() {
  loading.value = true
  try {
    const [p, o] = await Promise.all([
      merchantListHomepagePackages({ slot_type: slotType.value }),
      merchantListHomepageOrders({ slot_type: slotType.value, page: 1, page_size: 50 }),
    ])
    pkgs.value = p || []
    orders.value = o?.list || []
  } finally {
    loading.value = false
  }
}

async function searchArticles(keyword) {
  articleLoading.value = true
  try {
    const res = await listMyArticles({
      title: keyword || undefined,
      page: 1,
      page_size: 50,
      audit_status: 'approved',
    })
    articles.value = res?.list || []
  } finally {
    articleLoading.value = false
  }
}

async function buy(row) {
  buyPkg.value = row
  articleId.value = undefined
  articles.value = []
  buyVisible.value = true
  if (slotType.value === 'article') {
    await searchArticles('')
  }
}

async function confirmBuy() {
  if (slotType.value === 'article' && !articleId.value) {
    ElMessage.warning('请选择文章')
    return
  }
  try {
    await ElMessageBox.confirm(`确认支付 ¥${buyPkg.value.price}？`, '购买确认')
  } catch {
    return
  }
  buying.value = true
  try {
    await merchantBuyHomepage({
      package_id: buyPkg.value.id,
      target_id: slotType.value === 'article' ? articleId.value : 0,
    })
    ElMessage.success('购买成功')
    buyVisible.value = false
    reload()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    buying.value = false
  }
}

onMounted(reload)
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; align-items: center; margin-bottom: 16px; }
.toolbar h2 { margin-right: auto; }
.mt { margin-top: 16px; }
.tip { color: #888; font-size: 13px; }
</style>
