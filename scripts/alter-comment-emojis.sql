-- 评论表情包 + 菜单
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-comment-emojis.sql
USE mymall;

CREATE TABLE IF NOT EXISTS community_comment_emojis (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(64) NOT NULL DEFAULT '',
  image_url VARCHAR(500) NOT NULL DEFAULT '',
  sort INT NOT NULL DEFAULT 0,
  status TINYINT NOT NULL DEFAULT 1 COMMENT '1上架 0下架',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_status_sort (status, sort, id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(100, 90, '评论表情', 'menu', '/admin/articles/emojis', 'admin/articles/Emojis', 'Smile', 'community:article:emoji', 26, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1, parent_id=VALUES(parent_id);

INSERT IGNORE INTO sys_role_menu (role_id, menu_id) VALUES (1, 100);
