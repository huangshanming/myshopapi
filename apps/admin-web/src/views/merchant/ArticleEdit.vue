<template>
  <div>
    <div class="toolbar">
      <h2>{{ id ? (readonly ? '查看文章' : '编辑文章') : '发布文章' }}</h2>
      <el-button @click="$router.back()">返回</el-button>
      <el-button v-if="!readonly" type="primary" :loading="saving" @click="save">提交审核</el-button>
    </div>
    <el-alert
      v-if="readonly"
      type="info"
      :closable="false"
      style="margin-bottom: 12px"
      :title="readonlyTip"
    />
    <el-form label-width="110px" style="max-width: 900px" :disabled="readonly">
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
          style="width: 280px"
        />
      </el-form-item>
      <el-form-item label="封面">
        <ImageUploader v-model="coverImgs" :upload-fn="uploadMyArticleImage" />
      </el-form-item>
      <el-form-item label="定时发布">
        <el-date-picker v-model="form.schedule_publish_at" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" clearable />
      </el-form-item>
      <el-form-item label="正文">
        <RichTextEditor v-model="form.content" placeholder="文章正文" :upload-fn="uploadMyArticleImage" />
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import RichTextEditor from '../../components/merchant/RichTextEditor.vue'
import ImageUploader from '../../components/merchant/ImageUploader.vue'
import {
  createMyArticle, updateMyArticle, getMyArticle,
  listMyArticleCategories, uploadMyArticleImage,
} from '../../api/merchant-article'

const route = useRoute()
const router = useRouter()
const id = computed(() => route.params.id)
const saving = ref(false)
const catTree = ref([])
const coverImgs = ref([])
const auditStatus = ref('pending')
const rejectReason = ref('')
const form = reactive({
  category_id: undefined,
  title: '',
  content: '',
  schedule_publish_at: '',
})

const readonly = computed(() => !!id.value && auditStatus.value !== 'pending')
const readonlyTip = computed(() => {
  if (auditStatus.value === 'approved') return '已审核通过，内容只读'
  if (auditStatus.value === 'rejected') return `已驳回：${rejectReason.value || '无理由'}，内容只读`
  return '只读'
})

async function load() {
  const cats = await listMyArticleCategories()
  catTree.value = cats || []
  if (!id.value) return
  const res = await getMyArticle(id.value)
  const a = res?.article || {}
  form.title = a.title
  form.content = a.content || ''
  form.category_id = a.category_id || undefined
  form.schedule_publish_at = a.schedule_publish_at || ''
  auditStatus.value = a.audit_status
  rejectReason.value = a.reject_reason || ''
  if (a.cover_url) coverImgs.value = [{ url: a.cover_url, typ: 'main' }]
}

async function save() {
  if (!form.title) return ElMessage.warning('请填写标题')
  saving.value = true
  try {
    const payload = {
      ...form,
      cover_url: coverImgs.value[0]?.url || '',
      image_urls: coverImgs.value.map((x) => x.url),
      submit: true,
    }
    if (id.value) await updateMyArticle(id.value, payload)
    else await createMyArticle(payload)
    ElMessage.success('已提交，等待平台审核')
    router.push('/merchant/articles')
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 16px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
</style>
