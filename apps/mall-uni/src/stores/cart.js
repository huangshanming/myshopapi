const CART_KEY = 'mymall_cart'
const CHECKOUT_KEY = 'mymall_checkout'

function readCart() {
  try {
    const raw = uni.getStorageSync(CART_KEY)
    if (!raw) return []
    const list = typeof raw === 'string' ? JSON.parse(raw) : raw
    return Array.isArray(list) ? list : []
  } catch {
    return []
  }
}

function writeCart(list) {
  uni.setStorageSync(CART_KEY, JSON.stringify(list || []))
}

export function getCartItems() {
  return readCart()
}

export function getCartCount() {
  return readCart().reduce((n, it) => n + (Number(it.quantity) || 0), 0)
}

export function addToCart(item) {
  const list = readCart()
  const pid = Number(item.product_id)
  const sid = Number(item.sku_id || 0)
  const seckill = Number(item.seckill_entry_id || 0)
  const idx = list.findIndex(
    (x) =>
      Number(x.product_id) === pid &&
      Number(x.sku_id || 0) === sid &&
      Number(x.seckill_entry_id || 0) === seckill,
  )
  if (idx >= 0) {
    list[idx].quantity = (Number(list[idx].quantity) || 0) + (Number(item.quantity) || 1)
  } else {
    list.push({
      product_id: pid,
      sku_id: sid,
      name: item.name || '',
      image: item.image || '',
      price: Number(item.price) || 0,
      quantity: Number(item.quantity) || 1,
      shop_id: Number(item.shop_id) || 0,
      seckill_entry_id: seckill || undefined,
      selected: true,
    })
  }
  writeCart(list)
  return list
}

export function updateCartQty(productId, skuId, seckillEntryId, quantity) {
  const list = readCart()
  const idx = list.findIndex(
    (x) =>
      Number(x.product_id) === Number(productId) &&
      Number(x.sku_id || 0) === Number(skuId || 0) &&
      Number(x.seckill_entry_id || 0) === Number(seckillEntryId || 0),
  )
  if (idx < 0) return list
  if (quantity <= 0) list.splice(idx, 1)
  else list[idx].quantity = quantity
  writeCart(list)
  return list
}

export function removeCartItem(productId, skuId, seckillEntryId) {
  return updateCartQty(productId, skuId, seckillEntryId, 0)
}

export function setCartSelected(productId, skuId, seckillEntryId, selected) {
  const list = readCart()
  const item = list.find(
    (x) =>
      Number(x.product_id) === Number(productId) &&
      Number(x.sku_id || 0) === Number(skuId || 0) &&
      Number(x.seckill_entry_id || 0) === Number(seckillEntryId || 0),
  )
  if (item) item.selected = !!selected
  writeCart(list)
  return list
}

export function setAllSelected(selected) {
  const list = readCart().map((x) => ({ ...x, selected: !!selected }))
  writeCart(list)
  return list
}

export function clearCartItems(keys) {
  if (!keys?.length) return readCart()
  const set = new Set(keys.map((k) => `${k.product_id}_${k.sku_id || 0}_${k.seckill_entry_id || 0}`))
  const list = readCart().filter(
    (x) => !set.has(`${x.product_id}_${x.sku_id || 0}_${x.seckill_entry_id || 0}`),
  )
  writeCart(list)
  return list
}

export function setCheckoutPayload(payload) {
  uni.setStorageSync(CHECKOUT_KEY, JSON.stringify(payload || {}))
}

export function getCheckoutPayload() {
  try {
    const raw = uni.getStorageSync(CHECKOUT_KEY)
    if (!raw) return null
    return typeof raw === 'string' ? JSON.parse(raw) : raw
  } catch {
    return null
  }
}

export function clearCheckoutPayload() {
  uni.removeStorageSync(CHECKOUT_KEY)
}
