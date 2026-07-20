-- 积分商城商品表
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-points-mall.sql
USE mymall;

CREATE TABLE IF NOT EXISTS points_products (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '商品名称',
    cover_url VARCHAR(512) NOT NULL DEFAULT '' COMMENT '封面图',
    description VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '说明',
    points_price INT NOT NULL DEFAULT 0 COMMENT '兑换所需积分',
    stock INT NOT NULL DEFAULT 0 COMMENT '库存',
    per_user_limit INT NOT NULL DEFAULT 0 COMMENT '每人限兑次数,0=不限',
    status VARCHAR(16) NOT NULL DEFAULT 'off' COMMENT 'on|off',
    sort INT NOT NULL DEFAULT 0 COMMENT '排序,越大越靠前',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_status_sort (status, sort),
    KEY idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分商城商品';
