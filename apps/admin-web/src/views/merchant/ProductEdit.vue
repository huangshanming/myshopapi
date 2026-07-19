<template>
  <div class="edit-page" v-loading="loading">
    <div class="toolbar">
      <h2>{{ id ? '编辑商品' : '发布商品' }}</h2>
      <div>
        <el-button @click="$router.back()">返回</el-button>
        <el-button @click="save('draft')">存草稿</el-button>
        <el-button type="primary" @click="save('on_sale')">保存并上架</el-button>
      </div>
    </div>

    <el-form label-width="100px" class="form">
      <el-form-item label="名称" required>
        <el-input v-model="form.name" maxlength="200" show-word-limit />
      </el-form-item>
      <el-form-item label="副标题">
        <el-input v-model="form.subtitle" />
      </el-form-item>
      <el-form-item label="分类" required>
        <el-tree-select
          v-model="form.category_id"
          :data="catTree"
          :props="{ label: 'name', value: 'id', children: 'children' }"
          check-strictly
          filterable
          clearable
          placeholder="请选择商品分类"
          style="width: 320px"
        />
      </el-form-item>
      <el-form-item label="类型">
        <el-radio-group v-model="form.product_type">
          <el-radio-button value="physical">实物</el-radio-button>
          <el-radio-button value="fresh">生鲜</el-radio-button>
          <el-radio-button value="virtual">虚拟</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item v-if="form.product_type === 'fresh'" label="保质期(天)">
        <el-input-number v-model="form.shelf_life" :min="0" />
      </el-form-item>
      <el-form-item v-if="form.product_type === 'fresh'" label="存储条件">
        <el-input v-model="form.storage_condition" />
      </el-form-item>
      <el-form-item label="商品图片">
        <ImageUploader v-model="form.images" />
      </el-form-item>
      <el-form-item label="详情描述">
        <RichTextEditor v-model="form.description" />
      </el-form-item>

      <el-divider>规格与 SKU</el-divider>
      <SpecEditor v-model="form.spec_json" @generate="generateSkus" />
      <div class="sku-block">
        <SkuTable v-model="form.skus" />
      </div>

      <template v-if="!form.spec_json?.length">
        <el-form-item label="默认售价">
          <el-input-number v-model="form.sale_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="默认库存">
          <el-input-number v-model="form.stock" :min="0" />
        </el-form-item>
      </template>
    </el-form>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import SpecEditor from '../../components/merchant/SpecEditor.vue'
import SkuTable from '../../components/merchant/SkuTable.vue'
import ImageUploader from '../../components/merchant/ImageUploader.vue'
import RichTextEditor from '../../components/merchant/RichTextEditor.vue'
import { createProduct, getProduct, updateProduct, listProductCategories } from '../../api/merchant-product'
import { pickList } from '../../utils/list'
import { buildCategoryTree } from '../../utils/categoryTree'

const route = useRoute()
const router = useRouter()
const id = Number(route.params.id || 0)
const loading = ref(false)
const catTree = ref([])

const form = reactive({
  name: '',
  subtitle: '',
  description: '',
  category_id: undefined,
  product_type: 'physical',
  shelf_life: 0,
  storage_condition: '',
  sale_price: 0,
  stock: 0,
  spec_json: [],
  skus: [],
  images: [],
  attrs: [],
  tag_ids: [],
})

function cartesian(specs) {
  const named = (specs || []).filter((s) => s.name && s.values?.length)
  if (!named.length) return [{}]
  return named.reduce(
    (acc, spec) => {
      const next = []
      for (const row of acc) {
        for (const v of spec.values) {
          next.push({ ...row, [spec.name]: v })
        }
      }
      return next
    },
    [{}]
  )
}

function generateSkus() {
  const combos = cartesian(form.spec_json)
  const oldMap = new Map(
    form.skus.map((s) => {
      const key = JSON.stringify(normalizeSpec(s.spec_values))
      return [key, s]
    })
  )
  form.skus = combos.map((spec) => {
    const key = JSON.stringify(normalizeSpec(spec))
    const old = oldMap.get(key)
    return {
      id: old?.id || 0,
      spec_values: spec,
      sale_price: old?.sale_price ?? form.sale_price ?? 0,
      stock: old?.stock ?? form.stock ?? 0,
      stock_warn: old?.stock_warn ?? 10,
      status: old?.status || 'enabled',
      market_price: old?.market_price || 0,
      cost_price: old?.cost_price || 0,
      barcode: old?.barcode || '',
    }
  })
  ElMessage.success(`已生成 ${form.skus.length} 个 SKU`)
}

function normalizeSpec(sv) {
  if (!sv) return {}
  if (typeof sv === 'string') {
    try {
      sv = JSON.parse(sv)
    } catch {
      return {}
    }
  }
  const keys = Object.keys(sv).sort()
  const out = {}
  keys.forEach((k) => {
    out[k] = sv[k]
  })
  return out
}

function toRichHtml(text) {
  if (!text) return ''
  const t = String(text).trim()
  if (!t) return ''
  // 已是 HTML 则直接用
  if (/<[a-z][\s\S]*>/i.test(t)) return t
  // 纯文本按空行分段，保留换行
  return t
    .split(/\n{2,}/)
    .map((block) => {
      const body = block
        .split('\n')
        .map((line) => line.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;'))
        .join('<br/>')
      return `<p>${body}</p>`
    })
    .join('')
}

async function load() {
  if (!id) return
  loading.value = true
  try {
    const res = await getProduct(id)
    const d = res || {}
    const p = d.product || d
    Object.assign(form, {
      name: p.name || '',
      subtitle: p.subtitle || '',
      description: toRichHtml(p.description || ''),
      category_id: p.category_id || undefined,
      product_type: p.product_type || 'physical',
      shelf_life: p.shelf_life || 0,
      storage_condition: p.storage_condition || '',
      sale_price: p.sale_price || 0,
      stock: p.stock || 0,
    })
    let specs = p.spec_json
    if (typeof specs === 'string') {
      try {
        specs = JSON.parse(specs)
      } catch {
        specs = []
      }
    }
    form.spec_json = Array.isArray(specs) ? specs : []
    form.skus = (d.skus || []).map((s) => ({
      ...s,
      spec_values: typeof s.spec_values === 'string' ? JSON.parse(s.spec_values || '{}') : s.spec_values,
    }))
    form.images = (d.images || []).map((img, i) => ({
      url: img.url,
      typ: img.typ || 'gallery',
      sort: img.sort ?? i,
    }))
    form.attrs = d.attrs || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function save(status) {
  if (!form.name?.trim()) {
    ElMessage.warning('请填写商品名称')
    return
  }
  if (!form.category_id) {
    ElMessage.warning('请选择商品分类')
    return
  }
  if (!form.skus.length) {
    // 无规格时生成默认 SKU
    form.skus = [
      {
        id: 0,
        spec_values: {},
        sale_price: form.sale_price,
        stock: form.stock,
        stock_warn: 10,
        status: 'enabled',
      },
    ]
  }
  const payload = {
    name: form.name,
    subtitle: form.subtitle,
    description: form.description,
    category_id: form.category_id,
    product_type: form.product_type,
    shelf_life: form.shelf_life,
    storage_condition: form.storage_condition,
    status,
    sale_price: form.sale_price,
    stock: form.stock,
    main_image: form.images.find((i) => i.typ === 'main')?.url || form.images[0]?.url || '',
    spec_json: form.spec_json,
    skus: form.skus,
    images: form.images,
    attrs: form.attrs,
    tag_ids: form.tag_ids,
  }
  loading.value = true
  try {
    if (id) await updateProduct(id, payload)
    else await createProduct(payload)
    ElMessage.success('保存成功')
    router.push('/merchant/products')
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  try {
    const res = await listProductCategories()
    catTree.value = buildCategoryTree(pickList(res))
  } catch (_) {
    catTree.value = []
  }
}

onMounted(async () => {
  await loadCategories()
  await load()
})
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.form { max-width: 960px; background: #fff; padding: 20px; border-radius: 8px; }
.form :deep(.el-form-item__content) { flex: 1; min-width: 0; }
.sku-block { margin: 16px 0 24px; }
</style>
