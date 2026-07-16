<template>
  <div>
    <div class="toolbar">
      <h2>店铺管理</h2>
      <el-input v-model="name" placeholder="店铺名称" clearable style="width: 180px" @clear="load" @keyup.enter="load" />
      <el-select v-model="status" clearable placeholder="状态" style="width: 120px" @change="load">
        <el-option label="已通过" value="approved" />
        <el-option label="已停用" value="disabled" />
        <el-option label="待审核" value="pending" />
      </el-select>
      <el-button @click="load">查询</el-button>
      <el-button v-permission="'business:shop:add'" type="primary" @click="openCreate">新建店铺</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column prop="category" label="类目" width="100" />
      <el-table-column prop="contact_name" label="联系人" width="100" />
      <el-table-column prop="contact_phone" label="电话" width="120" />
      <el-table-column prop="owner_user_id" label="店主UID" width="90" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button v-permission="'business:shop:edit'" link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button v-permission="'business:shop:reset'" link @click="openReset(row)">重置密码</el-button>
          <el-button
            v-if="row.status !== 'disabled'"
            size="small"
            type="warning"
            @click="disable(row)"
          >停用</el-button>
          <el-button v-else size="small" type="success" @click="enable(row)">启用</el-button>
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

    <el-dialog v-model="formVisible" :title="isCreate ? '新建店铺' : '编辑店铺'" width="640px" destroy-on-close>
      <el-form label-width="100px">
        <template v-if="isCreate">
          <el-divider content-position="left">店主账号</el-divider>
          <el-form-item label="手机号" required>
            <el-input v-model="form.owner_mobile" maxlength="11" placeholder="已有用户则绑定，否则新建" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="form.owner_password" type="password" show-password placeholder="新建用户时必填" />
          </el-form-item>
          <el-form-item label="昵称">
            <el-input v-model="form.owner_nickname" />
          </el-form-item>
        </template>
        <el-divider content-position="left">店铺资料</el-divider>
        <el-form-item label="店铺名称" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="经营类目">
          <el-input v-model="form.category" />
        </el-form-item>
        <el-form-item label="Logo URL">
          <el-input v-model="form.logo" />
        </el-form-item>
        <el-form-item label="联系人">
          <el-input v-model="form.contact_name" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="form.contact_phone" maxlength="11" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="省 / 市 / 区">
          <div class="addr-row">
            <el-input v-model="form.province" placeholder="省" />
            <el-input v-model="form.city" placeholder="市" />
            <el-input v-model="form.district" placeholder="区" />
          </div>
        </el-form-item>
        <el-form-item label="详细地址">
          <el-input v-model="form.address" />
        </el-form-item>
        <el-form-item label="执照号">
          <el-input v-model="form.business_license_no" />
        </el-form-item>
        <el-form-item label="法人">
          <el-input v-model="form.legal_person" />
        </el-form-item>
        <el-form-item label="执照图 URL">
          <el-input v-model="form.license_image" />
        </el-form-item>
        <el-form-item label="门头图 URL">
          <el-input v-model="form.storefront_image" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="saveForm">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="pwdVisible" title="重置店主密码" width="420px">
      <el-input v-model="newPassword" type="password" show-password placeholder="新密码" />
      <template #footer>
        <el-button @click="pwdVisible = false">取消</el-button>
        <el-button type="primary" @click="savePwd">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  createShop, disableShop, enableShop, fetchShops, resetShopOwnerPassword, updateShop,
} from '../../api/merchant'

const emptyForm = () => ({
  name: '',
  logo: '',
  contact_name: '',
  contact_phone: '',
  description: '',
  category: '',
  province: '',
  city: '',
  district: '',
  address: '',
  business_license_no: '',
  legal_person: '',
  license_image: '',
  storefront_image: '',
  owner_mobile: '',
  owner_password: '123456',
  owner_nickname: '',
})

const list = ref([])
const loading = ref(false)
const name = ref('')
const status = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)
const formVisible = ref(false)
const pwdVisible = ref(false)
const isCreate = ref(true)
const form = ref(emptyForm())
const currentId = ref(0)
const newPassword = ref('123456')

async function load() {
  loading.value = true
  try {
    const res = await fetchShops({
      page: page.value,
      page_size: pageSize,
      name: name.value || undefined,
      status: status.value || undefined,
    })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

function openCreate() {
  isCreate.value = true
  form.value = emptyForm()
  formVisible.value = true
}

function openEdit(row) {
  isCreate.value = false
  currentId.value = row.id
  form.value = {
    ...emptyForm(),
    name: row.name || '',
    logo: row.logo || '',
    contact_name: row.contact_name || '',
    contact_phone: row.contact_phone || '',
    description: row.description || '',
    category: row.category || '',
    province: row.province || '',
    city: row.city || '',
    district: row.district || '',
    address: row.address || '',
    business_license_no: row.business_license_no || '',
    legal_person: row.legal_person || '',
    license_image: row.license_image || '',
    storefront_image: row.storefront_image || '',
  }
  formVisible.value = true
}

async function saveForm() {
  try {
    if (isCreate.value) {
      await createShop(form.value)
      ElMessage.success('已创建')
    } else {
      const { owner_mobile, owner_password, owner_nickname, ...payload } = form.value
      await updateShop(currentId.value, payload)
      ElMessage.success('已更新')
    }
    formVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

function openReset(row) {
  currentId.value = row.id
  newPassword.value = '123456'
  pwdVisible.value = true
}

async function savePwd() {
  try {
    await resetShopOwnerPassword(currentId.value, newPassword.value)
    ElMessage.success('已重置')
    pwdVisible.value = false
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function disable(row) {
  try {
    await disableShop(row.id, '')
    ElMessage.success('已停用')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

async function enable(row) {
  try {
    await enableShop(row.id)
    ElMessage.success('已启用')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; align-items: center; margin-bottom: 12px; flex-wrap: wrap; }
.toolbar h2 { margin-right: auto; }
.pager { margin-top: 12px; }
.addr-row { display: flex; gap: 8px; width: 100%; }
</style>
