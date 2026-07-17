<template>
  <div class="uploader">
    <div
      v-for="(img, i) in modelValue"
      :key="img.url + i"
      class="thumb"
      draggable="true"
      @dragstart="onDragStart(i)"
      @dragover.prevent
      @drop="onDrop(i)"
    >
      <el-image :src="img.url" fit="cover" style="width: 88px; height: 88px" />
      <div class="actions">
        <el-tag size="small">{{ img.typ || 'gallery' }}</el-tag>
        <el-button link type="danger" size="small" @click="remove(i)">删</el-button>
      </div>
    </div>
    <el-upload
      :show-file-list="false"
      :http-request="doUpload"
      accept="image/*"
      multiple
    >
      <div class="add">+</div>
    </el-upload>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { uploadImage as defaultUpload } from '../../api/merchant-product'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  uploadFn: { type: Function, default: null },
})
const emit = defineEmits(['update:modelValue'])

const dragFrom = ref(-1)

async function doUpload({ file }) {
  try {
    const up = props.uploadFn || defaultUpload
    const res = await up(file)
    const url = res.data?.url || res.data
    const next = [...props.modelValue, { url, typ: props.modelValue.length ? 'gallery' : 'main', sort: props.modelValue.length }]
    if (next.length === 1) next[0].typ = 'main'
    emit('update:modelValue', next)
  } catch (e) {
    ElMessage.error(e.message || '上传失败')
  }
}

function remove(i) {
  const next = props.modelValue.filter((_, idx) => idx !== i).map((img, idx) => ({ ...img, sort: idx }))
  if (next.length && !next.some((x) => x.typ === 'main')) next[0].typ = 'main'
  emit('update:modelValue', next)
}

function onDragStart(i) {
  dragFrom.value = i
}

function onDrop(i) {
  const from = dragFrom.value
  if (from < 0 || from === i) return
  const next = [...props.modelValue]
  const [item] = next.splice(from, 1)
  next.splice(i, 0, item)
  emit(
    'update:modelValue',
    next.map((img, idx) => ({ ...img, sort: idx }))
  )
  dragFrom.value = -1
}
</script>

<style scoped>
.uploader { display: flex; flex-wrap: wrap; gap: 10px; align-items: flex-start; }
.thumb { position: relative; border: 1px solid #e5e7eb; padding: 4px; background: #fff; cursor: grab; }
.actions { display: flex; justify-content: space-between; margin-top: 4px; }
.add {
  width: 88px; height: 88px; border: 1px dashed #cbd5e1; display: flex; align-items: center;
  justify-content: center; font-size: 28px; color: #94a3b8; cursor: pointer;
}
</style>
