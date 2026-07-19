-- C 端用户钱包
-- mysql -u homestead -p --default-character-set=utf8mb4 mymall < scripts/alter-user-wallet.sql

USE mymall;

CREATE TABLE IF NOT EXISTS user_wallets (
    user_id BIGINT UNSIGNED NOT NULL PRIMARY KEY COMMENT '用户ID',
    balance DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '可用余额',
    frozen_balance DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '冻结余额',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户钱包';

CREATE TABLE IF NOT EXISTS user_wallet_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    change_type VARCHAR(32) NOT NULL COMMENT 'admin_adjust/order_freeze/order_unfreeze/order_settle',
    amount DECIMAL(12,2) NOT NULL COMMENT '变动金额',
    balance_after DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    frozen_after DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    remark VARCHAR(255) NOT NULL DEFAULT '',
    operator_user_id BIGINT UNSIGNED NULL,
    ref_type VARCHAR(32) NOT NULL DEFAULT '',
    ref_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_user_created (user_id, created_at),
    KEY idx_ref (ref_type, ref_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户钱包流水';
