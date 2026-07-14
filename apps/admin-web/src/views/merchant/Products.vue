<template>
  <div>
    <div class="toolbar">
      <h2>商品管理</h2>
      <el-button type="primary" @click="openCreate">新建商品</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="sale_price" label="售价" width="100" />
      <el-table-column prop="stock" label="库存" width="80" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button size="small" @click="edit(row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="form.id ? '编辑商品' : '新建商品'" width="480px">
      <el-form label-width="80px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="售价"><el-input-number v-model="form.sale_price" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="库存"><el-input-number v-model="form.stock" :min="0" /></el-form-item>
        <el-form-item label="分类ID"><el-input-number v-model="form.category_id" :min="1" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status">
            <el-option label="上架" value="on_sale" />
            <el-option label="下架" value="off_sale" />
            <el-option label="草稿" value="draft" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import http from '../../api/http'

const list = ref([])
const loading = ref(false)
const visible = ref(false)
const form = reactive({ id: 0, name: '', sale_price: 0, stock: 0, category_id: 1, status: 'on_sale' })

async function load() {
  loading.value = true
  try {
    const res = await http.get('/api/v1/merchant/products', { params: { page: 1, page_size: 50 } })
    list.value = res.data?.list || res.data || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, { id: 0, name: '', sale_price: 0, stock: 0, category_id: 1, status: 'on_sale' })
  visible.value = true
}

function edit(row) {
  Object.assign(form, {
    id: row.id,
    name: row.name,
    sale_price: row.sale_price,
    stock: row.stock,
    category_id: row.category_id || 1,
    status: row.status || 'on_sale',
  })
  visible.value = true
}

async function save() {
  const payload = {
    name: form.name,
    sale_price: form.sale_price,
    stock: form.stock,
    category_id: form.category_id,
    status: form.status,
  }
  if (form.id) {
    await http.put(`/api/v1/merchant/products/${form.id}`, payload)
  } else {
    await http.post('/api/v1/merchant/products', payload)
  }
  ElMessage.success('已保存')
  visible.value = false
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
</style>
