-- 商家端商品中台表结构
-- 执行：mysql -u homestead -psecret mymall < scripts/alter-product-center.sql
USE mymall;
SET @db := DATABASE();

-- 迁移旧状态
UPDATE products SET status = 'draft' WHERE status IN ('pending','approved','rejected');
UPDATE products SET status = 'deleted' WHERE deleted_at IS NOT NULL AND status <> 'deleted';

-- products 增量字段（幂等）
SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='products' AND COLUMN_NAME='product_type');
SET @sql := IF(@exists=0, "ALTER TABLE products ADD COLUMN product_type ENUM('physical','fresh','virtual') NOT NULL DEFAULT 'physical' COMMENT '商品类型' AFTER category_id", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='products' AND COLUMN_NAME='spec_json');
SET @sql := IF(@exists=0, "ALTER TABLE products ADD COLUMN spec_json JSON NULL COMMENT '规格定义' AFTER product_type", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='products' AND COLUMN_NAME='schedule_on_at');
SET @sql := IF(@exists=0, "ALTER TABLE products ADD COLUMN schedule_on_at DATETIME NULL COMMENT '定时上架' AFTER publish_time", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='products' AND COLUMN_NAME='schedule_off_at');
SET @sql := IF(@exists=0, "ALTER TABLE products ADD COLUMN schedule_off_at DATETIME NULL COMMENT '定时下架' AFTER schedule_on_at", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='products' AND COLUMN_NAME='copy_from_id');
SET @sql := IF(@exists=0, "ALTER TABLE products ADD COLUMN copy_from_id BIGINT UNSIGNED NULL COMMENT '复制来源' AFTER schedule_off_at", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 收窄 status 枚举（MySQL 允许 MODIFY）
ALTER TABLE products MODIFY COLUMN status ENUM('draft','on_sale','off_sale','deleted','pending','approved','rejected') NOT NULL DEFAULT 'draft' COMMENT '商品状态';

CREATE TABLE IF NOT EXISTS product_skus (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  shop_id BIGINT UNSIGNED NOT NULL COMMENT '店铺ID',
  sku_no VARCHAR(64) NOT NULL COMMENT 'SKU编码',
  spec_values JSON NOT NULL COMMENT '规格值',
  spec_key VARCHAR(255) NOT NULL COMMENT '规格唯一键',
  sale_price DECIMAL(10,2) NOT NULL COMMENT '售价',
  market_price DECIMAL(10,2) NULL,
  cost_price DECIMAL(10,2) NULL,
  stock INT NOT NULL DEFAULT 0,
  stock_warn INT NOT NULL DEFAULT 10,
  barcode VARCHAR(64) NULL,
  status ENUM('enabled','disabled') NOT NULL DEFAULT 'enabled',
  sold_count INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL DEFAULT NULL,
  UNIQUE KEY uk_sku_no (sku_no),
  UNIQUE KEY uk_product_spec (product_id, spec_key),
  KEY idx_shop_stock (shop_id, stock),
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_images (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  product_id BIGINT UNSIGNED NOT NULL,
  shop_id BIGINT UNSIGNED NOT NULL,
  url VARCHAR(500) NOT NULL,
  typ ENUM('main','gallery','detail') NOT NULL DEFAULT 'gallery',
  sort INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_tags (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  name VARCHAR(50) NOT NULL,
  color VARCHAR(20) DEFAULT '',
  status TINYINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_shop_name (shop_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_tag_rels (
  product_id BIGINT UNSIGNED NOT NULL,
  tag_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (product_id, tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_attr_templates (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  name VARCHAR(100) NOT NULL,
  attrs_json JSON NOT NULL,
  status TINYINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_attrs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  product_id BIGINT UNSIGNED NOT NULL,
  template_id BIGINT UNSIGNED NULL,
  attr_key VARCHAR(64) NOT NULL,
  attr_label VARCHAR(100) NOT NULL,
  attr_value VARCHAR(500) NOT NULL,
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_schedules (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  product_id BIGINT UNSIGNED NOT NULL,
  shop_id BIGINT UNSIGNED NOT NULL,
  action ENUM('on_sale','off_sale') NOT NULL,
  run_at DATETIME NOT NULL,
  status ENUM('pending','done','cancelled') NOT NULL DEFAULT 'pending',
  locked_at DATETIME NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_due (status, run_at),
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_batch_jobs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  job_type VARCHAR(32) NOT NULL,
  payload_json JSON NOT NULL,
  progress INT NOT NULL DEFAULT 0,
  total INT NOT NULL DEFAULT 0,
  status ENUM('pending','running','success','failed','partial') NOT NULL DEFAULT 'pending',
  result_msg VARCHAR(1000) DEFAULT '',
  operator_id BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_shop (shop_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_op_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  product_id BIGINT UNSIGNED NULL,
  operator_id BIGINT UNSIGNED NOT NULL,
  action VARCHAR(64) NOT NULL,
  before_json JSON NULL,
  after_json JSON NULL,
  ip VARCHAR(64) DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_shop_time (shop_id, created_at),
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 商家端 RBAC
CREATE TABLE IF NOT EXISTS shop_roles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(100) NOT NULL,
  status TINYINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_shop_code (shop_id, code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_menus (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  name VARCHAR(100) NOT NULL,
  type VARCHAR(20) NOT NULL DEFAULT 'menu',
  path VARCHAR(200) DEFAULT '',
  component VARCHAR(200) DEFAULT '',
  icon VARCHAR(50) DEFAULT '',
  perms VARCHAR(100) DEFAULT '',
  sort INT NOT NULL DEFAULT 0,
  visible TINYINT NOT NULL DEFAULT 1,
  status TINYINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_role_menus (
  role_id BIGINT UNSIGNED NOT NULL,
  menu_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (role_id, menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_user_roles (
  shop_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  role_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (shop_id, user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 订单行 SKU 快照
SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='order_items' AND COLUMN_NAME='sku_id');
SET @sql := IF(@exists=0, "ALTER TABLE order_items ADD COLUMN sku_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'SKU ID' AFTER product_id", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='order_items' AND COLUMN_NAME='sku_snapshot');
SET @sql := IF(@exists=0, "ALTER TABLE order_items ADD COLUMN sku_snapshot JSON NULL COMMENT 'SKU规格快照' AFTER product_name", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 种子：商家菜单（按模块分层：目录 + 子菜单 + 按钮）
INSERT INTO shop_menus (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
-- 一级目录
(100, 0, '商品中心', 'dir', '', '', 'Goods', '', 10, 1, 1),
(101, 0, '库存管理', 'dir', '', '', 'Box', '', 20, 1, 1),
(102, 0, '订单中心', 'dir', '', '', 'List', '', 30, 1, 1),
(103, 0, '店铺设置', 'dir', '', '', 'Setting', '', 90, 1, 1),
-- 商品中心
(1, 100, '商品列表', 'menu', '/merchant/products', 'merchant/Products', 'Goods', 'product:list', 10, 1, 1),
(2, 100, '发布商品', 'menu', '/merchant/products/edit', 'merchant/ProductEdit', 'Edit', 'product:edit', 11, 1, 1),
(3, 100, '回收站', 'menu', '/merchant/products/recycle', 'merchant/ProductRecycle', 'Delete', 'product:recycle', 12, 1, 1),
(7, 100, '操作日志', 'menu', '/merchant/products/op-logs', 'merchant/OpLogs', 'Document', 'product:list', 13, 1, 1),
-- 库存 / 订单 / 设置
(4, 101, '库存预警', 'menu', '/merchant/stocks/warnings', 'merchant/StockWarnings', 'Warning', 'stock:warn', 10, 1, 1),
(5, 102, '店铺订单', 'menu', '/merchant/orders', 'merchant/Orders', 'List', 'order:list', 10, 1, 1),
(6, 103, '员工权限', 'menu', '/merchant/staff', 'merchant/Staff', 'User', 'shop:staff', 10, 1, 1),
-- 按钮挂在商品列表下
(11, 1, '商品新增', 'button', '', '', '', 'product:add', 1, 1, 1),
(12, 1, '商品编辑', 'button', '', '', '', 'product:edit', 2, 1, 1),
(13, 1, '商品上下架', 'button', '', '', '', 'product:status', 3, 1, 1),
(14, 1, '批量操作', 'button', '', '', '', 'product:batch', 4, 1, 1),
(15, 1, '导入导出', 'button', '', '', '', 'product:import', 5, 1, 1)
ON DUPLICATE KEY UPDATE
  parent_id=VALUES(parent_id), name=VALUES(name), type=VALUES(type),
  path=VALUES(path), perms=VALUES(perms), sort=VALUES(sort), visible=VALUES(visible), status=VALUES(status);
