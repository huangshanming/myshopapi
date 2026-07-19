-- 首页 Banner 表（幂等）
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-homepage-banners.sql
USE mymall;

CREATE TABLE IF NOT EXISTS homepage_banners (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(100) NOT NULL DEFAULT '',
  image_url VARCHAR(500) NOT NULL DEFAULT '',
  link_type ENUM('none','product','article') NOT NULL DEFAULT 'none',
  link_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  sort INT NOT NULL DEFAULT 0,
  status ENUM('on','off') NOT NULL DEFAULT 'on',
  start_at DATETIME NULL,
  end_at DATETIME NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_status_sort (status, sort),
  KEY idx_time (start_at, end_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
