<script>
export default {
  onLaunch() {},
  onShow() {},
  onHide() {},
}
</script>

<style lang="scss">
page {
  background: #fafafa;
  color: #1d1d1f;
  font-size: 28rpx;
  font-family: system-ui, -apple-system, "PingFang SC", sans-serif;
}

:root,
page {
  --shop-gold: #c8a876;
  --shop-gold-light: #e6d5bc;
  --shop-red: #d83636;
  --shop-text: #1d1d1f;
  --shop-sub: #71717a;
  --shop-bg: #fafafa;
  --shop-card: #ffffff;
  --shop-line: #f0f0f0;
}

.line-1 {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.line-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
/* H5 横向滚动条隐藏 */
.scroll-hide,
.scroll-hide ::-webkit-scrollbar,
uni-scroll-view::-webkit-scrollbar,
.uni-scroll-view::-webkit-scrollbar,
.uni-scroll-view-content::-webkit-scrollbar {
  display: none !important;
  width: 0 !important;
  height: 0 !important;
  background: transparent !important;
}
.scroll-hide,
uni-scroll-view,
.uni-scroll-view {
  scrollbar-width: none; /* firefox */
  -ms-overflow-style: none;
}
.gold-gradient {
  background: linear-gradient(135deg, #bfa472 0%, #d4b890 100%);
}
.sale-gradient {
  background: linear-gradient(135deg, #d83636, #f25757);
}
.section-card {
  margin: 24rpx 32rpx 0;
  background: var(--shop-card);
  border-radius: 32rpx;
  box-shadow: 0 4rpx 24rpx rgba(200, 168, 118, 0.08);
}
.sec-title {
  font-size: 34rpx;
  font-weight: 600;
  color: var(--shop-text);
}
.sec-more {
  font-size: 22rpx;
  color: var(--shop-sub);
}
.price {
  color: var(--shop-red);
  font-weight: 700;
}
.price-old {
  color: var(--shop-sub);
  font-size: 20rpx;
  text-decoration: line-through;
  margin-left: 8rpx;
}
</style>
