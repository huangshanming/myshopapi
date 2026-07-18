-- 订单履约字段 + 售后表
USE mymall;

ALTER TABLE orders
  MODIFY COLUMN status ENUM(
    'pending','confirmed','failed','cancelled','shipped','completed'
  ) NOT NULL DEFAULT 'pending';

-- 下列 ADD 若列已存在会报错，可忽略后继续
ALTER TABLE orders ADD COLUMN receiver_name VARCHAR(64) NOT NULL DEFAULT '' AFTER total_amount;
ALTER TABLE orders ADD COLUMN receiver_phone VARCHAR(20) NOT NULL DEFAULT '' AFTER receiver_name;
ALTER TABLE orders ADD COLUMN receiver_address VARCHAR(255) NOT NULL DEFAULT '' AFTER receiver_phone;
ALTER TABLE orders ADD COLUMN ship_company VARCHAR(64) NOT NULL DEFAULT '' AFTER receiver_address;
ALTER TABLE orders ADD COLUMN ship_no VARCHAR(64) NOT NULL DEFAULT '' AFTER ship_company;
ALTER TABLE orders ADD COLUMN shipped_at TIMESTAMP NULL DEFAULT NULL AFTER ship_no;
ALTER TABLE orders ADD COLUMN completed_at TIMESTAMP NULL DEFAULT NULL AFTER shipped_at;
ALTER TABLE orders ADD COLUMN remark VARCHAR(255) NOT NULL DEFAULT '' AFTER completed_at;

CREATE TABLE IF NOT EXISTS order_after_sales (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    order_no VARCHAR(64) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    type ENUM('refund','return_refund') NOT NULL DEFAULT 'refund',
    reason VARCHAR(500) NOT NULL DEFAULT '',
    amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    status ENUM('pending','approved','rejected','refunded','closed') NOT NULL DEFAULT 'pending',
    admin_remark VARCHAR(500) NOT NULL DEFAULT '',
    handled_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order_id (order_id),
    INDEX idx_shop_id (shop_id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    CONSTRAINT fk_after_sales_order FOREIGN KEY (order_id) REFERENCES orders(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
