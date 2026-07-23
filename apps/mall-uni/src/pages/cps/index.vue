<template>
  <view class="page">
    <view class="card" @tap="goRedpacket">
      <view class="text">
        <text class="title">平台优惠红包</text>
        <text class="sub">外卖打车等联盟活动 · 复制链接去领取</text>
      </view>
      <text class="arrow">›</text>
    </view>
    <view class="card goods" @tap="goProducts">
      <view class="text">
        <text class="title">优惠商品</text>
        <text class="sub">淘宝 / 京东 / 拼多多 / 唯品会 · 搜索好物复制推广链接</text>
      </view>
      <text class="arrow">›</text>
    </view>
  </view>
</template>

<script setup>
function goRedpacket() {
  uni.navigateTo({ url: '/pages/cps/redpacket' })
}

function goProducts() {
  uni.navigateTo({ url: '/pages/cps/products' })
}
</script>

<style scoped>
.page { padding: 32rpx; min-height: 100vh; background: #f7f3ec; }
.card {
  display: flex; align-items: center; justify-content: space-between;
  background: linear-gradient(135deg, #efe0c8 0%, #d4b890 100%);
  border-radius: 28rpx; padding: 40rpx 32rpx; margin-bottom: 24rpx;
  box-shadow: 0 8rpx 28rpx rgba(200,168,118,.18);
}
.card.goods {
  background: linear-gradient(135deg, #f3ebe0 0%, #e8d5b5 55%, #c8a876 100%);
}
.text { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 10rpx; }
.title { font-size: 34rpx; font-weight: 700; color: #5c4a32; }
.sub { font-size: 24rpx; color: rgba(92,74,50,.82); line-height: 1.4; }
.arrow { font-size: 44rpx; color: #8b7355; font-weight: 300; padding-left: 16rpx; }
</style>
