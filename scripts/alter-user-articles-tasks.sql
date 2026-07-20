-- C 端用户发文：author_user_id
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-user-articles-tasks.sql
USE mymall;

SET @col := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'community_article' AND COLUMN_NAME = 'author_user_id'
);
SET @sql := IF(@col = 0,
  'ALTER TABLE community_article ADD COLUMN author_user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT ''C端作者用户ID,商户文为0'' AFTER shop_id, ADD INDEX idx_author_user (author_user_id)',
  'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
