<template>
  <div class="page">
    <h2>秒杀规则</h2>
    <el-form v-loading="loading" label-width="140px" style="max-width: 480px">
      <el-form-item label="场次时长(小时)">
        <el-input-number v-model="form.duration_hours" :min="1" :max="720" />
      </el-form-item>
      <el-form-item label="报名费用(元)">
        <el-input-number v-model="form.apply_fee" :min="0" :precision="2" :step="1" />
      </el-form-item>
      <el-form-item label="每店每场上限">
        <el-input-number v-model="form.max_entries_per_shop" :min="1" :max="100" />
      </el-form-item>
      <el-form-item label="规则状态">
        <el-switch v-model="enabled" active-text="启用" inactive-text="停用" />
      </el-form-item>
      <el-form-item>
        <el-button v-permission="'marketing:seckill:rule'" type="primary" @click="save">保存</el-button>
      </el-form-item>
    </el-form>
    <p class="hint">说明：修改时长只影响之后新开的场次；当前进行中的场次结束后会按新规则自动开启下一轮。</p>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getSeckillRule, updateSeckillRule } from '../../api/seckill-wallet'

const loading = ref(false)
const form = reactive({
  duration_hours: 24,
  apply_fee: 10,
  max_entries_per_shop: 5,
  status: 1,
})
const enabled = computed({
  get: () => form.status === 1,
  set: (v) => { form.status = v ? 1 : 0 },
})

async function load() {
  loading.value = true
  try {
    const res = await getSeckillRule()
    Object.assign(form, {
      duration_hours: res?.duration_hours ?? 24,
      apply_fee: Number(res?.apply_fee ?? 10),
      max_entries_per_shop: res?.max_entries_per_shop ?? 5,
      status: res?.status ?? 1,
    })
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function save() {
  try {
    await updateSeckillRule({ ...form })
    ElMessage.success('已保存')
    load()
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.hint { color: #64748b; font-size: 13px; margin-top: 16px; max-width: 560px; }
</style>
