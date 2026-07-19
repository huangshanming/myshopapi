<template>
  <div>
    <div class="toolbar">
      <h2>回收站</h2>
      <el-button @click="$router.push('/merchant/products')">返回列表</el-button>
    </div>
    <div class="batch" v-if="selected.length">
      <el-button type="primary" @click="restore">恢复</el-button>
      <el-button type="danger" @click="purge">永久删除</el-button>
    </div>
    <el-table :data="list" v-loading="loading" @selection-change="(r) => (selected = r)">
      <el-table-column type="selection" width="48" />
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="product_no" label="货号" width="140" />
      <el-table-column prop="sale_price" label="售价" width="90" />
    </el-table>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { pickList } from '../../utils/list'
import { deleteRecycle, listProducts, restoreRecycle } from '../../api/merchant-product'

const list = ref([])
const loading = ref(false)
const selected = ref([])

async function load() {
  loading.value = true
  try {
    const res = await listProducts({ page: 1, page_size: 100, recycle: 1 })
    list.value = pickList(res)
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function restore() {
  const ids = selected.value.map((r) => r.id)
  await restoreRecycle(ids)
  ElMessage.success('已恢复')
  load()
}

async function purge() {
  await ElMessageBox.confirm('永久删除后不可恢复，确认？', '提示', { type: 'warning' })
  const ids = selected.value.map((r) => r.id)
  await deleteRecycle(ids)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.batch { margin: 12px 0; display: flex; gap: 8px; }
</style>
