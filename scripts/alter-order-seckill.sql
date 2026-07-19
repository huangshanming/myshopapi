-- 订单项关联秒杀报名（可空，0 表示非秒杀）
-- mysql -u homestead -p --default-character-set=utf8mb4 mymall < scripts/alter-order-seckill.sql

USE mymall;

SET @exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'order_items' AND COLUMN_NAME = 'seckill_entry_id'
);
SET @sql := IF(@exists = 0,
  "ALTER TABLE order_items ADD COLUMN seckill_entry_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '秒杀报名ID，0表示非秒杀' AFTER quantity",
  'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
