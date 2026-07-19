-- 首页付费展位
USE mymall;

CREATE TABLE IF NOT EXISTS homepage_slot_packages (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    slot_type VARCHAR(32) NOT NULL COMMENT 'brand_shop/quality_shop/article',
    name VARCHAR(100) NOT NULL,
    price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    duration_days INT NOT NULL DEFAULT 1,
    status VARCHAR(16) NOT NULL DEFAULT 'on' COMMENT 'on/off',
    sort INT NOT NULL DEFAULT 0,
    remark VARCHAR(255) NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_type_status (slot_type, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='首页展位套餐';

CREATE TABLE IF NOT EXISTS homepage_slot_settings (
    slot_type VARCHAR(32) NOT NULL PRIMARY KEY,
    home_limit INT NOT NULL DEFAULT 8 COMMENT '首页条带最多展示条数',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='首页展位展示设置';

INSERT INTO homepage_slot_settings (slot_type, home_limit) VALUES
('brand_shop', 8),
('quality_shop', 6),
('article', 6)
ON DUPLICATE KEY UPDATE home_limit=VALUES(home_limit);

CREATE TABLE IF NOT EXISTS homepage_slot_orders (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    shop_id BIGINT UNSIGNED NOT NULL,
    slot_type VARCHAR(32) NOT NULL,
    package_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    target_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '店铺=shop_id；文章=article_id',
    amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    duration_days INT NOT NULL DEFAULT 0,
    start_at DATETIME NOT NULL,
    end_at DATETIME NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'active' COMMENT 'active/expired/cancelled',
    pay_source VARCHAR(16) NOT NULL DEFAULT 'wallet' COMMENT 'wallet/admin',
    wallet_log_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    operator_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_type_status_end (slot_type, status, end_at),
    KEY idx_target_type (target_id, slot_type),
    KEY idx_shop (shop_id, slot_type, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='首页展位订单';

-- 示例套餐
INSERT INTO homepage_slot_packages (slot_type, name, price, duration_days, status, sort, remark) VALUES
('brand_shop', '品牌店·10天', 500.00, 10, 'on', 10, ''),
('brand_shop', '品牌店·30天', 1200.00, 30, 'on', 20, ''),
('quality_shop', '优质商户·10天', 300.00, 10, 'on', 10, ''),
('quality_shop', '优质商户·30天', 800.00, 30, 'on', 20, ''),
('article', '种草置顶·7天', 200.00, 7, 'on', 10, ''),
('article', '种草置顶·15天', 400.00, 15, 'on', 20, '');
