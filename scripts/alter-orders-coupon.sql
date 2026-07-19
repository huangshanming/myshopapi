-- 订单优惠券金额字段
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-orders-coupon.sql
USE mymall;

SET @col := (SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='orders' AND COLUMN_NAME='goods_amount');
SET @sql := IF(@col=0,
  'ALTER TABLE orders ADD COLUMN goods_amount DECIMAL(12,2) NOT NULL DEFAULT 0 AFTER total_amount, ADD COLUMN discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0 AFTER goods_amount, ADD COLUMN pay_amount DECIMAL(12,2) NOT NULL DEFAULT 0 AFTER discount_amount, ADD COLUMN user_coupon_id BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER pay_amount',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 历史单对齐
UPDATE orders
SET goods_amount = total_amount,
    pay_amount = total_amount,
    discount_amount = 0
WHERE pay_amount = 0 AND goods_amount = 0 AND total_amount > 0;
