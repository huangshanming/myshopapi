<template>
  <el-table :data="modelValue" border size="small">
    <el-table-column label="规格" min-width="160">
      <template #default="{ row }">
        {{ formatSpec(row.spec_values) }}
      </template>
    </el-table-column>
    <el-table-column label="售价" width="130">
      <template #default="{ row }">
        <el-input-number v-model="row.sale_price" :min="0" :precision="2" :controls="false" @change="emitChange" />
      </template>
    </el-table-column>
    <el-table-column label="库存" width="120">
      <template #default="{ row }">
        <el-input-number v-model="row.stock" :min="0" :controls="false" @change="emitChange" />
      </template>
    </el-table-column>
    <el-table-column label="预警" width="100">
      <template #default="{ row }">
        <el-input-number v-model="row.stock_warn" :min="0" :controls="false" @change="emitChange" />
      </template>
    </el-table-column>
    <el-table-column label="状态" width="110">
      <template #default="{ row }">
        <el-select v-model="row.status" @change="emitChange">
          <el-option label="启用" value="enabled" />
          <el-option label="禁用" value="disabled" />
        </el-select>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup>
const props = defineProps({
  modelValue: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:modelValue'])

function emitChange() {
  emit('update:modelValue', [...props.modelValue])
}

function formatSpec(sv) {
  if (!sv) return '默认'
  if (typeof sv === 'string') {
    try {
      sv = JSON.parse(sv)
    } catch {
      return sv
    }
  }
  return Object.entries(sv)
    .map(([k, v]) => `${k}:${v}`)
    .join(' / ')
}
</script>
