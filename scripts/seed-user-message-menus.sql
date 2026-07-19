-- 超管用户消息菜单（业务管理 id=1）
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/seed-user-message-menus.sql
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(124, 1, '用户消息', 'menu', '/admin/messages', 'admin/Messages', 'Bell', 'business:message:send', 55, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1, parent_id=VALUES(parent_id);

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (124);
