<template>
  <div>
    <div class="toolbar">
      <h2>主题坑位</h2>
    </div>
    <p class="tip">购买首页「主题好物集市」固定坑位曝光，可自定义标题/封面，并选择跳转店铺、分类或商品。</p>

    <h3>可选坑位</h3>
    <el-table :data="slots" stripe>
      <el-table-column prop="position" label="位置" width="70" />
      <el-table-column prop="name" label="主题" />
      <el-table-column label="占用" min-width="160">
        <template #default="{ row }">
          <span v-if="row.has_active">排队至 {{ row.occupied_until }}</span>
          <span v-else class="ok">空闲，购买后立即生效</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button v-permission="'theme:buy'" type="primary" size="small" @click="openBuy(row)">购买</el-button>
        </template>
      </el-table-column>
    </el-table>

    <h3 class="mt">我的订单</h3>
    <el-table :data="orders" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="坑位" width="120">
        <template #default="{ row }">{{ row.theme_slot_name || ('#' + row.theme_slot_id) }}</template>
      </el-table-column>
      <el-table-column prop="title" label="标题" />
      <el-table-column label="跳转" width="100">
        <template #default="{ row }">{{ themeLinkLabel(row.link_type) }}</template>
      </el-table-column>
      <el-table-column prop="amount" label="金额" width="90" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">{{ themeStatusLabel(row.status) }}</template>
      </el-table-column>
      <el-table-column prop="start_at" label="开始" width="160" />
      <el-table-column prop="end_at" label="结束" width="160" />
    </el-table>

    <el-dialog v-model="buyVisible" title="购买主题坑位" width="560px">
      <el-form label-width="90px">
        <el-form-item label="坑位">{{ buySlot?.name }}（位置 {{ buySlot?.position }}）</el-form-item>
        <el-form-item label="套餐" required>
          <el-select v-model="buyForm.package_id" style="width:100%">
            <el-option v-for="p in pkgs" :key="p.id" :label="`${p.name} · ¥${p.price} · ${p.duration_days}天`" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题" required><el-input v-model="buyForm.title" maxlength="100" /></el-form-item>
        <el-form-item label="副文案"><el-input v-model="buyForm.subtitle" maxlength="200" /></el-form-item>
        <el-form-item label="封面" required>
          <ImageUploader v-model="buyCover" />
        </el-form-item>
        <el-form-item label="跳转" required>
          <el-radio-group v-model="buyForm.link_type" @change="onLinkType">
            <el-radio-button v-for="o in THEME_LINK_OPTIONS" :key="o.value" :value="o.value">{{ o.label }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="buyForm.link_type === 'category'" label="分类" required>
          <el-select v-model="buyForm.link_id" filterable style="width:100%" placeholder="选择分类">
            <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="buyForm.link_type === 'product'" label="商品" required>
          <el-select
            v-model="buyForm.link_id"
            filterable
            remote
            :remote-method="searchProducts"
            :loading="productLoading"
            style="width:100%"
            placeholder="搜索本店商品"
          >
            <el-option v-for="p in products" :key="p.id" :label="`#${p.id} ${p.name}`" :value="p.id" />
          </el-select>
        </el-form-item>
        <p v-else-if="buyForm.link_type === 'shop'" class="tip">将跳转到本店店铺页。</p>
      </el-form>
      <template #footer>
        <el-button @click="buyVisible = false">取消</el-button>
        <el-button type="primary" :loading="buying" @click="confirmBuy">确认扣款</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  THEME_LINK_OPTIONS, merchantBuyTheme, merchantListThemeOrders, merchantListThemePackages,
  merchantListThemeSlots, themeLinkLabel, themeStatusLabel,
} from '../../api/theme'
import { listProductCategories, listProducts } from '../../api/merchant-product'
import ImageUploader from '../../components/merchant/ImageUploader.vue'
import { pickList } from '../../utils/list'

const slots = ref([])
const orders = ref([])
const pkgs = ref([])
const categories = ref([])
const products = ref([])
const productLoading = ref(false)
const buyVisible = ref(false)
const buySlot = ref(null)
const buying = ref(false)
const buyForm = reactive({
  package_id: undefined, title: '', subtitle: '', cover_url: '',
  link_type: 'shop', link_id: undefined,
})
const buyCover = computed({
  get: () => (buyForm.cover_url ? [{ url: buyForm.cover_url }] : []),
  set: (arr) => { buyForm.cover_url = arr?.[arr.length - 1]?.url || '' },
})

async function reload() {
  const [s, o] = await Promise.all([
    merchantListThemeSlots(),
    merchantListThemeOrders({ page: 1, page_size: 50 }),
  ])
  slots.value = pickList(s)
  orders.value = o?.list || []
}

async function openBuy(row) {
  buySlot.value = row
  Object.assign(buyForm, {
    package_id: undefined, title: row.name || '', subtitle: row.desc || '',
    cover_url: row.cover_url || '', link_type: 'shop', link_id: undefined,
  })
  pkgs.value = pickList(await merchantListThemePackages({ theme_slot_id: row.id }))
  if (!categories.value.length) {
    try {
      const res = await listProductCategories({ page: 1, page_size: 100 })
      categories.value = pickList(res)
    } catch { categories.value = [] }
  }
  buyVisible.value = true
}

function onLinkType() {
  buyForm.link_id = undefined
  products.value = []
}

async function searchProducts(kw) {
  productLoading.value = true
  try {
    const res = await listProducts({ page: 1, page_size: 30, name: kw || undefined, status: 'on_sale' })
    products.value = res?.list || []
  } finally { productLoading.value = false }
}

async function confirmBuy() {
  if (!buyForm.package_id || !buyForm.title || !buyForm.cover_url) {
    ElMessage.warning('请完善套餐、标题与封面')
    return
  }
  if (buyForm.link_type !== 'shop' && !buyForm.link_id) {
    ElMessage.warning('请选择跳转目标')
    return
  }
  const pkg = pkgs.value.find((p) => p.id === buyForm.package_id)
  try {
    await ElMessageBox.confirm(`确认支付 ¥${pkg?.price ?? '?'} 购买该坑位？`, '购买确认')
  } catch { return }
  buying.value = true
  try {
    await merchantBuyTheme({
      theme_slot_id: buySlot.value.id,
      package_id: buyForm.package_id,
      title: buyForm.title,
      subtitle: buyForm.subtitle,
      cover_url: buyForm.cover_url,
      link_type: buyForm.link_type,
      link_id: buyForm.link_type === 'shop' ? 0 : buyForm.link_id,
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
.toolbar { margin-bottom: 8px; }
.tip { color: #64748b; font-size: 13px; margin: 0 0 16px; }
.ok { color: #16a34a; }
.mt { margin-top: 24px; }
</style>
