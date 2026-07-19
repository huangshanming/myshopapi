<template>
  <div>
    <div class="toolbar">
      <h2>全站商品</h2>
      <el-select
        v-model="q.shop_id"
        filterable
        remote
        clearable
        reserve-keyword
        placeholder="商家"
        :remote-method="searchShops"
        :loading="shopLoading"
        style="width: 200px"
        @change="load"
      >
        <el-option v-for="s in shops" :key="s.id" :label="`${s.name} (#${s.id})`" :value="s.id" />
      </el-select>
      <el-input v-model="q.name" placeholder="商品名称" clearable style="width: 140px" @keyup.enter="load" />
      <el-input v-model="q.product_no" placeholder="货号" clearable style="width: 120px" @keyup.enter="load" />
      <el-select v-model="q.status" clearable placeholder="状态" style="width: 110px" @change="load">
        <el-option label="在售" value="on_sale" />
        <el-option label="下架" value="off_sale" />
        <el-option label="草稿" value="draft" />
      </el-select>
      <el-select v-model="q.product_type" clearable placeholder="类型" style="width: 110px" @change="load">
        <el-option label="实物" value="physical" />
        <el-option label="生鲜" value="fresh" />
        <el-option label="虚拟" value="virtual" />
      </el-select>
      <el-tree-select
        v-model="q.category_id"
        :data="catTree"
        :props="{ label: 'name', value: 'id', children: 'children' }"
        check-strictly
        filterable
        clearable
        placeholder="类目"
        style="width: 180px"
        @change="load"
      />
      <el-date-picker
        v-model="createdRange"
        type="daterange"
        value-format="YYYY-MM-DD"
        start-placeholder="创建起"
        end-placeholder="创建止"
        style="width: 240px"
      />
      <el-date-picker
        v-model="publishRange"
        type="daterange"
        value-format="YYYY-MM-DD"
        start-placeholder="上架起"
        end-placeholder="上架止"
        style="width: 240px"
      />
      <el-select v-model="q.order_by" clearable placeholder="排序" style="width: 140px" @change="load">
        <el-option label="收藏人数↓" value="collect_desc" />
        <el-option label="销量↓" value="sold_desc" />
        <el-option label="售价↓" value="sale_price_desc" />
      </el-select>
      <el-button type="primary" @click="load">查询</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="封面" width="72">
        <template #default="{ row }">
          <el-image v-if="row.main_image" :src="row.main_image" style="width: 48px; height: 48px" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column prop="shop_id" label="商家" width="80" />
      <el-table-column prop="product_no" label="货号" width="110" />
      <el-table-column prop="sale_price" label="售价" width="90" />
      <el-table-column prop="stock" label="库存" width="80" />
      <el-table-column prop="collect_count" label="收藏人数" width="100" />
      <el-table-column prop="avg_rating" label="均分" width="70" />
      <el-table-column prop="review_count" label="评价数" width="80" />
      <el-table-column prop="good_rate" label="好评率%" width="90" />
      <el-table-column prop="status" label="状态" width="90" />
      <el-table-column prop="product_type" label="类型" width="90" />
      <el-table-column prop="created_at" label="创建时间" width="160" />
      <el-table-column prop="publish_time" label="上架时间" width="160" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'on_sale'"
            v-permission="'business:product:off_sale'"
            link
            type="warning"
            @click="openRemark(row, 'off_sale')"
          >下架</el-button>
          <el-button
            v-permission="'business:product:delete'"
            link
            type="danger"
            @click="openRemark(row, 'delete')"
          >删除</el-button>
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

    <el-dialog v-model="remarkVisible" :title="remarkAction === 'delete' ? '删除并通知商家' : '下架并通知商家'" width="480px">
      <el-form label-width="80px">
        <el-form-item label="商品">
          <span>{{ remarkRow?.name }} (#{{ remarkRow?.id }})</span>
        </el-form-item>
        <el-form-item label="备注" required>
          <el-input v-model="remark" type="textarea" :rows="3" placeholder="将展示给商家，请说明原因" maxlength="500" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="remarkVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitRemark">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { fetchShops } from '../../../api/merchant'
import {
  listAdminProducts, forceOffSaleProduct, deleteAdminProduct, listProductCategories,
} from '../../../api/admin-product'
import { pickList } from '../../../utils/list'
import { buildCategoryTree } from '../../../utils/categoryTree'

const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const shops = ref([])
const shopLoading = ref(false)
const catTree = ref([])
const createdRange = ref([])
const publishRange = ref([])
const q = reactive({
  shop_id: undefined,
  name: '',
  product_no: '',
  status: '',
  product_type: '',
  category_id: undefined,
  order_by: '',
})

const remarkVisible = ref(false)
const remarkAction = ref('off_sale')
const remarkRow = ref(null)
const remark = ref('')
const submitting = ref(false)

async function searchShops(keyword) {
  shopLoading.value = true
  try {
    const res = await fetchShops({ name: keyword || undefined, page: 1, page_size: 50 })
    shops.value = res?.list || res?.items || []
  } finally {
    shopLoading.value = false
  }
}

async function loadCats() {
  try {
    const res = await listProductCategories()
    catTree.value = buildCategoryTree(pickList(res))
  } catch (_) {
    catTree.value = []
  }
}

async function load() {
  loading.value = true
  try {
    const res = await listAdminProducts({
      page: page.value,
      page_size: pageSize,
      shop_id: q.shop_id || undefined,
      name: q.name || undefined,
      product_no: q.product_no || undefined,
      status: q.status || undefined,
      product_type: q.product_type || undefined,
      category_id: q.category_id || undefined,
      created_from: createdRange.value?.[0],
      created_to: createdRange.value?.[1],
      publish_from: publishRange.value?.[0],
      publish_to: publishRange.value?.[1],
      order_by: q.order_by || undefined,
    })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openRemark(row, action) {
  remarkRow.value = row
  remarkAction.value = action
  remark.value = ''
  remarkVisible.value = true
}

async function submitRemark() {
  if (!remark.value.trim()) {
    ElMessage.warning('请填写备注')
    return
  }
  submitting.value = true
  try {
    if (remarkAction.value === 'delete') {
      await deleteAdminProduct(remarkRow.value.id, remark.value.trim())
      ElMessage.success('已删除并通知商家')
    } else {
      await forceOffSaleProduct(remarkRow.value.id, remark.value.trim())
      ElMessage.success('已下架并通知商家')
    }
    remarkVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  await searchShops('')
  await loadCats()
  load()
})
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 12px; }
.toolbar h2 { margin: 0 12px 0 0; font-size: 18px; }
.pager { margin-top: 12px; justify-content: flex-end; }
</style>
