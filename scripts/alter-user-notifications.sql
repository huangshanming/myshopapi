-- C 端站内信
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-user-notifications.sql
USE mymall;

CREATE TABLE IF NOT EXISTS user_notifications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
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
  KEY idx_user_read_id (user_id, is_read, id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_notification_batches (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(200) NOT NULL DEFAULT '',
  content VARCHAR(1000) NOT NULL DEFAULT '',
  target VARCHAR(16) NOT NULL DEFAULT 'users',
  user_count INT NOT NULL DEFAULT 0,
  success_count INT NOT NULL DEFAULT 0,
  link_type VARCHAR(16) NOT NULL DEFAULT 'none',
  link_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  sender_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
