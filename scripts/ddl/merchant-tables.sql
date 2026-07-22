-- goctl model DDL for merchant-service (extracted; authority remains scripts/*.sql)

CREATE TABLE IF NOT EXISTS shops (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    name VARCHAR(100) NOT NULL COMMENT '店铺名称',
    logo VARCHAR(500) DEFAULT '' COMMENT 'Logo',
    contact_name VARCHAR(50) DEFAULT '' COMMENT '联系人',
    contact_phone CHAR(11) DEFAULT '' COMMENT '联系电话',
    description VARCHAR(500) DEFAULT '' COMMENT '简介',
    category VARCHAR(50) DEFAULT '' COMMENT '经营类目',
    province VARCHAR(50) DEFAULT '',
    city VARCHAR(50) DEFAULT '',
    district VARCHAR(50) DEFAULT '',
    address VARCHAR(255) DEFAULT '' COMMENT '详细地址',
    business_license_no VARCHAR(64) DEFAULT '' COMMENT '营业执照号',
    legal_person VARCHAR(50) DEFAULT '' COMMENT '法人',
    license_image VARCHAR(500) DEFAULT '' COMMENT '执照图',
    storefront_image VARCHAR(500) DEFAULT '' COMMENT '门头图',
    owner_user_id BIGINT UNSIGNED NOT NULL COMMENT '店主用户ID',
    status ENUM('pending','approved','rejected','disabled') NOT NULL DEFAULT 'pending' COMMENT '状态',
    reject_reason VARCHAR(255) DEFAULT '' COMMENT '拒绝原因',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_owner (owner_user_id),
    INDEX idx_status (status),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_applications (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '申请人',
    shop_name VARCHAR(100) NOT NULL,
    contact_name VARCHAR(50) NOT NULL,
    contact_phone CHAR(11) NOT NULL,
    description VARCHAR(500) DEFAULT '',
    category VARCHAR(50) DEFAULT '',
    province VARCHAR(50) DEFAULT '',
    city VARCHAR(50) DEFAULT '',
    district VARCHAR(50) DEFAULT '',
    address VARCHAR(255) DEFAULT '',
    business_license_no VARCHAR(64) DEFAULT '',
    legal_person VARCHAR(50) DEFAULT '',
    license_image VARCHAR(500) DEFAULT '',
    storefront_image VARCHAR(500) DEFAULT '',
    status ENUM('pending','approved','rejected') NOT NULL DEFAULT 'pending',
    reject_reason VARCHAR(255) DEFAULT '',
    reviewed_by BIGINT UNSIGNED DEFAULT NULL,
    reviewed_at TIMESTAMP NULL DEFAULT NULL,
    shop_id BIGINT UNSIGNED DEFAULT NULL COMMENT '审核通过后关联店铺',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_status (status),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_members (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    shop_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    member_role ENUM('owner','staff') NOT NULL DEFAULT 'owner',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_shop_user (shop_id, user_id),
    INDEX idx_user (user_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_wallets (
    shop_id BIGINT UNSIGNED NOT NULL  COMMENT '店铺ID',
    balance DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '可用余额',
    frozen_balance DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '冻结余额',
    deposit DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '保证金',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (shop_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='店铺钱包';

CREATE TABLE IF NOT EXISTS shop_wallet_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
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
    KEY idx_shop_created (shop_id, created_at),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='钱包流水';

CREATE TABLE IF NOT EXISTS seckill_rules (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    duration_hours INT NOT NULL DEFAULT 24 COMMENT '场次时长(小时)',
    apply_fee DECIMAL(12,2) NOT NULL DEFAULT 10.00 COMMENT '报名费',
    max_entries_per_shop INT NOT NULL DEFAULT 5 COMMENT '每店每场最多报名数',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '1启用 0停用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀规则';

CREATE TABLE IF NOT EXISTS seckill_sessions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    rule_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    start_at DATETIME NOT NULL,
    end_at DATETIME NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'active' COMMENT 'active/ended',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_status_end (status, end_at),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀场次';

CREATE TABLE IF NOT EXISTS seckill_entries (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
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
    IF NOT EXISTS auto_renew TINYINT NOT NULL DEFAULT 0 COMMENT '1=场次到期自动续报下一场',
    auto_renew TINYINT NOT NULL DEFAULT 0 COMMENT '1=场次到期自动续报下一场',
    KEY idx_session (session_id),
    KEY idx_shop_session (shop_id, session_id),
    UNIQUE KEY uk_session_shop_product (session_id, shop_id, product_id),
    PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀报名';;

CREATE TABLE IF NOT EXISTS homepage_slot_packages (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    slot_type VARCHAR(32) NOT NULL COMMENT 'brand_shop/quality_shop/article',
    name VARCHAR(100) NOT NULL,
    price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    duration_days INT NOT NULL DEFAULT 1,
    status VARCHAR(16) NOT NULL DEFAULT 'on' COMMENT 'on/off',
    sort INT NOT NULL DEFAULT 0,
    remark VARCHAR(255) NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_type_status (slot_type, status),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='首页展位套餐';

CREATE TABLE IF NOT EXISTS homepage_slot_settings (
    slot_type VARCHAR(32) NOT NULL ,
    home_limit INT NOT NULL DEFAULT 8 COMMENT '首页条带最多展示条数',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (slot_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='首页展位展示设置';

CREATE TABLE IF NOT EXISTS homepage_slot_orders (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
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
    KEY idx_shop (shop_id, slot_type, status),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='首页展位订单';

CREATE TABLE IF NOT EXISTS homepage_theme_slots (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
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
  UNIQUE KEY uk_position (position),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS homepage_theme_packages (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
  theme_slot_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=通用任意坑位',
  name VARCHAR(100) NOT NULL,
  price DECIMAL(12,2) NOT NULL DEFAULT 0,
  duration_days INT NOT NULL DEFAULT 7,
  status ENUM('on','off') NOT NULL DEFAULT 'on',
  sort INT NOT NULL DEFAULT 0,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_slot_status (theme_slot_id, status),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS homepage_theme_orders (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
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
  KEY idx_shop (shop_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coupons (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
  name VARCHAR(100) NOT NULL DEFAULT '',
  issuer_type ENUM('platform','shop') NOT NULL DEFAULT 'platform',
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  coupon_type ENUM('full_reduce','no_threshold','category','product','discount') NOT NULL DEFAULT 'full_reduce',
  threshold_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  discount_rate DECIMAL(6,4) NOT NULL DEFAULT 0 COMMENT '0.80=八折',
  max_discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  scope_type ENUM('all','category','product') NOT NULL DEFAULT 'all',
  total_count INT NOT NULL DEFAULT 0 COMMENT '0=不限',
  claimed_count INT NOT NULL DEFAULT 0,
  per_user_limit INT NOT NULL DEFAULT 1,
  valid_type ENUM('fixed','relative') NOT NULL DEFAULT 'fixed',
  valid_start DATETIME NULL,
  valid_end DATETIME NULL,
  valid_days INT NOT NULL DEFAULT 0,
  stackable TINYINT NOT NULL DEFAULT 0,
  user_identity ENUM('all','new','old') NOT NULL DEFAULT 'all',
  channels JSON NULL,
  status ENUM('draft','on','off','expired') NOT NULL DEFAULT 'draft',
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_issuer (issuer_type, shop_id, status),
  KEY idx_status (status),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coupon_scopes (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
  coupon_id BIGINT UNSIGNED NOT NULL,
  ref_type ENUM('category','product') NOT NULL,
  ref_id BIGINT UNSIGNED NOT NULL,
  KEY idx_coupon (coupon_id),
  UNIQUE KEY uk_coupon_ref (coupon_id, ref_type, ref_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_coupons (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
  coupon_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  status ENUM('unused','locked','used','expired') NOT NULL DEFAULT 'unused',
  source ENUM('direct','order_gift','popup','targeted') NOT NULL DEFAULT 'direct',
  valid_start DATETIME NOT NULL,
  valid_end DATETIME NOT NULL,
  order_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  locked_at DATETIME NULL,
  used_at DATETIME NULL,
  claim_batch_no VARCHAR(64) NOT NULL DEFAULT '',
  discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0 COMMENT '核销时实际抵扣',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_user_status (user_id, status),
  KEY idx_coupon_user (coupon_id, user_id),
  KEY idx_order (order_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coupon_grants (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
  coupon_id BIGINT UNSIGNED NOT NULL,
  operator_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  issuer_type ENUM('platform','shop') NOT NULL DEFAULT 'platform',
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  user_count INT NOT NULL DEFAULT 0,
  success_count INT NOT NULL DEFAULT 0,
  batch_no VARCHAR(64) NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_coupon (coupon_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coupon_redeem_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
  user_coupon_id BIGINT UNSIGNED NOT NULL,
  coupon_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  order_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  action ENUM('redeem','unlock','return') NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_coupon (coupon_id),
  KEY idx_order (order_id),
  KEY idx_user (user_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

