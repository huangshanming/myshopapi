<template>
  <div v-loading="loading" class="page">
    <div class="toolbar">
      <h2>商品详情</h2>
      <el-button @click="$router.back()">返回</el-button>
    </div>
    <template v-if="product">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ product.id }}</el-descriptions-item>
        <el-descriptions-item label="店铺">#{{ product.shop_id }}</el-descriptions-item>
        <el-descriptions-item label="名称" :span="2">{{ product.name }}</el-descriptions-item>
        <el-descriptions-item label="副标题" :span="2">{{ product.subtitle || '-' }}</el-descriptions-item>
        <el-descriptions-item label="货号">{{ product.product_no || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ product.status }}</el-descriptions-item>
        <el-descriptions-item label="售价">{{ product.sale_price }}</el-descriptions-item>
        <el-descriptions-item label="库存">{{ product.stock }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ product.product_type }}</el-descriptions-item>
        <el-descriptions-item label="类目">#{{ product.category_id }}</el-descriptions-item>
      </el-descriptions>
      <div v-if="product.main_image" class="cover">
        <el-image :src="product.main_image" fit="contain" style="max-width: 320px; max-height: 320px" />
      </div>
    </template>
    <el-empty v-else-if="!loading" description="商品不存在或已删除" />
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getProductDetail } from '../../../api/order'

const route = useRoute()
const loading = ref(false)
const product = ref(null)

async function load() {
  const id = Number(route.params.id)
  if (!id) {
    product.value = null
    return
  }
  loading.value = true
  try {
    const res = await getProductDetail(id)
    product.value = res || null
  } catch (e) {
    product.value = null
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(() => route.params.id, load)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
.toolbar h2 { margin: 0; font-size: 18px; }
.cover { margin-top: 16px; }
</style>
