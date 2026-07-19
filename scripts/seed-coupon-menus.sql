-- 优惠券超管菜单（营销中心 id=16）
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/seed-coupon-menus.sql
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(121, 16, '优惠券', 'menu', '/admin/coupons', 'admin/Coupons', 'Ticket', 'marketing:coupon:list', 30, 1, 1),
(122, 121, '编辑', 'button', '', '', '', 'marketing:coupon:edit', 1, 1, 1),
(123, 121, '发放', 'button', '', '', '', 'marketing:coupon:grant', 2, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1, parent_id=VALUES(parent_id);

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (121, 122, 123);
