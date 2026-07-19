-- 商品收藏
USE mymall;

CREATE TABLE IF NOT EXISTS product_favorites (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_product (user_id, product_id),
    KEY idx_user_created (user_id, created_at),
    KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户商品收藏';

SET @db := DATABASE();
SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='products' AND COLUMN_NAME='avg_rating');
SET @sql := IF(@exists=0, "ALTER TABLE products ADD COLUMN avg_rating DECIMAL(3,2) NOT NULL DEFAULT 0.00 COMMENT '平均评分' AFTER collect_count, ADD COLUMN review_count INT NOT NULL DEFAULT 0 COMMENT '评价数' AFTER avg_rating, ADD COLUMN good_rate DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT '好评率%' AFTER review_count", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
