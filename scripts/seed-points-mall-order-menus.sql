-- 积分商城订单菜单（营销玩法 id=16）
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/seed-points-mall-order-menus.sql
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(129, 16, '积分订单', 'menu', '/admin/points-orders', 'admin/PointsOrders', 'Tickets', 'marketing:points_mall:list', 55, 1, 1),
(145, 129, '处理积分订单', 'button', '', '', '', 'marketing:points_mall:edit', 1, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1, parent_id=VALUES(parent_id);

-- 商品菜单名称更清晰
UPDATE sys_menu SET name='积分商品', sort=50 WHERE id=127;

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (129, 145);
