<template>
  <el-cascader
    v-model="inner"
    :options="options"
    :props="cascaderProps"
    clearable
    filterable
    placeholder="请选择省 / 市 / 区"
    style="width: 100%"
    @change="onChange"
  />
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import http from '../../api/http'

const props = defineProps({
  province: { type: String, default: '' },
  city: { type: String, default: '' },
  district: { type: String, default: '' },
  provinceCode: { type: String, default: '' },
  cityCode: { type: String, default: '' },
  districtCode: { type: String, default: '' },
})

const emit = defineEmits([
  'update:province', 'update:city', 'update:district',
  'update:provinceCode', 'update:cityCode', 'update:districtCode',
])

const cascaderProps = { value: 'code', label: 'name', children: 'children', emitPath: true }
const options = ref([])
const inner = ref([])
const flat = ref(new Map())

function walk(nodes, parentPath = []) {
  for (const n of nodes || []) {
    const path = [...parentPath, n]
    flat.value.set(n.code, path)
    if (n.children?.length) walk(n.children, path)
  }
}

function syncFromProps() {
  if (props.districtCode && flat.value.has(props.districtCode)) {
    inner.value = flat.value.get(props.districtCode).map((x) => x.code)
    return
  }
  if (props.province && props.city && props.district) {
    for (const [, path] of flat.value) {
      if (
        path.length === 3
        && path[0].name === props.province
        && path[1].name === props.city
        && path[2].name === props.district
      ) {
        inner.value = path.map((x) => x.code)
        return
      }
    }
  }
  inner.value = []
}

function onChange(codes) {
  if (!codes || codes.length < 3) {
    emit('update:province', '')
    emit('update:city', '')
    emit('update:district', '')
    emit('update:provinceCode', '')
    emit('update:cityCode', '')
    emit('update:districtCode', '')
    return
  }
  const path = flat.value.get(codes[2])
  if (!path || path.length < 3) return
  emit('update:province', path[0].name)
  emit('update:city', path[1].name)
  emit('update:district', path[2].name)
  emit('update:provinceCode', path[0].code)
  emit('update:cityCode', path[1].code)
  emit('update:districtCode', path[2].code)
}

onMounted(async () => {
  try {
    const res = await http.get('/api/v1/regions/tree')
    options.value = res || []
    flat.value = new Map()
    walk(options.value)
    syncFromProps()
  } catch {
    options.value = []
  }
})

watch(
  () => [props.provinceCode, props.cityCode, props.districtCode, props.province, props.city, props.district],
  () => {
    if (options.value.length) syncFromProps()
  },
)
</script>
