<template>
  <div>
    <div class="toolbar">
      <h2>积分商城</h2>
      <el-button v-permission="'marketing:points_mall:edit'" type="primary" @click="openEdit()">新建商品</el-button>
    </div>
    <p class="tip">配置可用积分兑换的商品。本期仅总后台管理；C 端兑换接口需另行对接。</p>

    <div class="filters">
      <el-input v-model="keyword" clearable placeholder="搜索名称" style="width: 220px" @keyup.enter="reload" />
      <el-select v-model="status" clearable placeholder="状态" style="width: 120px" @change="reload">
        <el-option label="上架" value="on" />
        <el-option label="下架" value="off" />
      </el-select>
      <el-button @click="reload">查询</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column label="封面" width="90">
        <template #default="{ row }">
          <el-image v-if="row.cover_url" :src="row.cover_url" fit="cover" style="width: 56px; height: 56px; border-radius: 6px" />
          <span v-else class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column prop="points_price" label="积分价" width="90" />
      <el-table-column prop="stock" label="库存" width="80" />
      <el-table-column label="每人限兑" width="100">
        <template #default="{ row }">{{ row.per_user_limit > 0 ? row.per_user_limit : '不限' }}</template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" width="70" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 'on' ? 'success' : 'info'" size="small">
            {{ row.status === 'on' ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'marketing:points_mall:edit'" link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button
            v-permission="'marketing:points_mall:edit'"
            link
            type="primary"
            @click="toggleStatus(row)"
          >{{ row.status === 'on' ? '下架' : '上架' }}</el-button>
          <el-button v-permission="'marketing:points_mall:edit'" link type="danger" @click="onDelete(row)">删除</el-button>
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

    <el-dialog v-model="visible" :title="form.id ? '编辑积分商品' : '新建积分商品'" width="560px" destroy-on-close>
      <el-form label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" maxlength="100" />
        </el-form-item>
        <el-form-item label="封面">
          <SingleImageField v-model="form.cover_url" :upload-fn="uploadPointsProductImage" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="form.description" type="textarea" :rows="3" maxlength="1000" />
        </el-form-item>
        <el-form-item label="积分价" required>
          <el-input-number v-model="form.points_price" :min="0" :max="9999999" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="form.stock" :min="0" :max="9999999" />
        </el-form-item>
        <el-form-item label="每人限兑">
          <el-input-number v-model="form.per_user_limit" :min="0" :max="999" />
          <span class="hint">0 = 不限</span>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="statusOn" active-text="上架" inactive-text="下架" />
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
import SingleImageField from '../../components/admin/SingleImageField.vue'
import {
  createPointsProduct,
  deletePointsProduct,
  listPointsProducts,
  setPointsProductStatus,
  updatePointsProduct,
  uploadPointsProductImage,
} from '../../api/points-mall'

const list = ref([])
const loading = ref(false)
const visible = ref(false)
const saving = ref(false)
const page = ref(1)
const total = ref(0)
const keyword = ref('')
const status = ref('')

const form = reactive({
  id: 0,
  name: '',
  cover_url: '',
  description: '',
  points_price: 0,
  stock: 0,
  per_user_limit: 0,
  sort: 0,
  status: 'off',
})

const statusOn = computed({
  get: () => form.status === 'on',
  set: (v) => { form.status = v ? 'on' : 'off' },
})

function resetForm() {
  Object.assign(form, {
    id: 0, name: '', cover_url: '', description: '',
    points_price: 0, stock: 0, per_user_limit: 0, sort: 0, status: 'off',
  })
}

function openEdit(row) {
  if (!row) {
    resetForm()
  } else {
    Object.assign(form, {
      id: row.id,
      name: row.name,
      cover_url: row.cover_url || '',
      description: row.description || '',
      points_price: row.points_price,
      stock: row.stock,
      per_user_limit: row.per_user_limit,
      sort: row.sort,
      status: row.status || 'off',
    })
  }
  visible.value = true
}

function reload() {
  page.value = 1
  load()
}

async function load() {
  loading.value = true
  try {
    const res = await listPointsProducts({
      page: page.value,
      page_size: 20,
      keyword: keyword.value || undefined,
      status: status.value || undefined,
    })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写名称')
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.name.trim(),
      cover_url: form.cover_url,
      description: form.description,
      points_price: form.points_price,
      stock: form.stock,
      per_user_limit: form.per_user_limit,
      sort: form.sort,
      status: form.status,
    }
    if (form.id) await updatePointsProduct(form.id, payload)
    else await createPointsProduct(payload)
    ElMessage.success('已保存')
    visible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function toggleStatus(row) {
  const next = row.status === 'on' ? 'off' : 'on'
  try {
    await setPointsProductStatus(row.id, next)
    ElMessage.success(next === 'on' ? '已上架' : '已下架')
    load()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function onDelete(row) {
  try {
    await ElMessageBox.confirm(`确认删除「${row.name}」？`, '删除积分商品', { type: 'warning' })
    await deletePointsProduct(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || '删除失败')
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; }
.tip { color: #64748b; font-size: 13px; margin: 0 0 16px; max-width: 720px; line-height: 1.5; }
.filters { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.pager { margin-top: 16px; justify-content: flex-end; }
.hint { margin-left: 8px; color: #94a3b8; font-size: 12px; }
.muted { color: #94a3b8; }
</style>
