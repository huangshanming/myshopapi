<template>
  <view class="empty">
    <text class="emoji">🛒</text>
    <text class="title">购物车</text>
    <text class="sub">即将开放，请先在商品详情页立即购买</text>
  </view>
</template>

<style scoped>
.empty {
  display: flex; flex-direction: column; align-items: center;
  padding-top: 200rpx; color: #71717a;
}
.emoji { font-size: 80rpx; margin-bottom: 24rpx; }
.title { font-size: 32rpx; color: #1d1d1f; font-weight: 600; }
.sub { margin-top: 12rpx; font-size: 24rpx; text-align: center; padding: 0 48rpx; }
</style>
