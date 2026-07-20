-- 任务中心菜单（营销中心 id=16）
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/seed-task-center-menus.sql
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(125, 16, '任务中心', 'menu', '/admin/tasks', 'admin/TaskCenter', 'Trophy', 'marketing:task:list', 45, 1, 1),
(126, 125, '编辑任务', 'button', '', '', '', 'marketing:task:edit', 1, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1, parent_id=VALUES(parent_id);

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (125, 126);
