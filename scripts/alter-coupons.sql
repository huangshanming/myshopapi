-- 优惠券：模板 / 范围 / 用户券 / 定向发放 / 核销流水
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-coupons.sql
USE mymall;

CREATE TABLE IF NOT EXISTS coupons (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL DEFAULT '',
  issuer_type ENUM('platform','shop') NOT NULL DEFAULT 'platform',
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  coupon_type ENUM('full_reduce','no_threshold','category','product','discount') NOT NULL DEFAULT 'full_reduce',
  threshold_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  discount_rate DECIMAL(6,4) NOT NULL DEFAULT 0 COMMENT '0.80=八折',
  max_discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  scope_type ENUM('all','category','product') NOT NULL DEFAULT 'all',
  total_count INT NOT NULL DEFAULT 0 COMMENT '0=不限',
  claimed_count INT NOT NULL DEFAULT 0,
  per_user_limit INT NOT NULL DEFAULT 1,
  valid_type ENUM('fixed','relative') NOT NULL DEFAULT 'fixed',
  valid_start DATETIME NULL,
  valid_end DATETIME NULL,
  valid_days INT NOT NULL DEFAULT 0,
  stackable TINYINT NOT NULL DEFAULT 0,
  user_identity ENUM('all','new','old') NOT NULL DEFAULT 'all',
  channels JSON NULL,
  status ENUM('draft','on','off','expired') NOT NULL DEFAULT 'draft',
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_issuer (issuer_type, shop_id, status),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coupon_scopes (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  coupon_id BIGINT UNSIGNED NOT NULL,
  ref_type ENUM('category','product') NOT NULL,
  ref_id BIGINT UNSIGNED NOT NULL,
  KEY idx_coupon (coupon_id),
  UNIQUE KEY uk_coupon_ref (coupon_id, ref_type, ref_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_coupons (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  coupon_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  status ENUM('unused','locked','used','expired') NOT NULL DEFAULT 'unused',
  source ENUM('direct','order_gift','popup','targeted') NOT NULL DEFAULT 'direct',
  valid_start DATETIME NOT NULL,
  valid_end DATETIME NOT NULL,
  order_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  locked_at DATETIME NULL,
  used_at DATETIME NULL,
  claim_batch_no VARCHAR(64) NOT NULL DEFAULT '',
  discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0 COMMENT '核销时实际抵扣',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_user_status (user_id, status),
  KEY idx_coupon_user (coupon_id, user_id),
  KEY idx_order (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coupon_grants (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  coupon_id BIGINT UNSIGNED NOT NULL,
  operator_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  issuer_type ENUM('platform','shop') NOT NULL DEFAULT 'platform',
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  user_count INT NOT NULL DEFAULT 0,
  success_count INT NOT NULL DEFAULT 0,
  batch_no VARCHAR(64) NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_coupon (coupon_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coupon_redeem_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_coupon_id BIGINT UNSIGNED NOT NULL,
  coupon_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  order_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  action ENUM('redeem','unlock','return') NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_coupon (coupon_id),
  KEY idx_order (order_id),
  KEY idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
