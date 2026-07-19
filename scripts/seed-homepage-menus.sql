-- 首页展位超管菜单
-- 导入时务必指定 utf8mb4，例如：
--   mysql -uhomestead -p --default-character-set=utf8mb4 mymall < scripts/seed-homepage-menus.sql
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(112, 1, '首页展位', 'menu', '/admin/homepage', 'admin/HomepageSlots', 'Promotion', 'business:homepage:list', 36, 1, 1),
(113, 112, '套餐管理', 'button', '', '', '', 'business:homepage:package', 1, 1, 1),
(114, 112, '代开通', 'button', '', '', '', 'business:homepage:grant', 2, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1;

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (112, 113, 114);
