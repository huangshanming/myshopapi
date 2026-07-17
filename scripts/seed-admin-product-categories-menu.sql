-- 已有库增量：商品分类菜单
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(6, 1, '商品分类', 'menu', '/admin/products/categories', 'admin/products/Categories', 'Menu', 'business:category:list', 15, 1, 1),
(76, 6, '分类新增', 'button', '', '', '', 'business:category:add', 1, 1, 1),
(77, 6, '分类编辑', 'button', '', '', '', 'business:category:edit', 2, 1, 1),
(78, 6, '分类删除', 'button', '', '', '', 'business:category:delete', 3, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1;

INSERT IGNORE INTO sys_role_menu (role_id, menu_id) VALUES (1, 6), (1, 76), (1, 77), (1, 78);
