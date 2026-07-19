-- 商家钱包 + 限时秒杀
-- 执行：mysql -u homestead -p mymall < scripts/alter-seckill-wallet.sql

USE mymall;

CREATE TABLE IF NOT EXISTS shop_wallets (
    shop_id BIGINT UNSIGNED NOT NULL PRIMARY KEY COMMENT '店铺ID',
    balance DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '可用余额',
    frozen_balance DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '冻结余额',
    deposit DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '保证金',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='店铺钱包';

CREATE TABLE IF NOT EXISTS shop_wallet_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    shop_id BIGINT UNSIGNED NOT NULL,
    change_type VARCHAR(32) NOT NULL COMMENT 'admin_adjust/seckill_apply',
    amount DECIMAL(12,2) NOT NULL COMMENT '变动金额，正加负减',
    balance_after DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    frozen_after DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    deposit_after DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    remark VARCHAR(255) NOT NULL DEFAULT '',
    operator_user_id BIGINT UNSIGNED NULL,
    ref_type VARCHAR(32) NOT NULL DEFAULT '',
    ref_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_shop_created (shop_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='钱包流水';

CREATE TABLE IF NOT EXISTS seckill_rules (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    duration_hours INT NOT NULL DEFAULT 24 COMMENT '场次时长(小时)',
    apply_fee DECIMAL(12,2) NOT NULL DEFAULT 10.00 COMMENT '报名费',
    max_entries_per_shop INT NOT NULL DEFAULT 5 COMMENT '每店每场最多报名数',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1启用 0停用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀规则';

CREATE TABLE IF NOT EXISTS seckill_sessions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    rule_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    start_at DATETIME NOT NULL,
    end_at DATETIME NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'active' COMMENT 'active/ended',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_status_end (status, end_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀场次';

CREATE TABLE IF NOT EXISTS seckill_entries (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    session_id BIGINT UNSIGNED NOT NULL,
    shop_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    product_name VARCHAR(200) NOT NULL DEFAULT '',
    product_image VARCHAR(500) NOT NULL DEFAULT '',
    origin_price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    seckill_price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    seckill_stock INT NOT NULL DEFAULT 0,
    fee_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_session (session_id),
    KEY idx_shop_session (shop_id, session_id),
    UNIQUE KEY uk_session_shop_product (session_id, shop_id, product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀报名';

INSERT INTO seckill_rules (id, duration_hours, apply_fee, max_entries_per_shop, status)
SELECT 1, 24, 10.00, 5, 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM seckill_rules LIMIT 1);
