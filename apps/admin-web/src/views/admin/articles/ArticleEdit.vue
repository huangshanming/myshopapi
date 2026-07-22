<template>
  <div>
    <div class="toolbar">
      <h2>{{ id ? '编辑文章' : '发布文章' }}</h2>
      <el-button @click="$router.back()">返回</el-button>
      <el-button type="primary" :loading="saving" @click="save">保存并发布</el-button>
    </div>
    <el-alert
      type="info"
      :closable="false"
      show-icon
      class="hint"
      title="超管发布跳过商家审核，可直接上架或定时发布。归属选「平台官方」为平台自有文章；也可代商家发布。"
    />
    <el-form label-width="110px" style="max-width: 900px">
      <el-form-item label="归属" required>
        <el-select
          v-model="form.shop_id"
          filterable
          remote
          clearable
          reserve-keyword
          placeholder="平台官方 / 商家"
          :remote-method="searchShops"
          :loading="shopLoading"
          style="width: 320px"
        >
          <el-option label="平台官方" :value="0" />
          <el-option v-for="s in shops" :key="s.id" :label="`${s.name} (#${s.id})`" :value="s.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="标题" required>
        <el-input v-model="form.title" maxlength="200" />
      </el-form-item>
      <el-form-item label="分类">
        <el-tree-select
          v-model="form.category_id"
          :data="catTree"
          :props="{ label: 'name', value: 'id', children: 'children' }"
          check-strictly
          clearable
          placeholder="选择分类"
          style="width: 280px"
        />
      </el-form-item>
      <el-form-item label="封面">
        <ImageUploader v-model="coverImgs" :upload-fn="upFn" />
      </el-form-item>
      <el-form-item label="定时发布">
        <el-date-picker v-model="form.schedule_publish_at" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" clearable />
      </el-form-item>
      <el-form-item label="置顶">
        <el-switch v-model="form.is_top" :active-value="1" :inactive-value="0" />
      </el-form-item>
      <el-form-item label="正文">
        <RichTextEditor v-model="form.content" placeholder="文章正文" :upload-fn="upFn" />
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import RichTextEditor from '../../../components/merchant/RichTextEditor.vue'
import ImageUploader from '../../../components/merchant/ImageUploader.vue'
import { fetchShops } from '../../../api/merchant'
import {
  createArticle, updateArticle, getArticle, listArticleCategories, uploadArticleImage,
} from '../../../api/admin-article'
import { pickList } from '../../../utils/list'

const route = useRoute()
const router = useRouter()
const id = computed(() => route.params.id)
const saving = ref(false)
const catTree = ref([])
const coverImgs = ref([])
const shops = ref([])
const shopLoading = ref(false)
const form = reactive({
  shop_id: 0,
  category_id: undefined,
  title: '',
  content: '',
  schedule_publish_at: '',
  is_top: 0,
})

function upFn(file) {
  return uploadArticleImage(file, form.shop_id || 0)
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

async function loadCats() {
  const res = await listArticleCategories()
  catTree.value = pickList(res)
}

async function loadDetail() {
  if (!id.value) return
  const res = await getArticle(id.value)
  const a = res?.article || {}
  form.shop_id = a.shop_id || 0
  form.category_id = a.category_id || undefined
  form.title = a.title
  form.content = a.content || ''
  form.schedule_publish_at = a.schedule_publish_at || ''
  form.is_top = a.is_top || 0
  if (a.cover_url) coverImgs.value = [{ url: a.cover_url, typ: 'main' }]
  if (form.shop_id) {
    await searchShops('')
    if (!shops.value.some((s) => s.id === form.shop_id)) {
      shops.value = [{ id: form.shop_id, name: `商家#${form.shop_id}` }, ...shops.value]
    }
  }
}

async function save() {
  if (form.shop_id === undefined || form.shop_id === null || form.shop_id === '') {
    ElMessage.warning('请选择归属')
    return
  }
  if (!form.title) {
    ElMessage.warning('请填写标题')
    return
  }
  saving.value = true
  try {
    const payload = {
      ...form,
      shop_id: Number(form.shop_id) || 0,
      cover_url: coverImgs.value[0]?.url || '',
      image_urls: coverImgs.value.map((x) => x.url),
      is_top: form.is_top,
    }
    if (id.value) await updateArticle(id.value, payload)
    else await createArticle(payload)
    ElMessage.success('已保存并发布')
    router.push('/admin/articles')
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await searchShops('')
  await loadCats()
  await loadDetail()
})
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.hint { margin-bottom: 16px; max-width: 900px; }
</style>
