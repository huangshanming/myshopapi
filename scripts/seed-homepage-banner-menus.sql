-- 首页 Banner 超管菜单（挂营销中心 id=16）
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/seed-homepage-banner-menus.sql
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(115, 16, '首页 Banner', 'menu', '/admin/banners', 'admin/Banners', 'Picture', 'marketing:banner:list', 15, 1, 1),
(116, 115, '新增编辑', 'button', '', '', '', 'marketing:banner:edit', 1, 1, 1),
(117, 115, '删除', 'button', '', '', '', 'marketing:banner:delete', 2, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1, parent_id=VALUES(parent_id);

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (115, 116, 117);
