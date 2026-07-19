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
(7,  1, '售后管理', 'menu', '/admin/after-sales', 'admin/after-sales/AfterSales', 'RefreshLeft', 'business:aftersale:list', 16, 1, 1),
(8,  1, '物流管理', 'menu', '/admin/logistics', 'admin/Logistics', 'Van', 'business:logistics:list', 17, 1, 1),
(16, 0, '营销中心', 'dir', '', '', 'Present', '', 30, 1, 1),
(17, 16, '秒杀规则', 'menu', '/admin/seckill/rule', 'admin/SeckillRule', 'Timer', 'marketing:seckill:rule', 31, 1, 1),
(18, 16, '秒杀场次', 'menu', '/admin/seckill/sessions', 'admin/SeckillSessions', 'Calendar', 'marketing:seckill:session', 32, 1, 1),
(5,  1, '全站商品', 'menu', '/admin/products', 'admin/products/Products', 'Goods', 'business:product:list', 14, 1, 1),
(6,  1, '商品分类', 'menu', '/admin/products/categories', 'admin/products/Categories', 'Menu', 'business:category:list', 15, 1, 1),
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
(42, 13, '用户编辑', 'button', '', '', '', 'system:user:edit', 2, 1, 1),
(43, 13, '用户重置密码', 'button', '', '', '', 'system:user:reset', 3, 1, 1),
(44, 13, '用户钱包', 'button', '', '', '', 'system:user:wallet', 4, 1, 1),
(51, 14, '管理员新增', 'button', '', '', '', 'system:admin:add', 1, 1, 1),
(52, 14, '分配角色', 'button', '', '', '', 'system:admin:assign', 2, 1, 1),
(53, 14, '重置密码', 'button', '', '', '', 'system:admin:reset', 3, 1, 1),
(61, 15, '配置保存', 'button', '', '', '', 'system:config:edit', 1, 1, 1),
(71, 3, '店铺新增', 'button', '', '', '', 'business:shop:add', 1, 1, 1),
(72, 3, '店铺编辑', 'button', '', '', '', 'business:shop:edit', 2, 1, 1),
(73, 3, '店主重置密码', 'button', '', '', '', 'business:shop:reset', 3, 1, 1),
(79, 3, '店铺钱包', 'button', '', '', '', 'business:shop:wallet', 4, 1, 1),
(19, 17, '保存秒杀规则', 'button', '', '', '', 'marketing:seckill:rule', 1, 1, 1),
-- 文章管理
(90, 0, '文章管理', 'dir', '', '', 'Document', '', 20, 1, 1),
(91, 90, '文章列表', 'menu', '/admin/articles', 'admin/articles/Articles', 'Document', 'community:article:list', 21, 1, 1),
(92, 90, '分类管理', 'menu', '/admin/articles/categories', 'admin/articles/Categories', 'Menu', 'community:article:category', 22, 1, 1),
(93, 90, '评论管理', 'menu', '/admin/articles/comments', 'admin/articles/Comments', 'ChatDotRound', 'community:article:comment', 23, 1, 1),
(94, 90, '文章回收站', 'menu', '/admin/articles/recycle', 'admin/articles/Recycle', 'Delete', 'community:article:recycle', 24, 1, 1),
(95, 90, '文章统计', 'menu', '/admin/articles/stats', 'admin/articles/Stats', 'DataAnalysis', 'community:article:stats', 25, 1, 1),
(96, 91, '文章新增', 'button', '', '', '', 'community:article:add', 1, 1, 1),
(97, 91, '文章编辑', 'button', '', '', '', 'community:article:edit', 2, 1, 1),
(98, 91, '文章审核', 'button', '', '', '', 'community:article:audit', 3, 1, 1),
(99, 91, '文章置顶', 'button', '', '', '', 'community:article:top', 4, 1, 1),
(74, 5, '商品下架', 'button', '', '', '', 'business:product:off_sale', 1, 1, 1),
(75, 5, '商品删除', 'button', '', '', '', 'business:product:delete', 2, 1, 1),
(76, 6, '分类新增', 'button', '', '', '', 'business:category:add', 1, 1, 1),
(77, 6, '分类编辑', 'button', '', '', '', 'business:category:edit', 2, 1, 1),
(78, 6, '分类删除', 'button', '', '', '', 'business:category:delete', 3, 1, 1),
(81, 4, '订单发货', 'button', '', '', '', 'business:order:ship', 1, 1, 1),
(82, 4, '订单完成', 'button', '', '', '', 'business:order:complete', 2, 1, 1),
(83, 4, '订单备注', 'button', '', '', '', 'business:order:remark', 3, 1, 1),
(84, 7, '售后处理', 'button', '', '', '', 'business:aftersale:handle', 1, 1, 1),
(85, 8, '物流新增', 'button', '', '', '', 'business:logistics:add', 1, 1, 1),
(86, 8, '物流编辑', 'button', '', '', '', 'business:logistics:edit', 2, 1, 1),
(87, 8, '物流删除', 'button', '', '', '', 'business:logistics:delete', 3, 1, 1),
(88, 8, '物流启停', 'button', '', '', '', 'business:logistics:status', 4, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component), icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort);

-- super_admin 拥有全部菜单
INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE status = 1;

-- 绑定种子超管（按手机号，不存在则跳过）
INSERT IGNORE INTO sys_user_role (user_id, role_id)
SELECT u.id, 1 FROM users u WHERE u.mobile = '13900000001' LIMIT 1;

UPDATE users SET role = 'platform_admin' WHERE mobile = '13900000001';
