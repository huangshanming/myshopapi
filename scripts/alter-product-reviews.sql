-- 订单评价 + 订单状态 reviewed
USE mymall;

ALTER TABLE orders
  MODIFY COLUMN status ENUM('pending','confirmed','failed','cancelled','shipped','completed','reviewed')
  NOT NULL DEFAULT 'pending';

SET @db := DATABASE();
SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='orders' AND COLUMN_NAME='reviewed_at');
SET @sql := IF(@exists=0, "ALTER TABLE orders ADD COLUMN reviewed_at DATETIME NULL AFTER completed_at", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS product_reviews (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    order_no VARCHAR(64) NOT NULL DEFAULT '',
    user_id BIGINT UNSIGNED NOT NULL,
    shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    product_id BIGINT UNSIGNED NOT NULL,
    order_item_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    sku_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    sku_snapshot VARCHAR(512) NOT NULL DEFAULT '',
    rating TINYINT NOT NULL,
    content VARCHAR(1000) NOT NULL DEFAULT '',
    is_anonymous TINYINT NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'visible' COMMENT 'visible/deleted',
    merchant_reply VARCHAR(500) NOT NULL DEFAULT '',
    replied_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_order (order_id),
    KEY idx_product_created (product_id, created_at),
    KEY idx_shop_rating (shop_id, rating, status),
    KEY idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品订单评价';

CREATE TABLE IF NOT EXISTS product_review_images (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    review_id BIGINT UNSIGNED NOT NULL,
    url VARCHAR(500) NOT NULL,
    sort INT NOT NULL DEFAULT 0,
    KEY idx_review (review_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评价图片';
