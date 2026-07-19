-- 用户收货地址
-- mysql -u homestead -p --default-character-set=utf8mb4 mymall < scripts/alter-user-addresses.sql

USE mymall;

CREATE TABLE IF NOT EXISTS user_addresses (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    receiver_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '收货人',
    receiver_phone VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号',
    province VARCHAR(64) NOT NULL DEFAULT '',
    city VARCHAR(64) NOT NULL DEFAULT '',
    district VARCHAR(64) NOT NULL DEFAULT '',
    detail VARCHAR(255) NOT NULL DEFAULT '' COMMENT '详细地址',
    is_default TINYINT NOT NULL DEFAULT 0 COMMENT '1默认',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_user (user_id),
    KEY idx_user_default (user_id, is_default)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收货地址';
