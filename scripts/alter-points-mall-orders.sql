-- 积分商城兑换订单
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-points-mall-orders.sql
USE mymall;

CREATE TABLE IF NOT EXISTS points_exchange_orders (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL COMMENT '兑换单号',
    user_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    product_name VARCHAR(100) NOT NULL DEFAULT '',
    product_cover VARCHAR(512) NOT NULL DEFAULT '',
    quantity INT NOT NULL DEFAULT 1,
    points_cost INT NOT NULL DEFAULT 0 COMMENT '消耗积分',
    status VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT 'pending|shipped|completed|cancelled',
    receiver_name VARCHAR(64) NOT NULL DEFAULT '',
    receiver_phone VARCHAR(32) NOT NULL DEFAULT '',
    receiver_address VARCHAR(255) NOT NULL DEFAULT '',
    ship_company VARCHAR(64) NOT NULL DEFAULT '',
    ship_no VARCHAR(64) NOT NULL DEFAULT '',
    admin_remark VARCHAR(255) NOT NULL DEFAULT '',
    shipped_at DATETIME NULL,
    completed_at DATETIME NULL,
    cancelled_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_order_no (order_no),
    KEY idx_user_created (user_id, created_at),
    KEY idx_status_created (status, created_at),
    KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分商城兑换订单';
