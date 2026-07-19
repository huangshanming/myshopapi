-- 文章评论多级回复字段 + 通知 link_type 支持 article
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-article-comments-reply.sql
USE mymall;

ALTER TABLE community_article_comment
  ADD COLUMN IF NOT EXISTS parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父评论ID，0=一级' AFTER user_id,
  ADD COLUMN IF NOT EXISTS root_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '根评论ID' AFTER parent_id,
  ADD COLUMN IF NOT EXISTS reply_to_user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '回复目标用户' AFTER root_id;

-- MySQL 8.0.12 前无 IF NOT EXISTS，兼容写法：
-- 若上面失败可手动执行：
-- ALTER TABLE community_article_comment
--   ADD COLUMN parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER user_id,
--   ADD COLUMN root_id BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER parent_id,
--   ADD COLUMN reply_to_user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER root_id;

ALTER TABLE community_article_comment
  ADD INDEX idx_article_root (article_id, root_id, id);

ALTER TABLE community_article
  ADD COLUMN IF NOT EXISTS comment_count INT NOT NULL DEFAULT 0 AFTER collect_count;

ALTER TABLE user_notifications
  MODIFY COLUMN link_type ENUM('none','order','url','article') NOT NULL DEFAULT 'none';
