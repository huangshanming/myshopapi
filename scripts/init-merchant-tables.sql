-- 多租户店铺表 + 用户角色 / 商品店铺 / 订单店铺
-- 执行：mysql -u homestead -p mymall < scripts/init-merchant-tables.sql
-- 可重复执行：列/索引已存在时会跳过（通过 information_schema 判断）

USE mymall;

SET @db := DATABASE();

-- users.role
SET @exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'users' AND COLUMN_NAME = 'role'
);
SET @sql := IF(@exists = 0,
  "ALTER TABLE users ADD COLUMN role VARCHAR(32) NOT NULL DEFAULT 'user' COMMENT '角色' AFTER status",
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- products.shop_id
SET @exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'products' AND COLUMN_NAME = 'shop_id'
);
SET @sql := IF(@exists = 0,
  "ALTER TABLE products ADD COLUMN shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属店铺ID' AFTER id, ADD INDEX idx_shop_id (shop_id)",
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- orders.shop_id
SET @exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'shop_id'
);
SET @sql := IF(@exists = 0,
  "ALTER TABLE orders ADD COLUMN shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属店铺ID' AFTER user_id, ADD INDEX idx_shop_id (shop_id)",
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS shops (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '店铺名称',
    logo VARCHAR(500) DEFAULT '' COMMENT 'Logo',
    contact_name VARCHAR(50) DEFAULT '' COMMENT '联系人',
    contact_phone CHAR(11) DEFAULT '' COMMENT '联系电话',
    description VARCHAR(500) DEFAULT '' COMMENT '简介',
    owner_user_id BIGINT UNSIGNED NOT NULL COMMENT '店主用户ID',
    status ENUM('pending','approved','rejected','disabled') NOT NULL DEFAULT 'pending' COMMENT '状态',
    reject_reason VARCHAR(255) DEFAULT '' COMMENT '拒绝原因',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_owner (owner_user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_applications (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '申请人',
    shop_name VARCHAR(100) NOT NULL,
    contact_name VARCHAR(50) NOT NULL,
    contact_phone CHAR(11) NOT NULL,
    description VARCHAR(500) DEFAULT '',
    status ENUM('pending','approved','rejected') NOT NULL DEFAULT 'pending',
    reject_reason VARCHAR(255) DEFAULT '',
    reviewed_by BIGINT UNSIGNED DEFAULT NULL,
    reviewed_at TIMESTAMP NULL DEFAULT NULL,
    shop_id BIGINT UNSIGNED DEFAULT NULL COMMENT '审核通过后关联店铺',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_members (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    shop_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    member_role ENUM('owner','staff') NOT NULL DEFAULT 'owner',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_shop_user (shop_id, user_id),
    INDEX idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
