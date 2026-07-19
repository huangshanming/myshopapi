<template>
  <div class="rich-wrap">
    <Toolbar
      v-if="editorRef"
      class="rich-toolbar"
      :editor="editorRef"
      :default-config="toolbarConfig"
      mode="default"
    />
    <Editor
      class="rich-editor"
      :model-value="modelValue"
      :default-config="editorConfig"
      mode="default"
      :style="{ height }"
      @onCreated="onCreated"
      @onChange="onChange"
    />
  </div>
</template>

<script setup>
import { onBeforeUnmount, shallowRef, watch } from 'vue'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import '@wangeditor/editor/dist/css/style.css'
import { ElMessage } from 'element-plus'
import { uploadImage as defaultUpload } from '../../api/merchant-product'

const props = defineProps({
  modelValue: { type: String, default: '' },
  height: { type: String, default: '360px' },
  placeholder: { type: String, default: '请输入商品详情，支持标题、列表、图片等' },
  uploadFn: { type: Function, default: null },
})
const emit = defineEmits(['update:modelValue'])

const editorRef = shallowRef()
let syncing = false

const toolbarConfig = {
  excludeKeys: ['fullScreen', 'group-video', 'insertTable'],
}

const editorConfig = {
  placeholder: props.placeholder,
  MENU_CONF: {
    uploadImage: {
      async customUpload(file, insertFn) {
        try {
          const up = props.uploadFn || defaultUpload
          const res = await up(file)
          const url = res?.url || res.data
          if (!url) throw new Error('上传未返回地址')
          insertFn(url, file.name || '图片', url)
        } catch (e) {
          ElMessage.error(e.message || '图片上传失败')
        }
      },
    },
  },
}

function onCreated(editor) {
  editorRef.value = editor
  if (props.modelValue) {
    syncing = true
    editor.setHtml(props.modelValue)
    syncing = false
  }
}

function onChange(editor) {
  if (syncing) return
  const html = editor.getHtml()
  const empty = !html || html === '<p><br></p>' || html === '<p></p>'
  emit('update:modelValue', empty ? '' : html)
}

watch(
  () => props.modelValue,
  (val) => {
    const ed = editorRef.value
    if (!ed) return
    const next = val || ''
    const cur = ed.getHtml()
    const curNorm = !cur || cur === '<p><br></p>' || cur === '<p></p>' ? '' : cur
    if (next === curNorm) return
    syncing = true
    ed.setHtml(next || '<p><br></p>')
    syncing = false
  }
)

onBeforeUnmount(() => {
  const ed = editorRef.value
  if (ed) ed.destroy()
})
</script>

<style scoped>
.rich-wrap {
  width: 100%;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  overflow: hidden;
  background: #fff;
}
.rich-toolbar {
  border-bottom: 1px solid #e5e7eb;
}
.rich-editor {
  overflow-y: hidden;
}
.rich-wrap :deep(.w-e-text-container) {
  background: #fff;
}
</style>
