-- RBAC 种子：超管角色 + 菜单树 + 绑定 13900000001
USE mymall;

INSERT INTO sys_role (id, code, name, status, remark)
VALUES (1, 'super_admin', '超级管理员', 1, '拥有全部权限')
ON DUPLICATE KEY UPDATE name=VALUES(name), status=1;

INSERT INTO sys_config (config_key, config_value, remark)
VALUES ('site_name', 'mymall 管理后台', '站点名称')
ON DUPLICATE KEY UPDATE config_value=VALUES(config_value);

-- 菜单：业务 + 系统管理（固定 id 便于角色关联）
INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(1,  0, '业务管理', 'dir',  '', '', 'Menu', '', 10, 1, 1),
(2,  1, '入驻审核', 'menu', '/admin/applications', 'admin/Applications', 'Document', 'business:application:list', 11, 1, 1),
(3,  1, '店铺管理', 'menu', '/admin/shops', 'admin/Shops', 'Shop', 'business:shop:list', 12, 1, 1),
(4,  1, '全站订单', 'menu', '/admin/orders', 'admin/Orders', 'List', 'business:order:list', 13, 1, 1),
(10, 0, '系统管理', 'dir',  '', '', 'Setting', '', 90, 1, 1),
(11, 10, '菜单管理', 'menu', '/admin/system/menus', 'admin/system/Menus', 'Menu', 'system:menu:list', 91, 1, 1),
(12, 10, '角色管理', 'menu', '/admin/system/roles', 'admin/system/Roles', 'UserFilled', 'system:role:list', 92, 1, 1),
(13, 10, '用户管理', 'menu', '/admin/system/users', 'admin/system/Users', 'User', 'system:user:list', 93, 1, 1),
(14, 10, '管理员设置', 'menu', '/admin/system/admins', 'admin/system/Admins', 'Avatar', 'system:admin:list', 94, 1, 1),
(15, 10, '系统设置', 'menu', '/admin/system/configs', 'admin/system/Configs', 'Tools', 'system:config:list', 95, 1, 1),
-- 按钮权限
(21, 11, '菜单新增', 'button', '', '', '', 'system:menu:add', 1, 1, 1),
(22, 11, '菜单编辑', 'button', '', '', '', 'system:menu:edit', 2, 1, 1),
(23, 11, '菜单删除', 'button', '', '', '', 'system:menu:delete', 3, 1, 1),
(31, 12, '角色新增', 'button', '', '', '', 'system:role:add', 1, 1, 1),
(32, 12, '角色编辑', 'button', '', '', '', 'system:role:edit', 2, 1, 1),
(33, 12, '角色删除', 'button', '', '', '', 'system:role:delete', 3, 1, 1),
(34, 12, '分配菜单', 'button', '', '', '', 'system:role:assign', 4, 1, 1),
(41, 13, '用户启停', 'button', '', '', '', 'system:user:status', 1, 1, 1),
(51, 14, '管理员新增', 'button', '', '', '', 'system:admin:add', 1, 1, 1),
(52, 14, '分配角色', 'button', '', '', '', 'system:admin:assign', 2, 1, 1),
(53, 14, '重置密码', 'button', '', '', '', 'system:admin:reset', 3, 1, 1),
(61, 15, '配置保存', 'button', '', '', '', 'system:config:edit', 1, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), perms=VALUES(perms), sort=VALUES(sort);

-- super_admin 拥有全部菜单
INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE status = 1;

-- 绑定种子超管（按手机号，不存在则跳过）
INSERT IGNORE INTO sys_user_role (user_id, role_id)
SELECT u.id, 1 FROM users u WHERE u.mobile = '13900000001' LIMIT 1;

UPDATE users SET role = 'platform_admin' WHERE mobile = '13900000001';
