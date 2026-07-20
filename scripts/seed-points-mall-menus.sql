-- 积分商城菜单（营销中心 id=16）
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/seed-points-mall-menus.sql
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(127, 16, '积分商城', 'menu', '/admin/points-mall', 'admin/PointsMall', 'Present', 'marketing:points_mall:list', 46, 1, 1),
(128, 127, '编辑积分商品', 'button', '', '', '', 'marketing:points_mall:edit', 1, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1, parent_id=VALUES(parent_id);

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (127, 128);
