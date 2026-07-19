-- 全局省市区（行政区划）
-- mysql -u homestead -p --default-character-set=utf8mb4 mymall < scripts/alter-regions.sql
-- 数据由 user-service 启动时从 pca-code.json 自动导入（表为空时）

USE mymall;

CREATE TABLE IF NOT EXISTS regions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(12) NOT NULL COMMENT '行政区划代码',
    name VARCHAR(64) NOT NULL COMMENT '名称',
    parent_code VARCHAR(12) NOT NULL DEFAULT '' COMMENT '父级代码，省级为空',
    level TINYINT NOT NULL COMMENT '1省 2市 3区县',
    sort INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code (code),
    KEY idx_parent (parent_code),
    KEY idx_level (level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='全国省市区';

-- 用户地址补充区划代码（可选，便于回显选择）
SET @db := DATABASE();
SET @exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=@db AND TABLE_NAME='user_addresses' AND COLUMN_NAME='province_code');
SET @sql := IF(@exists=0, "ALTER TABLE user_addresses ADD COLUMN province_code VARCHAR(12) NOT NULL DEFAULT '' AFTER detail, ADD COLUMN city_code VARCHAR(12) NOT NULL DEFAULT '' AFTER province_code, ADD COLUMN district_code VARCHAR(12) NOT NULL DEFAULT '' AFTER city_code", 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
