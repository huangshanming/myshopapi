<template>
  <div>
    <div class="toolbar">
      <h2>评价管理</h2>
      <el-select v-model="ratingLevel" clearable placeholder="评价等级" style="width: 140px" @change="search">
        <el-option label="好评" value="good" />
        <el-option label="中评" value="mid" />
        <el-option label="差评" value="bad" />
      </el-select>
      <el-button type="primary" @click="search">查询</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="order_no" label="订单号" min-width="150" />
      <el-table-column prop="product_id" label="商品ID" width="90" />
      <el-table-column prop="rating" label="评分" width="70" />
      <el-table-column prop="content" label="内容" min-width="180" show-overflow-tooltip />
      <el-table-column label="图片" width="120">
        <template #default="{ row }">
          <div class="imgs">
            <el-image
              v-for="img in (row.images || []).slice(0, 3)"
              :key="img.id || img.url"
              :src="img.url"
              :preview-src-list="(row.images || []).map((i) => i.url)"
              style="width: 32px; height: 32px"
              fit="cover"
            />
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="merchant_reply" label="商家回复" min-width="140" show-overflow-tooltip />
      <el-table-column prop="created_at" label="时间" width="160" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'product:review:reply'" link type="primary" @click="openReply(row)">回复</el-button>
          <el-button v-permission="'product:review:delete'" link type="danger" @click="onDelete(row)">删除</el-button>
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

    <el-dialog v-model="replyVisible" title="回复评价" width="480px">
      <el-input v-model="replyText" type="textarea" :rows="4" maxlength="500" show-word-limit />
      <template #footer>
        <el-button @click="replyVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitReply">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteReview, listReviews, replyReview } from '../../api/order'

const SCOPE = 'merchant'
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const ratingLevel = ref('')
const replyVisible = ref(false)
const replyText = ref('')
const replyRow = ref(null)
const submitting = ref(false)

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

function openReply(row) {
  replyRow.value = row
  replyText.value = row.merchant_reply || ''
  replyVisible.value = true
}

async function submitReply() {
  submitting.value = true
  try {
    await replyReview(replyRow.value.id, { reply: replyText.value })
    ElMessage.success('已回复')
    replyVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    submitting.value = false
  }
}

async function onDelete(row) {
  try {
    await ElMessageBox.confirm('确认删除该评价？（软删除，前台不再展示）', '提示')
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
.imgs { display: flex; gap: 4px; }
</style>
