-- 主题好物集市：固定坑位 + 套餐 + 订单（幂等）
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-homepage-theme-slots.sql
USE mymall;

CREATE TABLE IF NOT EXISTS homepage_theme_slots (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  slot_key VARCHAR(32) NOT NULL,
  position TINYINT UNSIGNED NOT NULL DEFAULT 1,
  name VARCHAR(100) NOT NULL DEFAULT '',
  `desc` VARCHAR(255) NOT NULL DEFAULT '',
  cover_url VARCHAR(500) NOT NULL DEFAULT '',
  default_link_type ENUM('none','category','shop','product') NOT NULL DEFAULT 'none',
  default_link_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  status ENUM('on','off') NOT NULL DEFAULT 'on',
  sort INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_slot_key (slot_key),
  UNIQUE KEY uk_position (position)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS homepage_theme_packages (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  theme_slot_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=通用任意坑位',
  name VARCHAR(100) NOT NULL,
  price DECIMAL(12,2) NOT NULL DEFAULT 0,
  duration_days INT NOT NULL DEFAULT 7,
  status ENUM('on','off') NOT NULL DEFAULT 'on',
  sort INT NOT NULL DEFAULT 0,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_slot_status (theme_slot_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS homepage_theme_orders (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  theme_slot_id BIGINT UNSIGNED NOT NULL,
  package_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  title VARCHAR(100) NOT NULL DEFAULT '',
  subtitle VARCHAR(255) NOT NULL DEFAULT '',
  cover_url VARCHAR(500) NOT NULL DEFAULT '',
  link_type ENUM('shop','category','product') NOT NULL DEFAULT 'shop',
  link_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  duration_days INT NOT NULL DEFAULT 7,
  start_at DATETIME NOT NULL,
  end_at DATETIME NOT NULL,
  status ENUM('active','expired','cancelled') NOT NULL DEFAULT 'active',
  pay_source ENUM('wallet','admin') NOT NULL DEFAULT 'wallet',
  wallet_log_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  operator_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_slot_status_time (theme_slot_id, status, start_at, end_at),
  KEY idx_shop (shop_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 默认 4 坑（分类 id 按现有种子：5生鲜 3服饰 4休闲食品 7玩具文创作数码替代）
INSERT INTO homepage_theme_slots (slot_key, position, name, `desc`, cover_url, default_link_type, default_link_id, status, sort)
SELECT 'fresh', 1, '生鲜果蔬专区', '200+生鲜商户当日达', '/uploads/shops/seed/1080.jpg', 'category', 5, 'on', 10
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM homepage_theme_slots WHERE slot_key='fresh');

INSERT INTO homepage_theme_slots (slot_key, position, name, `desc`, cover_url, default_link_type, default_link_id, status, sort)
SELECT 'digital', 2, '数码3C好物', '品牌数码官方直营', '/uploads/shops/seed/1054.jpg', 'category', 7, 'on', 20
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM homepage_theme_slots WHERE slot_key='digital');

INSERT INTO homepage_theme_slots (slot_key, position, name, `desc`, cover_url, default_link_type, default_link_id, status, sort)
SELECT 'fashion', 3, '男女潮流服饰', '上千款平价穿搭', '/uploads/shops/seed/1066.jpg', 'category', 3, 'on', 30
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM homepage_theme_slots WHERE slot_key='fashion');

INSERT INTO homepage_theme_slots (slot_key, position, name, `desc`, cover_url, default_link_type, default_link_id, status, sort)
SELECT 'snack', 4, '零食酒水集合', '全网低价追剧零食', '/uploads/products/1/cookie.png', 'category', 4, 'on', 40
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM homepage_theme_slots WHERE slot_key='snack');

INSERT INTO homepage_theme_packages (theme_slot_id, name, price, duration_days, status, sort, remark)
SELECT 0, '主题坑位·7天', 199.00, 7, 'on', 10, '通用任意坑位'
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM homepage_theme_packages WHERE name='主题坑位·7天');

INSERT INTO homepage_theme_packages (theme_slot_id, name, price, duration_days, status, sort, remark)
SELECT 0, '主题坑位·30天', 499.00, 30, 'on', 20, '通用任意坑位'
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM homepage_theme_packages WHERE name='主题坑位·30天');
