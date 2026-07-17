<template>
  <div>
    <div class="toolbar">
      <h2>{{ id ? '编辑文章' : '新增文章' }}</h2>
      <el-button @click="$router.back()">返回</el-button>
      <el-button type="primary" :loading="saving" @click="save">保存</el-button>
    </div>
    <el-form label-width="110px" style="max-width: 900px">
      <el-form-item label="商家ID" required>
        <el-input v-model.number="form.shop_id" type="number" placeholder="必填" style="width: 200px" />
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
import {
  createArticle, updateArticle, getArticle, listArticleCategories, uploadArticleImage,
} from '../../../api/admin-article'

const route = useRoute()
const router = useRouter()
const id = computed(() => route.params.id)
const saving = ref(false)
const catTree = ref([])
const coverImgs = ref([])
const form = reactive({
  shop_id: 0,
  category_id: undefined,
  title: '',
  content: '',
  schedule_publish_at: '',
  is_top: 0,
})

function upFn(file) {
  if (!form.shop_id) {
    return Promise.reject(new Error('请先填写商家ID'))
  }
  return uploadArticleImage(file, form.shop_id)
}

async function loadCats() {
  const res = await listArticleCategories()
  catTree.value = res.data || []
}

async function loadDetail() {
  if (!id.value) return
  const res = await getArticle(id.value)
  const a = res.data?.article || {}
  form.shop_id = a.shop_id
  form.category_id = a.category_id || undefined
  form.title = a.title
  form.content = a.content || ''
  form.schedule_publish_at = a.schedule_publish_at || ''
  form.is_top = a.is_top || 0
  if (a.cover_url) coverImgs.value = [{ url: a.cover_url, typ: 'main' }]
}

async function save() {
  if (!form.shop_id || !form.title) {
    ElMessage.warning('请填写商家与标题')
    return
  }
  saving.value = true
  try {
    const payload = {
      ...form,
      cover_url: coverImgs.value[0]?.url || '',
      image_urls: coverImgs.value.map((x) => x.url),
      is_top: form.is_top,
    }
    if (id.value) await updateArticle(id.value, payload)
    else await createArticle(payload)
    ElMessage.success('已保存')
    router.push('/admin/articles')
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await loadCats()
  await loadDetail()
})
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 16px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
</style>
