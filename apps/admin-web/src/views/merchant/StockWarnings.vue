<template>
  <div>
    <div class="toolbar">
      <h2>库存预警</h2>
      <el-button @click="load">刷新</el-button>
    </div>
    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="SKU ID" width="90" />
      <el-table-column prop="product_id" label="商品ID" width="90" />
      <el-table-column prop="sku_no" label="SKU编码" width="160" />
      <el-table-column label="规格" min-width="160">
        <template #default="{ row }">{{ formatSpec(row.spec_values) }}</template>
      </el-table-column>
      <el-table-column prop="stock" label="库存" width="80" />
      <el-table-column prop="stock_warn" label="预警值" width="90" />
      <el-table-column prop="sale_price" label="售价" width="90" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-input-number v-model="row._stock" :min="0" size="small" style="width: 100px" />
          <el-button size="small" type="primary" link @click="saveStock(row)">改库存</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { pickList } from '../../utils/list'
import { adjustSkuStock, stockWarnings } from '../../api/merchant-product'

const list = ref([])
const loading = ref(false)

function formatSpec(sv) {
  if (!sv) return '-'
  if (typeof sv === 'string') {
    try {
      sv = JSON.parse(sv)
    } catch {
      return sv
    }
  }
  return Object.entries(sv)
    .map(([k, v]) => `${k}:${v}`)
    .join(' / ')
}

async function load() {
  loading.value = true
  try {
    const res = await stockWarnings({ page: 1, page_size: 100 })
    list.value = pickList(res.data).map((r) => ({ ...r, _stock: r.stock }))
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function saveStock(row) {
  await adjustSkuStock(row.id, { stock: row._stock })
  ElMessage.success('已更新')
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
