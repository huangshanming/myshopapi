-- 九宫格抽奖
USE mymall;

CREATE TABLE IF NOT EXISTS lottery_activities (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  title VARCHAR(128) NOT NULL DEFAULT '',
  status TINYINT NOT NULL DEFAULT 0 COMMENT '0草稿 1上线 2下线',
  cost_points INT NOT NULL DEFAULT 0 COMMENT '每次消耗积分，0=免费',
  daily_limit INT NOT NULL DEFAULT 1 COMMENT '每用户每日次数，0=不限制',
  start_at DATETIME NULL,
  end_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_status_time (status, start_at, end_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS lottery_prizes (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  activity_id BIGINT UNSIGNED NOT NULL,
  slot TINYINT NOT NULL COMMENT '0-8',
  name VARCHAR(128) NOT NULL DEFAULT '',
  cover_url VARCHAR(512) NOT NULL DEFAULT '',
  prize_type VARCHAR(32) NOT NULL DEFAULT 'thanks' COMMENT 'points|thanks|physical',
  points_amount INT NOT NULL DEFAULT 0,
  weight INT NOT NULL DEFAULT 0 COMMENT '权重，0不可中',
  stock INT NOT NULL DEFAULT -1 COMMENT '-1无限',
  stock_strict TINYINT NOT NULL DEFAULT 0 COMMENT '1=强控库存 0=不强制',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_activity_slot (activity_id, slot),
  KEY idx_activity (activity_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS lottery_draw_records (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  activity_id BIGINT UNSIGNED NOT NULL,
  prize_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  slot TINYINT NOT NULL DEFAULT 0,
  prize_type VARCHAR(32) NOT NULL DEFAULT '',
  prize_name VARCHAR(128) NOT NULL DEFAULT '',
  points_amount INT NOT NULL DEFAULT 0,
  cost_points INT NOT NULL DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT 'pending|done|failed',
  fulfill_status VARCHAR(16) NOT NULL DEFAULT 'none' COMMENT 'none|need_address|pending|shipped',
  address_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  receiver_name VARCHAR(64) NOT NULL DEFAULT '',
  receiver_phone VARCHAR(32) NOT NULL DEFAULT '',
  receiver_address VARCHAR(512) NOT NULL DEFAULT '',
  ship_company VARCHAR(64) NOT NULL DEFAULT '',
  ship_no VARCHAR(64) NOT NULL DEFAULT '',
  shipped_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_user_activity_day (user_id, activity_id, created_at),
  KEY idx_activity (activity_id),
  KEY idx_fulfill (fulfill_status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
