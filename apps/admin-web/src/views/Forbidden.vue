<template>
  <div class="forbidden">
    <h2>403</h2>
    <p>无权限访问该页面</p>
    <el-button type="primary" @click="$router.push('/admin')">返回</el-button>
  </div>
</template>

<style scoped>
.forbidden { padding: 48px; text-align: center; }
</style>
