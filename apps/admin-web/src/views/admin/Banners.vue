<template>
  <div>
    <div class="toolbar">
      <h2>首页 Banner</h2>
      <el-button v-permission="'marketing:banner:edit'" type="primary" @click="openEdit()">新增</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column label="封面" width="120">
        <template #default="{ row }">
          <el-image v-if="row.image_url" :src="row.image_url" fit="cover" style="width: 96px; height: 48px" />
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="120" />
      <el-table-column label="跳转" width="90">
        <template #default="{ row }">{{ bannerLinkLabel(row.link_type) }}</template>
      </el-table-column>
      <el-table-column label="目标" min-width="160">
        <template #default="{ row }">
          <span v-if="row.link_type === 'none'">-</span>
          <span v-else>{{ row.link_name || ('#' + row.link_id) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" width="70" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">{{ row.status === 'on' ? '上架' : '下架' }}</template>
      </el-table-column>
      <el-table-column label="生效时间" min-width="200">
        <template #default="{ row }">
          <span v-if="!row.start_at && !row.end_at">长期</span>
          <span v-else>{{ row.start_at || '—' }} ~ {{ row.end_at || '—' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'marketing:banner:edit'" link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button v-permission="'marketing:banner:delete'" link type="danger" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pager"
      layout="prev,pager,next,total"
      :total="total"
      v-model:current-page="page"
      :page-size="20"
      @current-change="load"
    />

    <el-dialog v-model="visible" :title="form.id ? '编辑 Banner' : '新增 Banner'" width="560px" destroy-on-close>
      <el-form label-width="90px">
        <el-form-item label="标题"><el-input v-model="form.title" maxlength="100" /></el-form-item>
        <el-form-item label="图片" required>
          <ImageUploader v-model="coverImgs" :upload-fn="uploadBannerImage" />
        </el-form-item>
        <el-form-item label="跳转类型">
          <el-radio-group v-model="form.link_type" @change="onLinkTypeChange">
            <el-radio-button v-for="o in BANNER_LINK_OPTIONS" :key="o.value" :value="o.value">{{ o.label }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.link_type === 'product'" label="商品" required>
          <el-select
            v-model="form.link_id"
            filterable
            remote
            clearable
            :remote-method="searchProducts"
            :loading="productLoading"
            style="width: 100%"
            placeholder="搜索商品名称"
          >
            <el-option v-for="p in products" :key="p.id" :label="`#${p.id} ${p.name}`" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.link_type === 'article'" label="文章" required>
          <el-select
            v-model="form.link_id"
            filterable
            remote
            clearable
            :remote-method="searchArticles"
            :loading="articleLoading"
            style="width: 100%"
            placeholder="搜索文章标题"
          >
            <el-option v-for="a in articles" :key="a.id" :label="`#${a.id} ${a.title}`" :value="a.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width: 140px">
            <el-option label="上架" value="on" />
            <el-option label="下架" value="off" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker v-model="form.start_at" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" clearable placeholder="空=立即生效" />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker v-model="form.end_at" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" clearable placeholder="空=长期有效" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  BANNER_LINK_OPTIONS, bannerLinkLabel, createBanner, deleteBanner, listBanners, updateBanner, uploadBannerImage,
} from '../../api/admin-banner'
import { listAdminProducts } from '../../api/admin-product'
import { listArticles } from '../../api/admin-article'
import ImageUploader from '../../components/merchant/ImageUploader.vue'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const visible = ref(false)
const saving = ref(false)
const products = ref([])
const productLoading = ref(false)
const articles = ref([])
const articleLoading = ref(false)

const form = reactive({
  id: 0,
  title: '',
  image_url: '',
  link_type: 'none',
  link_id: undefined,
  sort: 0,
  status: 'on',
  start_at: '',
  end_at: '',
})

const coverImgs = computed({
  get: () => (form.image_url ? [{ url: form.image_url, typ: 'cover' }] : []),
  set: (arr) => {
    const last = arr?.length ? arr[arr.length - 1] : null
    form.image_url = last?.url || ''
  },
})

async function load() {
  loading.value = true
  try {
    const res = await listBanners({ page: page.value, page_size: 20 })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function onLinkTypeChange() {
  form.link_id = undefined
  products.value = []
  articles.value = []
}

async function searchProducts(keyword) {
  productLoading.value = true
  try {
    const res = await listAdminProducts({ page: 1, page_size: 30, name: keyword || undefined, status: 'on_sale' })
    products.value = res?.list || []
  } finally {
    productLoading.value = false
  }
}

async function searchArticles(keyword) {
  articleLoading.value = true
  try {
    const res = await listArticles({
      page: 1,
      page_size: 30,
      title: keyword || undefined,
      status: 'published',
      audit_status: 'approved',
    })
    articles.value = res?.list || []
  } finally {
    articleLoading.value = false
  }
}

function openEdit(row) {
  if (row) {
    Object.assign(form, {
      id: row.id,
      title: row.title || '',
      image_url: row.image_url || '',
      link_type: row.link_type || 'none',
      link_id: row.link_id || undefined,
      sort: row.sort || 0,
      status: row.status || 'on',
      start_at: row.start_at || '',
      end_at: row.end_at || '',
    })
    if (row.link_type === 'product' && row.link_id) {
      products.value = [{ id: row.link_id, name: row.link_name || `商品 #${row.link_id}` }]
    }
    if (row.link_type === 'article' && row.link_id) {
      articles.value = [{ id: row.link_id, title: row.link_name || `文章 #${row.link_id}` }]
    }
  } else {
    Object.assign(form, {
      id: 0, title: '', image_url: '', link_type: 'none', link_id: undefined,
      sort: 0, status: 'on', start_at: '', end_at: '',
    })
    products.value = []
    articles.value = []
  }
  visible.value = true
}

async function save() {
  if (!form.image_url) {
    ElMessage.warning('请上传图片')
    return
  }
  if (form.link_type !== 'none' && !form.link_id) {
    ElMessage.warning('请选择跳转目标')
    return
  }
  saving.value = true
  try {
    const payload = {
      title: form.title,
      image_url: form.image_url,
      link_type: form.link_type,
      link_id: form.link_type === 'none' ? 0 : form.link_id,
      sort: form.sort,
      status: form.status,
      start_at: form.start_at || '',
      end_at: form.end_at || '',
    }
    if (form.id) await updateBanner(form.id, payload)
    else await createBanner(payload)
    ElMessage.success('已保存')
    visible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

async function onDelete(row) {
  try {
    await ElMessageBox.confirm(`确认删除「${row.title || 'Banner'}」？`, '删除确认')
  } catch {
    return
  }
  try {
    await deleteBanner(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
.toolbar h2 { margin-right: auto; }
.pager { margin-top: 16px; justify-content: flex-end; }
</style>
