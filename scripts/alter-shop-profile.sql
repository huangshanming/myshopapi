-- 店铺入驻档资料字段
-- 执行：mysql -u homestead -p mymall < scripts/alter-shop-profile.sql
USE mymall;

SET @db := DATABASE();

-- shops 扩展字段
SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shops' AND COLUMN_NAME='category');
SET @sql := IF(@exists=0, "ALTER TABLE shops ADD COLUMN category VARCHAR(50) DEFAULT '' COMMENT '经营类目' AFTER description", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shops' AND COLUMN_NAME='province');
SET @sql := IF(@exists=0, "ALTER TABLE shops ADD COLUMN province VARCHAR(50) DEFAULT '' AFTER category", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shops' AND COLUMN_NAME='city');
SET @sql := IF(@exists=0, "ALTER TABLE shops ADD COLUMN city VARCHAR(50) DEFAULT '' AFTER province", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shops' AND COLUMN_NAME='district');
SET @sql := IF(@exists=0, "ALTER TABLE shops ADD COLUMN district VARCHAR(50) DEFAULT '' AFTER city", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shops' AND COLUMN_NAME='address');
SET @sql := IF(@exists=0, "ALTER TABLE shops ADD COLUMN address VARCHAR(255) DEFAULT '' COMMENT '详细地址' AFTER district", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shops' AND COLUMN_NAME='business_license_no');
SET @sql := IF(@exists=0, "ALTER TABLE shops ADD COLUMN business_license_no VARCHAR(64) DEFAULT '' COMMENT '营业执照号' AFTER address", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shops' AND COLUMN_NAME='legal_person');
SET @sql := IF(@exists=0, "ALTER TABLE shops ADD COLUMN legal_person VARCHAR(50) DEFAULT '' COMMENT '法人' AFTER business_license_no", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shops' AND COLUMN_NAME='license_image');
SET @sql := IF(@exists=0, "ALTER TABLE shops ADD COLUMN license_image VARCHAR(500) DEFAULT '' COMMENT '执照图' AFTER legal_person", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shops' AND COLUMN_NAME='storefront_image');
SET @sql := IF(@exists=0, "ALTER TABLE shops ADD COLUMN storefront_image VARCHAR(500) DEFAULT '' COMMENT '门头图' AFTER license_image", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- shop_applications 同步字段
SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shop_applications' AND COLUMN_NAME='category');
SET @sql := IF(@exists=0, "ALTER TABLE shop_applications ADD COLUMN category VARCHAR(50) DEFAULT '' AFTER description", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shop_applications' AND COLUMN_NAME='province');
SET @sql := IF(@exists=0, "ALTER TABLE shop_applications ADD COLUMN province VARCHAR(50) DEFAULT '' AFTER category", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shop_applications' AND COLUMN_NAME='city');
SET @sql := IF(@exists=0, "ALTER TABLE shop_applications ADD COLUMN city VARCHAR(50) DEFAULT '' AFTER province", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shop_applications' AND COLUMN_NAME='district');
SET @sql := IF(@exists=0, "ALTER TABLE shop_applications ADD COLUMN district VARCHAR(50) DEFAULT '' AFTER city", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shop_applications' AND COLUMN_NAME='address');
SET @sql := IF(@exists=0, "ALTER TABLE shop_applications ADD COLUMN address VARCHAR(255) DEFAULT '' AFTER district", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shop_applications' AND COLUMN_NAME='business_license_no');
SET @sql := IF(@exists=0, "ALTER TABLE shop_applications ADD COLUMN business_license_no VARCHAR(64) DEFAULT '' AFTER address", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shop_applications' AND COLUMN_NAME='legal_person');
SET @sql := IF(@exists=0, "ALTER TABLE shop_applications ADD COLUMN legal_person VARCHAR(50) DEFAULT '' AFTER business_license_no", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shop_applications' AND COLUMN_NAME='license_image');
SET @sql := IF(@exists=0, "ALTER TABLE shop_applications ADD COLUMN license_image VARCHAR(500) DEFAULT '' AFTER legal_person", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='shop_applications' AND COLUMN_NAME='storefront_image');
SET @sql := IF(@exists=0, "ALTER TABLE shop_applications ADD COLUMN storefront_image VARCHAR(500) DEFAULT '' AFTER license_image", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
