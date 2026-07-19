<template>
  <div>
    <div class="toolbar">
      <h2>主题集市</h2>
      <el-tabs v-model="tab">
        <el-tab-pane label="坑位" name="slots" />
        <el-tab-pane label="套餐" name="pkg" />
        <el-tab-pane label="订单" name="orders" />
      </el-tabs>
    </div>

    <template v-if="tab === 'slots'">
      <el-table :data="slots" v-loading="loading" stripe>
        <el-table-column prop="position" label="位置" width="70" />
        <el-table-column prop="slot_key" label="Key" width="100" />
        <el-table-column label="封面" width="100">
          <template #default="{ row }">
            <el-image v-if="row.cover_url" :src="row.cover_url" style="width:72px;height:48px" fit="cover" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="默认标题" min-width="120" />
        <el-table-column prop="desc" label="副文案" min-width="140" />
        <el-table-column label="默认跳转" width="120">
          <template #default="{ row }">{{ themeLinkLabel(row.default_link_type) }} #{{ row.default_link_id || '-' }}</template>
        </el-table-column>
        <el-table-column label="占用" width="160">
          <template #default="{ row }">
            <span v-if="row.has_active">至 {{ row.occupied_until }}</span>
            <span v-else class="muted">空闲</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">{{ themeStatusLabel(row.status) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button v-permission="'marketing:theme:list'" link type="primary" @click="openSlot(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>

    <template v-else-if="tab === 'pkg'">
      <div class="toolbar">
        <el-button v-permission="'marketing:theme:package'" type="primary" @click="openPkg()">新增套餐</el-button>
      </div>
      <el-table :data="pkgs" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="名称" />
        <el-table-column label="适用坑位" width="120">
          <template #default="{ row }">{{ row.theme_slot_id ? ('坑位#' + row.theme_slot_id) : '通用' }}</template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="100" />
        <el-table-column prop="duration_days" label="天数" width="80" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">{{ themeStatusLabel(row.status) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button v-permission="'marketing:theme:package'" link type="primary" @click="openPkg(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>

    <template v-else>
      <div class="toolbar">
        <el-button v-permission="'marketing:theme:grant'" type="warning" @click="openGrant">代开通</el-button>
        <el-button @click="loadOrders">刷新</el-button>
      </div>
      <el-table :data="orders" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="店铺" min-width="120">
          <template #default="{ row }">{{ row.shop_name || ('#' + row.shop_id) }}</template>
        </el-table-column>
        <el-table-column label="坑位" width="120">
          <template #default="{ row }">{{ row.theme_slot_name || ('#' + row.theme_slot_id) }}</template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="120" />
        <el-table-column label="跳转" width="110">
          <template #default="{ row }">{{ themeLinkLabel(row.link_type) }}</template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="90" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">{{ themeStatusLabel(row.status) }}</template>
        </el-table-column>
        <el-table-column prop="start_at" label="开始" width="160" />
        <el-table-column prop="end_at" label="结束" width="160" />
      </el-table>
      <el-pagination class="pager" layout="prev,pager,next,total" :total="orderTotal" v-model:current-page="orderPage" :page-size="20" @current-change="loadOrders" />
    </template>

    <el-dialog v-model="slotVisible" title="编辑坑位" width="520px">
      <el-form label-width="100px">
        <el-form-item label="标题"><el-input v-model="slotForm.name" /></el-form-item>
        <el-form-item label="副文案"><el-input v-model="slotForm.desc" /></el-form-item>
        <el-form-item label="封面">
          <ImageUploader v-model="slotCover" :upload-fn="uploadShopImage" />
        </el-form-item>
        <el-form-item label="默认跳转">
          <el-select v-model="slotForm.default_link_type" style="width:140px">
            <el-option v-for="o in THEME_DEFAULT_LINK_OPTIONS" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="slotForm.default_link_type === 'category'" label="分类ID">
          <el-input-number v-model="slotForm.default_link_id" :min="0" />
        </el-form-item>
        <el-form-item label="排序"><el-input-number v-model="slotForm.sort" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="slotForm.status" style="width:120px">
            <el-option label="上架" value="on" /><el-option label="下架" value="off" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="slotVisible = false">取消</el-button>
        <el-button type="primary" @click="saveSlot">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="pkgVisible" :title="pkgForm.id ? '编辑套餐' : '新增套餐'" width="480px">
      <el-form label-width="100px">
        <el-form-item label="名称"><el-input v-model="pkgForm.name" /></el-form-item>
        <el-form-item label="适用坑位">
          <el-select v-model="pkgForm.theme_slot_id" clearable placeholder="通用" style="width:100%">
            <el-option label="通用（任意坑位）" :value="0" />
            <el-option v-for="s in slots" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
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

    <el-dialog v-model="grantVisible" title="代开通主题坑位" width="560px">
      <el-form label-width="90px">
        <el-form-item label="坑位" required>
          <el-select v-model="grantForm.theme_slot_id" style="width:100%" @change="onGrantSlot">
            <el-option v-for="s in slots.filter(x=>x.status==='on')" :key="s.id" :label="`位置${s.position} · ${s.name}`" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="套餐" required>
          <el-select v-model="grantForm.package_id" style="width:100%">
            <el-option v-for="p in grantPkgs" :key="p.id" :label="`${p.name} · ¥${p.price} · ${p.duration_days}天`" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="店铺" required>
          <el-select v-model="grantForm.shop_id" filterable remote :remote-method="searchShops" :loading="shopLoading" style="width:100%" placeholder="搜索店铺">
            <el-option v-for="s in shops" :key="s.id" :label="`#${s.id} ${s.name}`" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题" required><el-input v-model="grantForm.title" /></el-form-item>
        <el-form-item label="副文案"><el-input v-model="grantForm.subtitle" /></el-form-item>
        <el-form-item label="封面" required>
          <ImageUploader v-model="grantCover" :upload-fn="uploadShopImage" />
        </el-form-item>
        <el-form-item label="跳转" required>
          <el-radio-group v-model="grantForm.link_type">
            <el-radio-button v-for="o in THEME_LINK_OPTIONS" :key="o.value" :value="o.value">{{ o.label }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="grantForm.link_type === 'category'" label="分类ID" required>
          <el-input-number v-model="grantForm.link_id" :min="1" />
        </el-form-item>
        <el-form-item v-if="grantForm.link_type === 'product'" label="商品ID" required>
          <el-input-number v-model="grantForm.link_id" :min="1" />
        </el-form-item>
        <p v-if="grantForm.link_type === 'shop'" class="hint">跳转本店：开通后自动绑定所选店铺。</p>
      </el-form>
      <template #footer>
        <el-button @click="grantVisible = false">取消</el-button>
        <el-button type="primary" @click="doGrant">开通</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  THEME_DEFAULT_LINK_OPTIONS, THEME_LINK_OPTIONS, createThemePackage, grantThemeOrder,
  listThemeOrders, listThemePackages, listThemeSlots, themeLinkLabel, themeStatusLabel,
  updateThemePackage, updateThemeSlot,
} from '../../api/theme'
import { fetchShops, uploadShopImage } from '../../api/merchant'
import ImageUploader from '../../components/merchant/ImageUploader.vue'

const tab = ref('slots')
const loading = ref(false)
const slots = ref([])
const pkgs = ref([])
const orders = ref([])
const orderTotal = ref(0)
const orderPage = ref(1)

const slotVisible = ref(false)
const slotForm = reactive({})
const slotCover = computed({
  get: () => (slotForm.cover_url ? [{ url: slotForm.cover_url }] : []),
  set: (arr) => { slotForm.cover_url = arr?.[arr.length - 1]?.url || '' },
})

const pkgVisible = ref(false)
const pkgForm = reactive({ id: 0, theme_slot_id: 0, name: '', price: 0, duration_days: 7, status: 'on', sort: 0, remark: '' })

const grantVisible = ref(false)
const grantForm = reactive({
  theme_slot_id: undefined, package_id: undefined, shop_id: undefined,
  title: '', subtitle: '', cover_url: '', link_type: 'shop', link_id: undefined,
})
const grantCover = computed({
  get: () => (grantForm.cover_url ? [{ url: grantForm.cover_url }] : []),
  set: (arr) => { grantForm.cover_url = arr?.[arr.length - 1]?.url || '' },
})
const grantPkgs = ref([])
const shops = ref([])
const shopLoading = ref(false)

async function loadSlots() {
  loading.value = true
  try { slots.value = await listThemeSlots() || [] }
  catch (e) { ElMessage.error(e.message) }
  finally { loading.value = false }
}

async function loadPkgs() {
  try { pkgs.value = await listThemePackages() || [] }
  catch (e) { ElMessage.error(e.message) }
}

async function loadOrders() {
  try {
    const res = await listThemeOrders({ page: orderPage.value, page_size: 20 })
    orders.value = res?.list || []
    orderTotal.value = res?.total || 0
  } catch (e) { ElMessage.error(e.message) }
}

function openSlot(row) {
  Object.assign(slotForm, { ...row })
  slotVisible.value = true
}

async function saveSlot() {
  try {
    await updateThemeSlot(slotForm.id, {
      name: slotForm.name, desc: slotForm.desc, cover_url: slotForm.cover_url,
      default_link_type: slotForm.default_link_type,
      default_link_id: slotForm.default_link_type === 'none' ? 0 : (slotForm.default_link_id || 0),
      status: slotForm.status, sort: slotForm.sort,
    })
    ElMessage.success('已保存')
    slotVisible.value = false
    loadSlots()
  } catch (e) { ElMessage.error(e.message) }
}

function openPkg(row) {
  if (row) Object.assign(pkgForm, { ...row })
  else Object.assign(pkgForm, { id: 0, theme_slot_id: 0, name: '', price: 0, duration_days: 7, status: 'on', sort: 0, remark: '' })
  pkgVisible.value = true
}

async function savePkg() {
  try {
    const payload = { ...pkgForm, theme_slot_id: pkgForm.theme_slot_id || 0 }
    if (pkgForm.id) await updateThemePackage(pkgForm.id, payload)
    else await createThemePackage(payload)
    ElMessage.success('已保存')
    pkgVisible.value = false
    loadPkgs()
  } catch (e) { ElMessage.error(e.message) }
}

async function onGrantSlot() {
  grantPkgs.value = await listThemePackages({ theme_slot_id: grantForm.theme_slot_id }) || []
}

async function searchShops(kw) {
  shopLoading.value = true
  try {
    const res = await fetchShops({ page: 1, page_size: 30, name: kw || undefined })
    shops.value = res?.list || []
  } finally { shopLoading.value = false }
}

function openGrant() {
  Object.assign(grantForm, {
    theme_slot_id: slots.value.find((s) => s.status === 'on')?.id,
    package_id: undefined, shop_id: undefined, title: '', subtitle: '', cover_url: '',
    link_type: 'shop', link_id: undefined,
  })
  onGrantSlot()
  grantVisible.value = true
}

async function doGrant() {
  if (!grantForm.theme_slot_id || !grantForm.package_id || !grantForm.shop_id) {
    ElMessage.warning('请完善坑位/套餐/店铺')
    return
  }
  if (!grantForm.title || !grantForm.cover_url) {
    ElMessage.warning('请填写标题并上传封面')
    return
  }
  const payload = {
    ...grantForm,
    link_id: grantForm.link_type === 'shop' ? grantForm.shop_id : grantForm.link_id,
  }
  try {
    await grantThemeOrder(payload)
    ElMessage.success('已开通')
    grantVisible.value = false
    loadOrders()
    loadSlots()
  } catch (e) { ElMessage.error(e.message) }
}

watch(tab, (t) => {
  if (t === 'slots') loadSlots()
  else if (t === 'pkg') loadPkgs()
  else loadOrders()
})

onMounted(() => { loadSlots(); loadPkgs() })
</script>

<style scoped>
.toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; flex-wrap: wrap; }
.toolbar h2 { margin-right: auto; }
.pager { margin-top: 16px; }
.muted { color: #94a3b8; }
.hint { color: #64748b; font-size: 13px; margin: 0; }
</style>
