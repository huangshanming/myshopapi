-- 评价管理菜单（平台 + 商家按钮权限在 catalog EnsureShopMenus）
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(110, 1, '评价管理', 'menu', '/admin/reviews', 'admin/Reviews', 'ChatDotRound', 'business:review:list', 35, 1, 1),
(111, 110, '删除违规评价', 'button', '', '', '', 'business:review:delete', 1, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1;

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (110, 111);
