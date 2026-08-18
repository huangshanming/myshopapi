-- 增量：九宫格抽奖菜单（已跑过 seed-rbac 的环境可单独执行）
-- 注意：id=100 已被「评论表情」占用，抽奖使用 146；记录 148；实物订单 147
USE mymall;
INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status)
VALUES (146, 16, '九宫格抽奖', 'menu', '/admin/lottery', 'admin/Lottery', 'Present', 'marketing:lottery:list', 36, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), parent_id=VALUES(parent_id), path=VALUES(path), component=VALUES(component), perms=VALUES(perms), sort=VALUES(sort), icon=VALUES(icon), visible=VALUES(visible), status=VALUES(status);
INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status)
VALUES (148, 16, '抽奖记录', 'menu', '/admin/lottery-records', 'admin/LotteryRecords', 'Tickets', 'marketing:lottery:list', 37, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), parent_id=VALUES(parent_id), path=VALUES(path), component=VALUES(component), perms=VALUES(perms), sort=VALUES(sort), icon=VALUES(icon), visible=VALUES(visible), status=VALUES(status);
INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status)
VALUES (147, 16, '抽奖实物订单', 'menu', '/admin/lottery-orders', 'admin/LotteryOrders', 'Van', 'marketing:lottery:order', 38, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), parent_id=VALUES(parent_id), path=VALUES(path), component=VALUES(component), perms=VALUES(perms), sort=VALUES(sort), icon=VALUES(icon), visible=VALUES(visible), status=VALUES(status);
INSERT IGNORE INTO sys_role_menu (role_id, menu_id) VALUES (1, 146), (1, 147), (1, 148);
