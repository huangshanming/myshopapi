<template>
  <div>
    <div class="toolbar">
      <h2>首页展位</h2>
      <el-tabs v-model="tab" @tab-change="onTab">
        <el-tab-pane label="套餐" name="pkg" />
        <el-tab-pane label="条带数量" name="settings" />
        <el-tab-pane label="订单" name="orders" />
      </el-tabs>
    </div>

    <template v-if="tab === 'pkg'">
      <div class="toolbar">
        <el-radio-group v-model="slotType" @change="loadPkgs">
          <el-radio-button v-for="o in SLOT_TYPE_OPTIONS" :key="o.value" :value="o.value">{{ o.label }}</el-radio-button>
        </el-radio-group>
        <el-button v-permission="'business:homepage:package'" type="primary" @click="openPkg()">新增套餐</el-button>
      </div>
      <el-table :data="pkgs" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="price" label="价格" width="100" />
        <el-table-column prop="duration_days" label="天数" width="80" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">{{ statusLabel(row.status) }}</template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="70" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button v-permission="'business:homepage:package'" link type="primary" @click="openPkg(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>

    <template v-else-if="tab === 'settings'">
      <el-table :data="settings" v-loading="loading" stripe>
        <el-table-column prop="slot_type" label="类型" width="140">
          <template #default="{ row }">{{ typeLabel(row.slot_type) }}</template>
        </el-table-column>
        <el-table-column label="首页条带数量">
          <template #default="{ row }">
            <el-input-number v-model="row.home_limit" :min="1" :max="50" />
          </template>
        </el-table-column>
      </el-table>
      <el-button class="mt" type="primary" @click="saveSettings">保存</el-button>
    </template>

    <template v-else>
      <div class="toolbar">
        <el-select v-model="orderType" clearable placeholder="类型" style="width: 140px" @change="loadOrders">
          <el-option v-for="o in SLOT_TYPE_OPTIONS" :key="o.value" :label="o.label" :value="o.value" />
        </el-select>
        <el-button v-permission="'business:homepage:grant'" type="warning" @click="openGrant">代开通</el-button>
        <el-button @click="loadOrders">刷新</el-button>
      </div>
      <el-table :data="orders" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="店铺" min-width="140">
          <template #default="{ row }">{{ row.shop_name || ('店铺 #' + row.shop_id) }}</template>
        </el-table-column>
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
      <el-pagination class="pager" layout="prev,pager,next,total" :total="orderTotal" v-model:current-page="orderPage" :page-size="20" @current-change="loadOrders" />
    </template>

    <el-dialog v-model="pkgVisible" :title="pkgForm.id ? '编辑套餐' : '新增套餐'" width="480px">
      <el-form label-width="90px">
        <el-form-item label="类型">
          <el-select v-model="pkgForm.slot_type" style="width: 100%">
            <el-option v-for="o in SLOT_TYPE_OPTIONS" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称"><el-input v-model="pkgForm.name" /></el-form-item>
        <el-form-item label="价格"><el-input-number v-model="pkgForm.price" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="天数"><el-input-number v-model="pkgForm.duration_days" :min="1" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="pkgForm.sort" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="pkgForm.status"><el-option label="上架" value="on" /><el-option label="下架" value="off" /></el-select>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="pkgForm.remark" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pkgVisible = false">取消</el-button>
        <el-button type="primary" @click="savePkg">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="grantVisible" title="代开通展位" width="520px" @open="onGrantOpen">
      <el-form label-width="90px">
        <el-form-item label="套餐" required>
          <el-select
            v-model="grantForm.package_id"
            filterable
            style="width: 100%"
            placeholder="请选择套餐"
            @change="onGrantPkgChange"
          >
            <el-option
              v-for="p in allPkgs"
              :key="p.id"
              :label="`${typeLabel(p.slot_type)} · ${p.name} · ¥${p.price} · ${p.duration_days}天`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="店铺" required>
          <el-select
            v-model="grantForm.shop_id"
            filterable
            remote
            clearable
            :remote-method="searchShops"
            :loading="shopLoading"
            style="width: 100%"
            placeholder="搜索店铺名称"
            @change="onGrantShopChange"
          >
            <el-option
              v-for="s in shops"
              :key="s.id"
              :label="`#${s.id} ${s.name}`"
              :value="s.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="grantSlotType === 'article'" label="文章" required>
          <el-select
            v-model="grantForm.target_id"
            filterable
            remote
            clearable
            :remote-method="searchArticles"
            :loading="articleLoading"
            :disabled="!grantForm.shop_id"
            style="width: 100%"
            placeholder="先选店铺，再搜索文章标题"
          >
            <el-option
              v-for="a in articles"
              :key="a.id"
              :label="`#${a.id} ${a.title}`"
              :value="a.id"
            />
          </el-select>
        </el-form-item>
        <p v-else class="hint">品牌/优质商户展位：开通后直接作用于所选店铺，无需再选文章。</p>
      </el-form>
      <template #footer>
        <el-button @click="grantVisible = false">取消</el-button>
        <el-button type="primary" @click="doGrant">开通</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  SLOT_TYPE_OPTIONS, createHomepagePackage, grantHomepageOrder, listHomepageOrders,
  listHomepagePackages, listHomepageSettings, updateHomepagePackage, updateHomepageSettings,
  slotPaySourceLabel, slotStatusLabel, slotTargetLabel, slotTypeLabel,
} from '../../api/homepage'
import { fetchShops } from '../../api/merchant'
import { listArticles } from '../../api/admin-article'

const tab = ref('pkg')
const slotType = ref('brand_shop')
const loading = ref(false)
const pkgs = ref([])
const settings = ref([])
const orders = ref([])
const orderTotal = ref(0)
const orderPage = ref(1)
const orderType = ref('')
const allPkgs = ref([])

const shops = ref([])
const shopLoading = ref(false)
const articles = ref([])
const articleLoading = ref(false)

const pkgVisible = ref(false)
const pkgForm = reactive({ id: 0, slot_type: 'brand_shop', name: '', price: 500, duration_days: 10, status: 'on', sort: 0, remark: '' })
const grantVisible = ref(false)
const grantForm = reactive({ shop_id: undefined, package_id: undefined, target_id: undefined })

const grantSlotType = computed(() => {
  const p = allPkgs.value.find((x) => x.id === grantForm.package_id)
  return p?.slot_type || ''
})

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

async function loadPkgs() {
  loading.value = true
  try {
    const res = await listHomepagePackages({ slot_type: slotType.value })
    pkgs.value = res || []
  } finally {
    loading.value = false
  }
}

async function loadSettings() {
  loading.value = true
  try {
    const res = await listHomepageSettings()
    settings.value = res || []
  } finally {
    loading.value = false
  }
}

async function loadOrders() {
  loading.value = true
  try {
    const res = await listHomepageOrders({ page: orderPage.value, page_size: 20, slot_type: orderType.value || undefined })
    orders.value = res?.list || []
    orderTotal.value = res?.total || 0
  } finally {
    loading.value = false
  }
}

function onTab() {
  if (tab.value === 'pkg') loadPkgs()
  else if (tab.value === 'settings') loadSettings()
  else loadOrders()
}

function openPkg(row) {
  if (row) Object.assign(pkgForm, row)
  else Object.assign(pkgForm, { id: 0, slot_type: slotType.value, name: '', price: 500, duration_days: 10, status: 'on', sort: 0, remark: '' })
  pkgVisible.value = true
}

async function savePkg() {
  try {
    if (pkgForm.id) await updateHomepagePackage(pkgForm.id, pkgForm)
    else await createHomepagePackage(pkgForm)
    ElMessage.success('已保存')
    pkgVisible.value = false
    loadPkgs()
    const res = await listHomepagePackages({})
    allPkgs.value = (res || []).filter((p) => p.status === 'on' || !p.status)
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function saveSettings() {
  try {
    await updateHomepageSettings(settings.value)
    ElMessage.success('已保存')
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function searchShops(keyword) {
  shopLoading.value = true
  try {
    const res = await fetchShops({ name: keyword || undefined, page: 1, page_size: 50 })
    shops.value = res?.list || res?.items || []
  } finally {
    shopLoading.value = false
  }
}

async function searchArticles(keyword) {
  if (!grantForm.shop_id) {
    articles.value = []
    return
  }
  articleLoading.value = true
  try {
    const res = await listArticles({
      shop_id: grantForm.shop_id,
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

function openGrant() {
  Object.assign(grantForm, { shop_id: undefined, package_id: undefined, target_id: undefined })
  articles.value = []
  grantVisible.value = true
}

async function onGrantOpen() {
  if (!shops.value.length) await searchShops('')
}

function onGrantPkgChange() {
  grantForm.target_id = undefined
  articles.value = []
  if (grantSlotType.value === 'article' && grantForm.shop_id) {
    searchArticles('')
  }
}

function onGrantShopChange() {
  grantForm.target_id = undefined
  articles.value = []
  if (grantSlotType.value === 'article' && grantForm.shop_id) {
    searchArticles('')
  }
}

async function doGrant() {
  if (!grantForm.package_id) {
    ElMessage.warning('请选择套餐')
    return
  }
  if (!grantForm.shop_id) {
    ElMessage.warning('请选择店铺')
    return
  }
  if (grantSlotType.value === 'article' && !grantForm.target_id) {
    ElMessage.warning('请选择文章')
    return
  }
  try {
    await grantHomepageOrder({
      shop_id: grantForm.shop_id,
      package_id: grantForm.package_id,
      target_id: grantSlotType.value === 'article' ? grantForm.target_id : 0,
    })
    ElMessage.success('已开通')
    grantVisible.value = false
    loadOrders()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(async () => {
  loadPkgs()
  const res = await listHomepagePackages({})
  allPkgs.value = res || []
})
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; align-items: center; margin-bottom: 12px; flex-wrap: wrap; }
.toolbar h2 { margin-right: auto; }
.pager { margin-top: 12px; }
.mt { margin-top: 16px; }
.hint { margin: 0 0 0 90px; color: #888; font-size: 12px; line-height: 1.5; }
</style>
