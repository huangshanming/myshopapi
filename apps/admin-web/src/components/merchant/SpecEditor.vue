<template>
  <div class="spec-editor">
    <div v-for="(spec, si) in modelValue" :key="si" class="spec-row">
      <el-input v-model="spec.name" placeholder="规格名（如颜色）" style="width: 140px" @change="emitChange" />
      <div class="values">
        <el-tag
          v-for="(v, vi) in spec.values"
          :key="vi"
          closable
          class="tag"
          @close="removeValue(si, vi)"
        >{{ v }}</el-tag>
        <el-input
          v-model="drafts[si]"
          size="small"
          placeholder="回车添加规格值"
          style="width: 140px"
          @keyup.enter="addValue(si)"
        />
      </div>
      <el-button link type="danger" @click="removeSpec(si)">删除规格</el-button>
    </div>
    <el-button type="primary" link @click="addSpec">+ 添加规格项</el-button>
    <el-button type="success" link @click="$emit('generate')">生成 SKU 组合</el-button>
  </div>
</template>

<script setup>
import { reactive, watch } from 'vue'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:modelValue', 'generate'])

const drafts = reactive({})

watch(
  () => props.modelValue.length,
  () => {
    props.modelValue.forEach((_, i) => {
      if (drafts[i] === undefined) drafts[i] = ''
    })
  },
  { immediate: true }
)

function emitChange() {
  emit('update:modelValue', [...props.modelValue])
}

function addSpec() {
  emit('update:modelValue', [...props.modelValue, { name: '', values: [] }])
}

function removeSpec(i) {
  const next = props.modelValue.filter((_, idx) => idx !== i)
  emit('update:modelValue', next)
}

function addValue(si) {
  const val = (drafts[si] || '').trim()
  if (!val) return
  const next = props.modelValue.map((s, i) => {
    if (i !== si) return s
    if (s.values.includes(val)) return s
    return { ...s, values: [...s.values, val] }
  })
  drafts[si] = ''
  emit('update:modelValue', next)
}

function removeValue(si, vi) {
  const next = props.modelValue.map((s, i) => {
    if (i !== si) return s
    return { ...s, values: s.values.filter((_, j) => j !== vi) }
  })
  emit('update:modelValue', next)
}
</script>

<style scoped>
.spec-editor { display: flex; flex-direction: column; gap: 12px; }
.spec-row { display: flex; align-items: flex-start; gap: 12px; flex-wrap: wrap; }
.values { display: flex; flex-wrap: wrap; gap: 6px; align-items: center; flex: 1; }
.tag { margin: 0; }
</style>
