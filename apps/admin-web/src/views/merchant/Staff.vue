<template>
  <div>
    <h2>员工与角色</h2>
    <el-tabs v-model="tab">
      <el-tab-pane label="员工" name="staff">
        <div class="toolbar">
          <el-button type="primary" @click="openStaff('create')">新建店员</el-button>
          <el-button @click="openStaff('bind')">绑定已有账号</el-button>
        </div>
        <el-table :data="staff" v-loading="loading">
          <el-table-column prop="user_id" label="用户ID" width="90" />
          <el-table-column prop="mobile" label="手机号" />
          <el-table-column prop="nickname" label="昵称" />
          <el-table-column prop="role_name" label="角色" />
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="角色" name="roles">
        <div class="toolbar">
          <el-button type="primary" @click="openRole()">新建角色</el-button>
        </div>
        <el-table :data="roles" v-loading="loading">
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="code" label="编码" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="remark" label="备注" />
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button link size="small" @click="openRole(row)">编辑</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog
      v-model="staffVisible"
      :title="staffForm.mode === 'create' ? '新建店员' : '绑定已有账号'"
      width="460px"
      @closed="resetStaffForm"
    >
      <el-form ref="staffFormRef" :model="staffForm" :rules="staffRules" label-width="90px">
        <el-form-item label="手机号" prop="mobile">
          <el-input v-model="staffForm.mobile" maxlength="11" placeholder="11 位手机号" />
        </el-form-item>
        <template v-if="staffForm.mode === 'create'">
          <el-form-item label="昵称" prop="nickname">
            <el-input v-model="staffForm.nickname" placeholder="可选，默认用手机号" />
          </el-form-item>
          <el-form-item label="登录密码" prop="password">
            <el-input v-model="staffForm.password" type="password" show-password placeholder="至少 6 位" />
          </el-form-item>
        </template>
        <el-form-item label="角色" prop="role_id">
          <el-select v-model="staffForm.role_id" style="width: 100%" placeholder="请选择角色">
            <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-alert
          v-if="staffForm.mode === 'create'"
          type="info"
          :closable="false"
          show-icon
          title="若手机号已注册，将提示改用「绑定已有账号」"
        />
      </el-form>
      <template #footer>
        <el-button @click="staffVisible = false">取消</el-button>
        <el-button type="primary" :loading="staffSaving" @click="submitStaff">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="roleVisible" title="角色权限" width="640px">
      <el-form label-width="80px">
        <el-form-item label="编码"><el-input v-model="roleForm.code" :disabled="!!roleForm.id" /></el-form-item>
        <el-form-item label="名称"><el-input v-model="roleForm.name" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="roleForm.remark" /></el-form-item>
        <el-form-item label="菜单权限">
          <div class="perm-tree">
            <div v-for="group in menuGroups" :key="group.id" class="perm-group">
              <div class="group-title">
                <el-checkbox
                  :model-value="isGroupChecked(group)"
                  :indeterminate="isGroupIndeterminate(group)"
                  @change="(v) => toggleGroup(group, v)"
                >{{ group.name }}</el-checkbox>
              </div>
              <el-checkbox-group v-model="roleForm.menu_ids" class="group-items">
                <el-checkbox
                  v-for="m in group.items"
                  :key="m.id"
                  :label="m.id"
                  :value="m.id"
                >{{ m.name }}<span v-if="m.perms" class="perm-code">{{ m.perms }}</span></el-checkbox>
              </el-checkbox-group>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleVisible = false">取消</el-button>
        <el-button type="primary" :loading="roleSaving" @click="saveRole">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  bindShopStaff,
  getRoleMenus,
  listShopMenus,
  listShopRoles,
  listShopStaff,
  saveShopRole,
} from '../../api/merchant-product'

const tab = ref('staff')
const loading = ref(false)
const staff = ref([])
const roles = ref([])
const menuTree = ref([])
const staffVisible = ref(false)
const roleVisible = ref(false)
const staffSaving = ref(false)
const roleSaving = ref(false)
const staffFormRef = ref()

const staffForm = reactive({
  mode: 'create',
  mobile: '',
  nickname: '',
  password: '',
  role_id: null,
})

const staffRules = {
  mobile: [
    { required: true, message: '请填写手机号', trigger: 'blur' },
    { pattern: /^1\d{10}$/, message: '手机号格式不正确', trigger: 'blur' },
  ],
  password: [
    {
      validator: (_r, v, cb) => {
        if (staffForm.mode !== 'create') return cb()
        if (!v || String(v).length < 6) return cb(new Error('密码至少 6 位'))
        cb()
      },
      trigger: 'blur',
    },
  ],
  role_id: [{ required: true, message: '请选择角色', trigger: 'change' }],
}

const roleForm = reactive({ id: 0, code: '', name: '', remark: '', menu_ids: [] })

const menuGroups = computed(() => {
  return (menuTree.value || [])
    .filter((n) => n.type === 'dir' || n.children?.length)
    .map((n) => {
      const items = []
      const walk = (nodes, depth = 0) => {
        ;(nodes || []).forEach((c) => {
          if (c.type === 'menu' || c.type === 'button') {
            items.push({ ...c, _depth: depth })
          }
          if (c.children?.length) walk(c.children, depth + 1)
        })
      }
      walk(n.children || [])
      if (!n.children?.length && (n.type === 'menu' || n.type === 'button')) {
        items.push(n)
      }
      return { id: n.id, name: n.name, items }
    })
    .filter((g) => g.items.length)
})

function groupIds(group) {
  return group.items.map((i) => i.id)
}

function isGroupChecked(group) {
  const ids = groupIds(group)
  return ids.length > 0 && ids.every((id) => roleForm.menu_ids.includes(id))
}

function isGroupIndeterminate(group) {
  const ids = groupIds(group)
  const n = ids.filter((id) => roleForm.menu_ids.includes(id)).length
  return n > 0 && n < ids.length
}

function toggleGroup(group, checked) {
  const ids = groupIds(group)
  if (checked) {
    roleForm.menu_ids = Array.from(new Set([...roleForm.menu_ids, ...ids, group.id]))
  } else {
    const drop = new Set([...ids, group.id])
    roleForm.menu_ids = roleForm.menu_ids.filter((id) => !drop.has(id))
  }
}

function errMsg(e) {
  return e?.message || e?.msg || '操作失败'
}

async function load() {
  loading.value = true
  try {
    const [s, r, m] = await Promise.all([listShopStaff(), listShopRoles(), listShopMenus()])
    staff.value = Array.isArray(s.data) ? s.data : []
    roles.value = Array.isArray(r.data) ? r.data : []
    menuTree.value = Array.isArray(m.data) ? m.data : []
  } catch (e) {
    ElMessage.error(errMsg(e))
  } finally {
    loading.value = false
  }
}

function resetStaffForm() {
  Object.assign(staffForm, { mode: 'create', mobile: '', nickname: '', password: '', role_id: null })
  staffFormRef.value?.clearValidate?.()
}

function openStaff(mode) {
  resetStaffForm()
  staffForm.mode = mode
  const ops = roles.value.find((r) => r.code === 'shop_ops') || roles.value.find((r) => r.code !== 'shop_owner')
  staffForm.role_id = ops?.id || roles.value[0]?.id || null
  staffVisible.value = true
}

async function submitStaff() {
  try {
    await staffFormRef.value?.validate?.()
  } catch {
    return
  }
  staffSaving.value = true
  try {
    const payload = {
      mode: staffForm.mode,
      mobile: staffForm.mobile.trim(),
      role_id: staffForm.role_id,
      nickname: staffForm.nickname.trim(),
      password: staffForm.password,
    }
    const res = await bindShopStaff(payload)
    ElMessage.success(res.msg || (staffForm.mode === 'create' ? '已创建并绑定' : '已绑定'))
    staffVisible.value = false
    load()
  } catch (e) {
    const msg = errMsg(e)
    ElMessage.error(msg)
    // 新建时若已存在，引导切到绑定
    if (staffForm.mode === 'create' && /已注册|已存在/.test(msg)) {
      try {
        await ElMessageBox.confirm(`${msg}，是否改为绑定该账号？`, '提示', {
          confirmButtonText: '去绑定',
          cancelButtonText: '取消',
          type: 'warning',
        })
        staffForm.mode = 'bind'
        staffForm.password = ''
      } catch {
        /* 用户取消 */
      }
    }
  } finally {
    staffSaving.value = false
  }
}

async function openRole(row) {
  if (row) {
    Object.assign(roleForm, {
      id: row.id,
      code: row.code,
      name: row.name,
      remark: row.remark || '',
      menu_ids: [],
    })
    try {
      const res = await getRoleMenus(row.id)
      roleForm.menu_ids = Array.isArray(res.data) ? res.data.map(Number) : []
    } catch (e) {
      roleForm.menu_ids = []
      ElMessage.error(errMsg(e))
    }
  } else {
    Object.assign(roleForm, { id: 0, code: 'shop_ops', name: '运营', remark: '', menu_ids: [] })
  }
  roleVisible.value = true
}

async function saveRole() {
  if (!roleForm.name?.trim()) {
    ElMessage.warning('请填写角色名称')
    return
  }
  roleSaving.value = true
  try {
    const dirIds = menuGroups.value
      .filter((g) => groupIds(g).some((id) => roleForm.menu_ids.includes(id)))
      .map((g) => g.id)
    const menu_ids = Array.from(new Set([...roleForm.menu_ids, ...dirIds]))
    await saveShopRole(
      {
        code: roleForm.code,
        name: roleForm.name,
        remark: roleForm.remark,
        menu_ids,
      },
      roleForm.id || undefined
    )
    ElMessage.success('已保存')
    roleVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(errMsg(e))
  } finally {
    roleSaving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { margin: 12px 0; display: flex; gap: 8px; }
.perm-tree { width: 100%; max-height: 360px; overflow: auto; border: 1px solid #e5e7eb; border-radius: 6px; padding: 8px 12px; }
.perm-group + .perm-group { margin-top: 12px; padding-top: 12px; border-top: 1px dashed #e5e7eb; }
.group-title { font-weight: 600; margin-bottom: 6px; }
.group-items { display: flex; flex-direction: column; gap: 4px; padding-left: 22px; }
.perm-code { margin-left: 8px; color: #94a3b8; font-size: 12px; }
</style>
