-- goctl model DDL for user-service (extracted; authority remains scripts/*.sql)

CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    mobile CHAR(11) NOT NULL COMMENT '登录手机号',
    password VARCHAR(255) NOT NULL COMMENT '登录密码',
    nickname VARCHAR(50) NOT NULL COMMENT '用户昵称',
    avatar VARCHAR(255) NOT NULL DEFAULT '' COMMENT '用户头像URL',
	gender TINYINT NOT NULL DEFAULT 0 COMMENT '性别：0-未知 1-男 2-女',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '账号状态：1-正常 0-禁用',
    role VARCHAR(32) NOT NULL DEFAULT 'user' COMMENT '角色',
    last_login_time TIMESTAMP NULL DEFAULT NULL COMMENT '最后登录时间',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    UNIQUE KEY uk_mobile (mobile),
    INDEX idx_deleted_at (deleted_at),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_addresses (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    user_id BIGINT UNSIGNED NOT NULL,
    receiver_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '收货人',
    receiver_phone VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号',
    province VARCHAR(64) NOT NULL DEFAULT '',
    city VARCHAR(64) NOT NULL DEFAULT '',
    district VARCHAR(64) NOT NULL DEFAULT '',
    detail VARCHAR(255) NOT NULL DEFAULT '' COMMENT '详细地址',
    province_code VARCHAR(12) NOT NULL DEFAULT '',
    city_code VARCHAR(12) NOT NULL DEFAULT '',
    district_code VARCHAR(12) NOT NULL DEFAULT '',
    is_default TINYINT NOT NULL DEFAULT 0 COMMENT '1默认',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_user (user_id),
    KEY idx_user_default (user_id, is_default),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收货地址';

CREATE TABLE IF NOT EXISTS user_wallets (
    user_id BIGINT UNSIGNED NOT NULL  COMMENT '用户ID',
    balance DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '可用余额',
    frozen_balance DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '冻结余额',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户钱包';

CREATE TABLE IF NOT EXISTS user_wallet_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
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
    KEY idx_ref (ref_type, ref_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户钱包流水';

CREATE TABLE IF NOT EXISTS user_points (
    user_id BIGINT UNSIGNED NOT NULL  COMMENT '用户ID',
    points BIGINT NOT NULL DEFAULT 0 COMMENT '当前积分',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分账户';

CREATE TABLE IF NOT EXISTS user_point_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    user_id BIGINT UNSIGNED NOT NULL,
    change_type VARCHAR(32) NOT NULL DEFAULT 'task_claim' COMMENT 'task_claim/admin_adjust',
    delta INT NOT NULL COMMENT '变动积分(可负)',
    points_after BIGINT NOT NULL DEFAULT 0,
    remark VARCHAR(255) NOT NULL DEFAULT '',
    ref_type VARCHAR(32) NOT NULL DEFAULT '',
    ref_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_user_created (user_id, created_at),
    KEY idx_ref (ref_type, ref_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分流水';

CREATE TABLE IF NOT EXISTS task_definitions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    code VARCHAR(64) NOT NULL COMMENT '任务编码',
    title VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL DEFAULT '',
    icon VARCHAR(64) NOT NULL DEFAULT '',
    period VARCHAR(16) NOT NULL DEFAULT 'daily' COMMENT 'daily|once',
    enabled TINYINT NOT NULL DEFAULT 1,
    reward_points INT NOT NULL DEFAULT 0 COMMENT '0=不奖励积分',
    target_count INT NOT NULL DEFAULT 1 COMMENT '完成所需进度',
    daily_limit INT NOT NULL DEFAULT 0 COMMENT '每日最多领取次数,0=不限; once任务忽略',
    sort INT NOT NULL DEFAULT 0,
    rules_json VARCHAR(1000) NOT NULL DEFAULT '{}',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code (code),
    KEY idx_enabled_sort (enabled, sort),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务定义';

CREATE TABLE IF NOT EXISTS user_task_progress (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    user_id BIGINT UNSIGNED NOT NULL,
    task_code VARCHAR(64) NOT NULL,
    biz_date DATE NOT NULL COMMENT '日任务按自然日; once=1970-01-01',
    progress INT NOT NULL DEFAULT 0,
    claim_count INT NOT NULL DEFAULT 0 COMMENT '当日已领取次数',
    status VARCHAR(16) NOT NULL DEFAULT 'ongoing' COMMENT 'ongoing|claimable|claimed',
    claimed_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_task_date (user_id, task_code, biz_date),
    KEY idx_user_status (user_id, status),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户任务进度';

CREATE TABLE IF NOT EXISTS user_task_dedupe (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    user_id BIGINT UNSIGNED NOT NULL,
    task_code VARCHAR(64) NOT NULL,
    biz_date DATE NOT NULL,
    ref_key VARCHAR(64) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_dedupe (user_id, task_code, biz_date, ref_key),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务事件去重';

CREATE TABLE IF NOT EXISTS points_products (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
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
    KEY idx_name (name),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分商城商品';

CREATE TABLE IF NOT EXISTS points_exchange_orders (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
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
    KEY idx_product (product_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分商城兑换订单';

CREATE TABLE IF NOT EXISTS user_notifications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
  user_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(200) NOT NULL DEFAULT '',
  content VARCHAR(1000) NOT NULL DEFAULT '',
  msg_type ENUM('announce','order','system') NOT NULL DEFAULT 'system',
  link_type ENUM('none','order','url') NOT NULL DEFAULT 'none',
  link_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  extra JSON NULL,
  is_read TINYINT NOT NULL DEFAULT 0,
  sender_type ENUM('admin','system') NOT NULL DEFAULT 'system',
  sender_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  batch_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_user_read_id (user_id, is_read, id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_notification_batches (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
  title VARCHAR(200) NOT NULL DEFAULT '',
  content VARCHAR(1000) NOT NULL DEFAULT '',
  target VARCHAR(16) NOT NULL DEFAULT 'users',
  user_count INT NOT NULL DEFAULT 0,
  success_count INT NOT NULL DEFAULT 0,
  link_type VARCHAR(16) NOT NULL DEFAULT 'none',
  link_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  sender_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_created (created_at),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS regions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    code VARCHAR(12) NOT NULL COMMENT '行政区划代码',
    name VARCHAR(64) NOT NULL COMMENT '名称',
    parent_code VARCHAR(12) NOT NULL DEFAULT '' COMMENT '父级代码，省级为空',
    level TINYINT NOT NULL COMMENT '1省 2市 3区县',
    sort INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code (code),
    KEY idx_parent (parent_code),
    KEY idx_level (level),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='全国省市区';

CREATE TABLE IF NOT EXISTS sys_menu (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    name VARCHAR(64) NOT NULL DEFAULT '',
    type VARCHAR(16) NOT NULL DEFAULT 'menu' COMMENT 'dir|menu|button',
    path VARCHAR(128) NOT NULL DEFAULT '',
    component VARCHAR(128) NOT NULL DEFAULT '',
    icon VARCHAR(64) NOT NULL DEFAULT '',
    perms VARCHAR(128) NOT NULL DEFAULT '' COMMENT '权限码，空则仅导航',
    sort INT NOT NULL DEFAULT 0,
    visible TINYINT NOT NULL DEFAULT 1,
    status TINYINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_sys_menu_parent (parent_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_role (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    code VARCHAR(64) NOT NULL,
    name VARCHAR(64) NOT NULL,
    status TINYINT NOT NULL DEFAULT 1,
    remark VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_sys_role_code (code),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_role_menu (
    role_id BIGINT UNSIGNED NOT NULL,
    menu_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (role_id, menu_id),
    KEY idx_sys_role_menu_menu (menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_user_role (
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, role_id),
    KEY idx_sys_user_role_role (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sys_config (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT ,
    config_key VARCHAR(64) NOT NULL,
    config_value VARCHAR(512) NOT NULL DEFAULT '',
    remark VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_sys_config_key (config_key),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

