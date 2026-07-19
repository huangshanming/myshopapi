<template>
  <div class="single-img">
    <div v-if="modelValue" class="preview">
      <el-image :src="modelValue" fit="cover" style="width: 96px; height: 96px" />
      <el-button class="clear" link type="danger" size="small" @click="clear">清除</el-button>
    </div>
    <el-upload
      v-else
      :show-file-list="false"
      :http-request="doUpload"
      accept="image/*"
    >
      <div class="add">+</div>
    </el-upload>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: { type: String, default: '' },
  uploadFn: { type: Function, required: true },
})
const emit = defineEmits(['update:modelValue'])

async function doUpload({ file }) {
  try {
    const res = await props.uploadFn(file)
    const url = res?.url || res.data
    if (!url) throw new Error('未返回图片地址')
    emit('update:modelValue', url)
  } catch (e) {
    ElMessage.error(e.message || '上传失败')
  }
}

function clear() {
  emit('update:modelValue', '')
}
</script>

<style scoped>
.single-img { display: flex; align-items: flex-start; }
.preview { display: flex; flex-direction: column; gap: 4px; align-items: flex-start; }
.add {
  width: 96px; height: 96px; border: 1px dashed #cbd5e1; display: flex; align-items: center;
  justify-content: center; font-size: 28px; color: #94a3b8; cursor: pointer; background: #fff;
}
</style>
