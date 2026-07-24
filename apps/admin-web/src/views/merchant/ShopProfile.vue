<template>
  <div class="page" v-loading="loading">
    <h2>店铺资料</h2>
    <el-form label-width="110px" class="form">
      <el-form-item label="店铺名称">
        <el-input v-model="form.name" maxlength="50" />
      </el-form-item>
      <el-form-item label="类目">
        <el-input v-model="form.category" maxlength="30" />
      </el-form-item>
      <el-form-item label="联系人">
        <el-input v-model="form.contact_name" maxlength="30" />
      </el-form-item>
      <el-form-item label="联系电话">
        <el-input v-model="form.contact_phone" maxlength="11" />
      </el-form-item>
      <el-form-item label="简介">
        <el-input v-model="form.description" type="textarea" :rows="3" maxlength="500" />
      </el-form-item>
      <el-form-item label="所在地区">
        <RegionCascader
          v-model:province="form.province"
          v-model:city="form.city"
          v-model:district="form.district"
        />
      </el-form-item>
      <el-form-item label="详细地址">
        <el-input v-model="form.address" maxlength="200" />
      </el-form-item>
      <el-form-item label="门头图">
        <el-upload :show-file-list="false" :http-request="uploadStorefront" accept="image/*">
          <img v-if="form.storefront_image" :src="form.storefront_image" class="cover" />
          <el-button v-else>上传门头图</el-button>
        </el-upload>
      </el-form-item>
      <el-form-item label="门店图片">
        <div class="imgs">
          <div v-for="(u, i) in form.images" :key="u + i" class="img-item">
            <img :src="u" />
            <el-button link type="danger" @click="form.images.splice(i, 1)">删除</el-button>
          </div>
          <el-upload :show-file-list="false" :http-request="uploadGallery" accept="image/*">
            <el-button>添加图片</el-button>
          </el-upload>
        </div>
      </el-form-item>
      <el-form-item label="地图选点">
        <TencentMapPicker
          :latitude="form.latitude"
          :longitude="form.longitude"
          @pick="onMapPick"
        />
        <div class="coord" v-if="form.latitude">
          坐标：{{ form.latitude.toFixed(6) }}, {{ form.longitude.toFixed(6) }}
        </div>
      </el-form-item>
      <el-form-item label="本地商家">
        <el-switch v-model="localOn" active-text="展示在 C 端本地商家专栏" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import RegionCascader from '../../components/admin/RegionCascader.vue'
import TencentMapPicker from '../../components/merchant/TencentMapPicker.vue'
import { uploadImage } from '../../api/merchant-product'
import { listMyShops, reverseGeocode, updateMyShop } from '../../api/shop-profile'
import { useAuthStore } from '../../stores/auth'
import { pickList } from '../../utils/list'

const auth = useAuthStore()
const loading = ref(false)
const saving = ref(false)
const shopId = ref(0)
const form = reactive({
  name: '',
  logo: '',
  contact_name: '',
  contact_phone: '',
  description: '',
  category: '',
  province: '',
  city: '',
  district: '',
  address: '',
  storefront_image: '',
  latitude: 0,
  longitude: 0,
  local_enabled: 0,
  images: [],
})

const localOn = computed({
  get: () => form.local_enabled === 1,
  set: (v) => { form.local_enabled = v ? 1 : 0 },
})

function unwrapUrl(res) {
  return res?.url || res?.data?.url || res?.path || (typeof res === 'string' ? res : '')
}

async function uploadStorefront({ file }) {
  try {
    const res = await uploadImage(file)
    form.storefront_image = unwrapUrl(res)
    if (!form.logo) form.logo = form.storefront_image
  } catch (e) {
    ElMessage.error(e.message || '上传失败')
  }
}

async function uploadGallery({ file }) {
  try {
    const res = await uploadImage(file)
    const url = unwrapUrl(res)
    if (url) form.images.push(url)
  } catch (e) {
    ElMessage.error(e.message || '上传失败')
  }
}

async function onMapPick({ latitude, longitude }) {
  form.latitude = latitude
  form.longitude = longitude
  try {
    const geo = await reverseGeocode(latitude, longitude)
    if (geo?.province) form.province = geo.province
    if (geo?.city) form.city = geo.city
    if (geo?.district) form.district = geo.district
    if (geo?.address && !form.address) form.address = geo.address
  } catch (_) {
    /* 逆地理失败不阻断选点 */
  }
}

async function load() {
  loading.value = true
  try {
    const res = await listMyShops()
    const list = pickList(res) || (Array.isArray(res) ? res : [])
    const sid = Number(auth.shopId || 0)
    const shop = list.find((s) => Number(s.id) === sid) || list[0]
    if (!shop) {
      ElMessage.warning('未找到店铺')
      return
    }
    shopId.value = shop.id
    Object.assign(form, {
      name: shop.name || '',
      logo: shop.logo || '',
      contact_name: shop.contact_name || '',
      contact_phone: shop.contact_phone || '',
      description: shop.description || '',
      category: shop.category || '',
      province: shop.province || '',
      city: shop.city || '',
      district: shop.district || '',
      address: shop.address || '',
      storefront_image: shop.storefront_image || '',
      latitude: Number(shop.latitude) || 0,
      longitude: Number(shop.longitude) || 0,
      local_enabled: Number(shop.local_enabled) || 0,
      images: Array.isArray(shop.images) ? [...shop.images] : [],
    })
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!shopId.value) return
  if (form.local_enabled === 1 && (!form.latitude || !form.longitude)) {
    ElMessage.warning('开启本地商家前请先在地图选点')
    return
  }
  saving.value = true
  try {
    await updateMyShop(shopId.value, { ...form })
    ElMessage.success('已保存')
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.form { max-width: 720px; margin-top: 16px; }
.cover { width: 160px; height: 100px; object-fit: cover; border-radius: 8px; }
.imgs { display: flex; flex-wrap: wrap; gap: 12px; align-items: flex-start; }
.img-item { display: flex; flex-direction: column; align-items: center; gap: 4px; }
.img-item img { width: 96px; height: 96px; object-fit: cover; border-radius: 8px; }
.coord { margin-top: 8px; font-size: 13px; color: #64748b; }
</style>
