-- 物流公司主数据
USE mymall;

CREATE TABLE IF NOT EXISTS logistics_companies (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL COMMENT '公司名称',
    code VARCHAR(32) NOT NULL COMMENT '编码',
    sort INT NOT NULL DEFAULT 0 COMMENT '排序，越小越靠前',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1启用 0停用',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code (code),
    INDEX idx_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT IGNORE INTO logistics_companies (name, code, sort, status) VALUES
('顺丰速运', 'SF', 10, 1),
('中通快递', 'ZTO', 20, 1),
('圆通速递', 'YTO', 30, 1),
('韵达快递', 'YD', 40, 1),
('申通快递', 'STO', 50, 1),
('EMS', 'EMS', 60, 1),
('京东物流', 'JD', 70, 1),
('德邦快递', 'DBL', 80, 1);
