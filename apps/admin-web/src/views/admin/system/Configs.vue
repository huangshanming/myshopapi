<template>
  <div>
    <div class="toolbar">
      <h2>系统设置</h2>
      <el-button v-permission="'system:config:edit'" type="primary" @click="save">保存</el-button>
    </div>
    <el-form label-width="120px" v-loading="loading" style="max-width: 520px">
      <el-form-item v-for="item in list" :key="item.config_key" :label="item.remark || item.config_key">
        <el-input v-model="item.config_value" />
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { fetchConfigs, saveConfigs } from '../../../api/system'

const list = ref([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await fetchConfigs()
    list.value = res || []
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

async function save() {
  try {
    await saveConfigs(list.value.map((i) => ({
      config_key: i.config_key,
      config_value: i.config_value,
      remark: i.remark,
    })))
    ElMessage.success('已保存')
  } catch (e) {
    ElMessage.error(e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
