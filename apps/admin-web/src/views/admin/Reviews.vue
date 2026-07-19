<template>
  <div>
    <div class="toolbar">
      <h2>评价管理</h2>
      <el-input v-model="shopId" placeholder="店铺ID" clearable style="width: 120px" @keyup.enter="search" />
      <el-select v-model="ratingLevel" clearable placeholder="评价等级" style="width: 140px" @change="search">
        <el-option label="好评" value="good" />
        <el-option label="中评" value="mid" />
        <el-option label="差评" value="bad" />
      </el-select>
      <el-button type="primary" @click="search">查询</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="shop_id" label="店铺" width="80" />
      <el-table-column prop="order_no" label="订单号" min-width="150" />
      <el-table-column prop="product_id" label="商品ID" width="90" />
      <el-table-column prop="user_id" label="用户" width="80" />
      <el-table-column prop="rating" label="评分" width="70" />
      <el-table-column prop="content" label="内容" min-width="180" show-overflow-tooltip />
      <el-table-column prop="merchant_reply" label="商家回复" min-width="140" show-overflow-tooltip />
      <el-table-column prop="created_at" label="时间" width="160" />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'business:review:delete'" link type="danger" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pager"
      layout="prev, pager, next, total"
      :total="total"
      v-model:current-page="page"
      :page-size="pageSize"
      @current-change="load"
    />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteReview, listReviews } from '../../api/order'

const SCOPE = 'admin'
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const ratingLevel = ref('')
const shopId = ref('')

function search() {
  page.value = 1
  load()
}

async function load() {
  loading.value = true
  try {
    const res = await listReviews(SCOPE, {
      page: page.value,
      page_size: pageSize,
      rating_level: ratingLevel.value || undefined,
      shop_id: shopId.value || undefined,
    })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
    list.value = []
  } finally {
    loading.value = false
  }
}

async function onDelete(row) {
  try {
    await ElMessageBox.confirm('确认删除违规评价？', '提示')
    await deleteReview(SCOPE, row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || String(e))
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin-right: auto; }
.pager { margin-top: 12px; }
</style>
