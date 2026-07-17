-- mymall 业务表（users / products / product_categories）
-- 执行：mysql -u homestead -p mymall < scripts/init-schema.sql

USE mymall;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
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
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_categories (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    parent_id INT NOT NULL DEFAULT 0 COMMENT '父分类ID',
    name VARCHAR(100) NOT NULL COMMENT '分类名称',
    icon VARCHAR(200) DEFAULT NULL COMMENT '分类图标',
    image VARCHAR(500) DEFAULT NULL COMMENT '分类图片',
    description VARCHAR(500) DEFAULT NULL COMMENT '分类描述',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序',
    level INT NOT NULL DEFAULT 1 COMMENT '分类层级',
    is_show TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否显示',
    product_count INT NOT NULL DEFAULT 0 COMMENT '商品数量',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属店铺ID',
    product_no VARCHAR(50) NOT NULL COMMENT '商品编号',
    name VARCHAR(200) NOT NULL COMMENT '商品名称',
    subtitle VARCHAR(500) DEFAULT NULL COMMENT '商品副标题',
    description TEXT DEFAULT NULL COMMENT '商品描述',
    main_image VARCHAR(500) DEFAULT NULL COMMENT '主图',
    image_list JSON DEFAULT NULL COMMENT '商品图片列表',
    video_url VARCHAR(500) DEFAULT NULL COMMENT '商品视频',
    market_price DECIMAL(10,2) DEFAULT NULL COMMENT '市场价',
    sale_price DECIMAL(10,2) NOT NULL COMMENT '销售价',
    cost_price DECIMAL(10,2) DEFAULT NULL COMMENT '成本价',
    discount DECIMAL(5,2) NOT NULL DEFAULT 100.00 COMMENT '折扣(百分比)',
    discount_price DECIMAL(10,2) DEFAULT NULL COMMENT '折后价',
    stock INT NOT NULL DEFAULT 0 COMMENT '库存数量',
    stock_warn INT NOT NULL DEFAULT 10 COMMENT '库存预警值',
    sold_count INT NOT NULL DEFAULT 0 COMMENT '已售数量',
    view_count INT NOT NULL DEFAULT 0 COMMENT '浏览数量',
    collect_count INT NOT NULL DEFAULT 0 COMMENT '收藏数量',
    pet_type ENUM('dog','cat','both','other') NOT NULL DEFAULT 'both' COMMENT '宠物类型',
    pet_age JSON DEFAULT NULL COMMENT '适用年龄',
    pet_size JSON DEFAULT NULL COMMENT '适用体型',
    weight DECIMAL(8,2) DEFAULT NULL COMMENT '重量(kg)',
    unit VARCHAR(20) DEFAULT NULL COMMENT '单位',
    brand_id INT DEFAULT NULL COMMENT '品牌ID',
    category_id INT NOT NULL COMMENT '分类ID',
    product_type ENUM('physical','fresh','virtual') NOT NULL DEFAULT 'physical' COMMENT '商品类型',
    spec_json JSON DEFAULT NULL COMMENT '规格定义',
    tags JSON DEFAULT NULL COMMENT '标签数组',
    nutrition_info JSON DEFAULT NULL COMMENT '营养成分表',
    ingredients TEXT DEFAULT NULL COMMENT '主要成分',
    feeding_guide TEXT DEFAULT NULL COMMENT '喂养指南',
    shelf_life INT DEFAULT NULL COMMENT '保质期(天)',
    storage_condition VARCHAR(200) DEFAULT NULL COMMENT '存储条件',
    status ENUM('draft','on_sale','off_sale','deleted') NOT NULL DEFAULT 'draft' COMMENT '状态',
    is_hot TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否热销',
    is_new TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否新品',
    is_recommend TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否推荐',
    is_prescription TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否处方粮',
    is_imported TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否进口',
    is_organic TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否有机',
    is_grain_free TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否无谷',
    publish_time DATETIME DEFAULT NULL COMMENT '上架时间',
    schedule_on_at DATETIME DEFAULT NULL COMMENT '定时上架',
    schedule_off_at DATETIME DEFAULT NULL COMMENT '定时下架',
    copy_from_id BIGINT UNSIGNED DEFAULT NULL COMMENT '复制来源',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    UNIQUE KEY uk_product_no (product_no),
    INDEX idx_category_id (category_id),
    INDEX idx_brand_id (brand_id),
    INDEX idx_pet_type (pet_type),
    INDEX idx_status (status),
    INDEX idx_is_hot (is_hot),
    INDEX idx_is_new (is_new),
    INDEX idx_is_recommend (is_recommend),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_shop_id (shop_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ========== 商品中台关联表（与 alter-product-center.sql 对齐）==========
CREATE TABLE IF NOT EXISTS product_skus (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  product_id BIGINT UNSIGNED NOT NULL,
  shop_id BIGINT UNSIGNED NOT NULL,
  sku_no VARCHAR(64) NOT NULL,
  spec_values JSON NOT NULL,
  spec_key VARCHAR(255) NOT NULL,
  sale_price DECIMAL(10,2) NOT NULL,
  market_price DECIMAL(10,2) NULL,
  cost_price DECIMAL(10,2) NULL,
  stock INT NOT NULL DEFAULT 0,
  stock_warn INT NOT NULL DEFAULT 10,
  barcode VARCHAR(64) NULL,
  status ENUM('enabled','disabled') NOT NULL DEFAULT 'enabled',
  sold_count INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL DEFAULT NULL,
  UNIQUE KEY uk_sku_no (sku_no),
  UNIQUE KEY uk_product_spec (product_id, spec_key),
  KEY idx_shop_stock (shop_id, stock),
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_images (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  product_id BIGINT UNSIGNED NOT NULL,
  shop_id BIGINT UNSIGNED NOT NULL,
  url VARCHAR(500) NOT NULL,
  typ ENUM('main','gallery','detail') NOT NULL DEFAULT 'gallery',
  sort INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_tags (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  name VARCHAR(50) NOT NULL,
  color VARCHAR(20) DEFAULT '',
  status TINYINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_shop_name (shop_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_tag_rels (
  product_id BIGINT UNSIGNED NOT NULL,
  tag_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (product_id, tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_attr_templates (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  name VARCHAR(100) NOT NULL,
  attrs_json JSON NOT NULL,
  status TINYINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_attrs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  product_id BIGINT UNSIGNED NOT NULL,
  template_id BIGINT UNSIGNED NULL,
  attr_key VARCHAR(64) NOT NULL,
  attr_label VARCHAR(100) NOT NULL,
  attr_value VARCHAR(500) NOT NULL,
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_schedules (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  product_id BIGINT UNSIGNED NOT NULL,
  shop_id BIGINT UNSIGNED NOT NULL,
  action ENUM('on_sale','off_sale') NOT NULL,
  run_at DATETIME NOT NULL,
  status ENUM('pending','done','cancelled') NOT NULL DEFAULT 'pending',
  locked_at DATETIME NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_due (status, run_at),
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_batch_jobs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  job_type VARCHAR(32) NOT NULL,
  payload_json JSON NOT NULL,
  progress INT NOT NULL DEFAULT 0,
  total INT NOT NULL DEFAULT 0,
  status ENUM('pending','running','success','failed','partial') NOT NULL DEFAULT 'pending',
  result_msg VARCHAR(1000) DEFAULT '',
  operator_id BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_shop (shop_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_op_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  product_id BIGINT UNSIGNED NULL,
  operator_id BIGINT UNSIGNED NOT NULL,
  action VARCHAR(64) NOT NULL,
  before_json JSON NULL,
  after_json JSON NULL,
  ip VARCHAR(64) DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_shop_time (shop_id, created_at),
  KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_roles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(100) NOT NULL,
  status TINYINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_shop_code (shop_id, code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_menus (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  name VARCHAR(100) NOT NULL,
  type VARCHAR(20) NOT NULL DEFAULT 'menu',
  path VARCHAR(200) DEFAULT '',
  component VARCHAR(200) DEFAULT '',
  icon VARCHAR(50) DEFAULT '',
  perms VARCHAR(100) DEFAULT '',
  sort INT NOT NULL DEFAULT 0,
  visible TINYINT NOT NULL DEFAULT 1,
  status TINYINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_role_menus (
  role_id BIGINT UNSIGNED NOT NULL,
  menu_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (role_id, menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shop_user_roles (
  shop_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  role_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (shop_id, user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- 社区文章相关表（幂等）
CREATE TABLE IF NOT EXISTS community_article_category (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  name VARCHAR(100) NOT NULL,
  sort INT NOT NULL DEFAULT 0,
  status TINYINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_parent (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS community_article (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  category_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  title VARCHAR(200) NOT NULL,
  cover_url VARCHAR(500) DEFAULT '',
  content MEDIUMTEXT,
  audit_status ENUM('pending','approved','rejected') NOT NULL DEFAULT 'pending',
  reject_reason VARCHAR(500) DEFAULT '',
  status ENUM('draft','scheduled','published','offline','deleted') NOT NULL DEFAULT 'draft',
  schedule_publish_at DATETIME NULL,
  is_top TINYINT NOT NULL DEFAULT 0,
  view_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
  like_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
  published_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_shop_status (shop_id, status),
  KEY idx_audit (audit_status),
  KEY idx_schedule (schedule_publish_at, status),
  KEY idx_top_pub (is_top, published_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS community_article_img (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  article_id BIGINT UNSIGNED NOT NULL,
  shop_id BIGINT UNSIGNED NOT NULL,
  url VARCHAR(500) NOT NULL,
  sort INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_article (article_id),
  KEY idx_shop (shop_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS community_article_comment (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  article_id BIGINT UNSIGNED NOT NULL,
  shop_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  content VARCHAR(1000) NOT NULL,
  status ENUM('visible','hidden','deleted') NOT NULL DEFAULT 'visible',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_article (article_id),
  KEY idx_shop (shop_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 商家消息通知表（幂等）
CREATE TABLE IF NOT EXISTS shop_notifications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  type VARCHAR(32) NOT NULL,
  title VARCHAR(200) NOT NULL,
  content VARCHAR(1000) NOT NULL DEFAULT '',
  link VARCHAR(255) DEFAULT '',
  ref_type VARCHAR(32) DEFAULT '',
  ref_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  is_read TINYINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_shop_read_time (shop_id, is_read, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
