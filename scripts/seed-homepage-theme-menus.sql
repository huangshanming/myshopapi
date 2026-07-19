-- 主题集市超管菜单（营销中心 id=16）
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/seed-homepage-theme-menus.sql
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(118, 16, '主题集市', 'menu', '/admin/themes', 'admin/ThemeSlots', 'Grid', 'marketing:theme:list', 20, 1, 1),
(119, 118, '坑位套餐', 'button', '', '', '', 'marketing:theme:package', 1, 1, 1),
(120, 118, '代开通', 'button', '', '', '', 'marketing:theme:grant', 2, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1, parent_id=VALUES(parent_id);

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (118, 119, 120);
