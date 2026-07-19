<template>
  <div>
    <div class="toolbar">
      <h2>评论表情</h2>
      <el-button v-permission="'community:article:emoji'" type="primary" @click="openEdit()">新增表情</el-button>
    </div>
    <p class="tip">配置 C 端评论可选表情包：上传图标、设置名称与排序；下架后前台不可见。</p>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column label="图标" width="90">
        <template #default="{ row }">
          <el-image :src="row.image_url" style="width:48px;height:48px" fit="contain" />
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column prop="sort" label="排序" width="90" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status === 1 ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="更新时间" width="180" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'community:article:emoji'" link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button v-permission="'community:article:emoji'" link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pager"
      layout="prev, pager, next, total"
      :total="total"
      v-model:current-page="page"
      :page-size="20"
      @current-change="load"
    />

    <el-dialog v-model="visible" :title="form.id ? '编辑表情' : '新增表情'" width="480px" destroy-on-close>
      <el-form label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="form.name" maxlength="32" placeholder="如：开心" />
        </el-form-item>
        <el-form-item label="图标" required>
          <ImageUploader v-model="imgs" :upload-fn="uploadArticleImage" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
          <span class="hint">数字越小越靠前</span>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio-button :value="1">上架</el-radio-button>
            <el-radio-button :value="0">下架</el-radio-button>
          </el-radio-group>
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
  createCommentEmoji, deleteCommentEmoji, listCommentEmojis, updateCommentEmoji, uploadArticleImage,
} from '../../../api/admin-article'
import ImageUploader from '../../../components/merchant/ImageUploader.vue'

const list = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const visible = ref(false)
const saving = ref(false)
const form = reactive({ id: 0, name: '', image_url: '', sort: 0, status: 1 })

const imgs = computed({
  get: () => (form.image_url ? [{ url: form.image_url }] : []),
  set: (arr) => {
    const last = arr?.length ? arr[arr.length - 1] : null
    form.image_url = last?.url || ''
  },
})

async function load() {
  loading.value = true
  try {
    const res = await listCommentEmojis({ page: page.value, page_size: 20 })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openEdit(row) {
  Object.assign(form, row
    ? { id: row.id, name: row.name, image_url: row.image_url, sort: row.sort, status: row.status }
    : { id: 0, name: '', image_url: '', sort: 0, status: 1 })
  visible.value = true
}

async function save() {
  if (!form.image_url) {
    ElMessage.warning('请上传表情图标')
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.name.trim() || '表情',
      image_url: form.image_url,
      sort: form.sort,
      status: form.status,
    }
    if (form.id) await updateCommentEmoji(form.id, payload)
    else await createCommentEmoji(payload)
    ElMessage.success('已保存')
    visible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`删除表情「${row.name || row.id}」？`, '确认')
  } catch { return }
  try {
    await deleteCommentEmoji(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; }
.tip { color: #64748b; font-size: 13px; margin: 0 0 16px; }
.pager { margin-top: 16px; justify-content: flex-end; }
.hint { margin-left: 8px; color: #94a3b8; font-size: 12px; }
</style>
